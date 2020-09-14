package worker

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jamesxuhaozhe/croncord/worker/model"
	"time"
)

type JobManager interface {
	watchJobEvents() error
	watchKillEvents() error
	createLock(jobName string) DistributedLock
}

type DefaultJobManager struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
	watcher clientv3.Watcher
}

var WorkerJobManager JobManager

func InitJobManager() error {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
		watcher clientv3.Watcher
		err error
	)

	config = clientv3.Config{
		Endpoints:WorkerConfig.EtcdEndpoints,
		DialTimeout:time.Duration(WorkerConfig.EtcdDialTimeout) * time.Millisecond,
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
		getResp *clientv3.GetResponse
		kvpair *mvccpb.KeyValue
		job *model.Job
		watchStartRevision int64
		watchChan clientv3.WatchChan
		watchResp clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobName string
		jobEvent *model.JobEvent
		err error
	)

	if getResp, err = djm.kv.Get(context.TODO(), "/croncord/jobs/", clientv3.WithPrefix()); err != nil {
		return err
	}

	for _, kvpair = range getResp.Kvs {
		if job, err = unpackJob(kvpair.Value); err == nil {

		}
	}


}

func unpackJob(value []byte) (*model.Job, error) {
	job := &model.Job{}
	var err error
	if err = json.Unmarshal(value, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (djm *DefaultJobManager) watchKillEvents() error {
	panic("implement me")
}

func (djm *DefaultJobManager) createLock(jobName string) DistributedLock {
	lock := InitJobLock(jobName, djm.kv, djm.lease)
	return lock
}




