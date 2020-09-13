package main

import (
	"fmt"
	"github.com/jamesxuhaozhe/croncord/master"
	"github.com/jamesxuhaozhe/croncord/worker"
	"github.com/spf13/pflag"
)

var (
	role = pflag.StringP("mode", "m", "", "define deploy mode, either master or worker, like -m=master or -m=worker")
	cfg = pflag.StringP("config", "c", "", "worker config file path. like -c=./config.yaml")
)

func main() {
	pflag.Parse()

	if *role == "worker" || *role == "master" {
		if *role == "worker" {
			if err := worker.InitWorker(*cfg); err != nil {
				fmt.Printf("failed to init the worker config: %s\n", err.Error())
			}
		} else {
			if err := master.InitMaster(*cfg); err != nil {
				fmt.Printf("failed to init the master config: %s\n", err.Error())
			}
		}
	} else {
		fmt.Println("deployment role is not defined")
		return
	}

}
