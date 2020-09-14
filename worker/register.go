package worker

import (
	"errors"
	"github.com/coreos/etcd/clientv3"
	"net"
	"time"
)

type Register struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	localIP string
}

var WorkerRegister *Register

func getLocalIP() (ipv4 string, err error) {
	var addrs []net.Addr
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, isIPNet := addr.(*net.IPNet); isIPNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("local ip not found")
}

// InitRegister register the worker node to the etcd config center
func InitRegister() (err error) {

	// init etcd client config
	config := clientv3.Config{
		Endpoints:            WorkerConfig.EtcdEndpoints,
		DialTimeout:          time.Duration(WorkerConfig.EtcdDialTimeout) * time.Millisecond,
	}

	// init the etcd client
	var client *clientv3.Client
	if client, err = clientv3.New(config); err != nil {
		return err
	}

	var localIP string
	if localIP, err = getLocalIP(); err != nil {
		return err
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)

	WorkerRegister = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIP,
	}

	return nil
}