package srd

import (
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/client/discovery"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"log"
)


type Manager struct {
	config *Config
	registry   registry.Registry
	resolver   discovery.Resolver
	serverAddr  string
	localAddr string
}

type Config struct {
	ServiceName string
	Consul *ConsulConfig
}

type ConsulConfig struct {
	ServerAddr string
	LocalPort int
	Registry bool
	Resolver bool
}

var manager *Manager


func Init(conf *Config) bool {
	if conf == nil {
		log.Fatal("srd config nil")
		return false
	}
	if manager == nil {
		manager = &Manager{
			config:conf,
		}
	}
	if conf.Consul != nil {
		err := manager.initConsul(conf.Consul)
		if err != nil {
			return false
		}
	}

	return true
}

func GetRegistry() (ok bool, option config.Option){
	if manager ==nil || manager.registry == nil {
		return
	}
	ok = true
	option = server.WithRegistry(manager.registry, &registry.Info{
		ServiceName: manager.config.ServiceName,
		Addr:        utils.NewNetAddr("tcp", manager.localAddr),
		Weight:      10,
		Tags:        nil,
	})
	return
}

func GetDiscovery() (ok bool, mw client.Middleware){
	if manager ==nil || manager.resolver == nil {
		return
	}
	ok = true
	mw = sd.Discovery(manager.resolver)
	return
}
