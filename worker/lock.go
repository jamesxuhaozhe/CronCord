package worker

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
)
type DistributedLock interface {
	TryLock() error
	Unlock() error
}

type EtcdLock struct {
	kv    clientv3.KV
	lease clientv3.Lease

	jobName    string
	cancelFunc context.CancelFunc
	leaseId    clientv3.LeaseID
	isLocked   bool
}

func (e *EtcdLock) TryLock() error {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
		err            error
	)

	if leaseGrantResp, err = e.lease.Grant(context.TODO(), 5); err != nil {
		return err
	}

	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	leaseId = leaseGrantResp.ID

	if keepRespChan, err = e.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}

	go func() {
		var keepResp *clientv3.LeaseKeepAliveResponse
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					break
				}
			}
		}
	}()

	txn = e.kv.Txn(context.TODO())

	lockKey = "/croncord/lock/" + e.jobName

	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	if !txnResp.Succeeded {
		err = errors.New("lock unable to acquire")
		goto FAIL
	}

	e.leaseId = leaseId
	e.cancelFunc = cancelFunc
	e.isLocked = true

FAIL:
	cancelFunc()
	e.lease.Revoke(context.TODO(), leaseId)
	return err
}

func (e *EtcdLock) Unlock() error{
	if e.isLocked {
		e.cancelFunc()
		if _, err := e.lease.Revoke(context.TODO(), e.leaseId); err != nil {
			return err
		}
	}
	return nil
}

func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) DistributedLock {
	lock := &EtcdLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return lock
}

