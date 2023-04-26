package srd

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
	"log"
	"github.com/xhigher/hzgo/env"
)

func (mgr *Manager) initConsul(conf *ConsulConfig) (err error) {
	if len(conf.ServerAddr) == 0 {
		log.Fatal(err)
		return
	}
	if conf.Registry {
		err = mgr.initConsulRegistry(conf)
		if err != nil {
			return
		}
	}

	if conf.Resolver {
		err = mgr.initConsulResolver(conf)
		if err != nil {
			return
		}
	}
	return
}

func (mgr *Manager) initConsulRegistry(conf *ConsulConfig) (err error) {
	config := consulapi.DefaultConfig()
	config.Address = conf.ServerAddr
	consulClient, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	// build a consul register with the consul client
	reg := consul.NewConsulRegister(consulClient)

	// run Hertz with the consul register
	localIP := env.GetLocalIP()
	if len(localIP) == 0 {
		log.Fatal(err)
		return
	}

	mgr.localAddr = fmt.Sprintf("%s:%d", localIP, conf.LocalPort)
	mgr.registry = reg
	return
}

func (mgr *Manager) initConsulResolver(conf *ConsulConfig) (err error) {
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = conf.ServerAddr
	consulClient, err := consulapi.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	// build a consul resolver with the consul client
	mgr.resolver = consul.NewConsulResolver(consulClient)
	return
}
