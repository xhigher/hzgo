package service

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/srd"
)

type Server struct {
	Hz   *server.Hertz
	Conf *config.ServerConfig
}

func NewServer(conf *config.ServerConfig) *Server {
	logger.Init(conf.Logger)
	mysql.Init(conf.Mysql)

	if conf.Srd != nil && srd.Init(conf.Srd) {
		if ok, option := srd.GetRegistry(); ok {
			return &Server{
				Hz: server.Default(server.WithHostPorts(conf.Addr),
					server.WithMaxRequestBodySize(conf.MaxReqSize),
					option),
				Conf: conf,
			}
		}
	}

	return &Server{
		Hz: server.Default(server.WithHostPorts(conf.Addr),
			server.WithMaxRequestBodySize(conf.MaxReqSize)),
		Conf: conf,
	}
}

func (s *Server) InitRouter(mgr RouterManager) {
	routers := mgr.Routers()
	for _, r := range routers {
		if r.Version == 0 {
			r.Version = 1
		}
		path := fmt.Sprintf("/%s/v%d/%s", mgr.Name(), r.Version, r.Path)
		s.Hz.POST(path, r.Handler)
	}
}

func (s *Server) Start() {

	s.Hz.Spin()
}
