package worker

import (
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"net"
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

func InitRegister() error {
	if ip, err := getLocalIP(); err != nil {
		return err
	} else {
		fmt.Printf("local ip: %s\n", ip)
		return nil
	}
}