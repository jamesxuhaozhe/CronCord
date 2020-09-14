package worker

import (
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
	)
}

func (djm *DefaultJobManager) watchKillEvents() error {
	panic("implement me")
}

func (djm *DefaultJobManager) createLock() DistributedLock {
	panic("implement me")
}




