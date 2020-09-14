package worker

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jamesxuhaozhe/croncord/worker/model"
	"strings"
	"time"
)

type JobManager interface {
	watchJobEvents() error
	watchKillEvents() error
	createLock(jobName string) DistributedLock
}

type DefaultJobManager struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var WorkerJobManager JobManager

func InitJobManager() error {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
		err     error
	)

	config = clientv3.Config{
		Endpoints:   WorkerConfig.EtcdEndpoints,
		DialTimeout: time.Duration(WorkerConfig.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return err
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	WorkerJobManager = &DefaultJobManager{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	WorkerJobManager.watchJobEvents()

	WorkerJobManager.watchKillEvents()

	return err
}

func (djm *DefaultJobManager) watchJobEvents() error {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *model.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *model.JobEvent
		err                error
	)

	if getResp, err = djm.kv.Get(context.TODO(), "/croncord/jobs/", clientv3.WithPrefix()); err != nil {
		return err
	}

	for _, kvpair = range getResp.Kvs {
		if job, err = unpackJob(kvpair.Value); err == nil {
			jobEvent = buildJobEvent(model.JobEventSave, job)
			// TODO need to push the even to scheduler

		}
	}

	go func() {

		watchStartRevision = getResp.Header.Revision + 1

		watchChan = djm.watcher.Watch(context.TODO(), model.JobSaveDir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())

		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = unpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobEvent = buildJobEvent(model.JobEventSave, job)
				case mvccpb.DELETE:
					jobName = extractJobName(string(watchEvent.Kv.Key))

					job = &model.Job{Name:jobName}

					jobEvent = buildJobEvent(model.JobEventDelete, job)
				}

				// todo need to push the event to scheduler

			}
		}
	}()
	return nil
}

func unpackJob(value []byte) (*model.Job, error) {
	job := &model.Job{}
	var err error
	if err = json.Unmarshal(value, job); err != nil {
		return nil, err
	}
	return job, nil
}

func buildJobEvent(eventType int, job *model.Job) (jobEvent *model.JobEvent) {
	return &model.JobEvent{
		EventType: eventType,
		Job: job,
	}
}

func extractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, model.JobSaveDir)
}

func extractKillerName(killerKey string) string {
	return strings.TrimPrefix(killerKey, model.JobKillerDir)
}

func (djm *DefaultJobManager) watchKillEvents() error {
	var (
		watchChan clientv3.WatchChan
		watchResp clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobEvent *model.JobEvent
		jobName string
		job *model.Job
	)

	go func() {

		watchChan = djm.watcher.Watch(context.TODO(), model.JobKillerDir, clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					jobName = extractKillerName(string(watchEvent.Kv.Key))
					job = &model.Job{Name:jobName}
					jobEvent = buildJobEvent(model.JobEventKill, job)
					// todo need to push this even to scheduler
					case mvccpb.DELETE:

				}
			}
		}
	}()
	return nil
}

func (djm *DefaultJobManager) createLock(jobName string) DistributedLock {
	lock := InitJobLock(jobName, djm.kv, djm.lease)
	return lock
}
