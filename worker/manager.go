package worker

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jamesxuhaozhe/croncord/worker/model"
)

type JobManager interface {
	watchJobEvents() error
	watchKillEvents() error
	createLock() DistributedLock
}

type DefaultJobManager struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
	watcher clientv3.Watcher
}

var WorkerJobManager JobManager

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

func (djm *DefaultJobManager) createLock() DistributedLock {
	panic("implement me")
}




