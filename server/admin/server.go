package admin

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/utils"
)

type HzgoServer struct {
	Conf *config.ServerConfig
	Hz   *server.Hertz
	Auth *Auth
}

func NewServer(conf *config.ServerConfig) *HzgoServer {
	logger.Init(conf.Logger)
	mysql.Init(conf.Mysql)
	fmt.Println("server config: ", utils.JSONString(conf))
	return &HzgoServer{
		Conf: conf,
		Hz: server.Default(server.WithHostPorts(conf.Addr),
			server.WithExitWaitTime(consts.TimeSecond1),
			server.WithMaxRequestBodySize(conf.MaxReqSize)),
		Auth: NewAuth(conf.JWT),
	}
}

func (s *HzgoServer) InitRouters(mgr ModuleManager) {
	s.initPlatformModuleRouters(mgr)

	for _, m := range mgr.Modules() {
		for _, r := range m.Routers() {
			if r.Version == 0 {
				r.Version = 1
			}
			r.mergeRoles(m.Roles())

			path := r.FullPath(m.Name())
			handlers := []app.HandlerFunc{
				s.Auth.Handler(),
				r.PermissionFunc(m),
				r.Handler,
			}
			switch r.Method {
			case consts.MethodPost:
				s.Hz.POST(path, handlers...)
			case consts.MethodGet:
				s.Hz.GET(path, handlers...)
			}
		}
	}
}

func (s *HzgoServer) initPlatformModuleRouters(mgr ModuleManager){
	m := PlatformModule{
		ctrl: mgr.BaseController(),
		handler: mgr.PlatformHandler(),
	}
	for _, r := range m.Routers() {
		if r.Version == 0 {
			r.Version = 1
		}
		path := r.FullPath(m.Name())
		handlers := []app.HandlerFunc{
			s.Auth.Handler(),
			r.Handler,
		}
		if r.NoAuth {
			handlers = []app.HandlerFunc{
				r.Handler,
			}
		}
		switch r.Method {
		case consts.MethodPost:
			s.Hz.POST(path, handlers...)
		case consts.MethodGet:
			s.Hz.GET(path, handlers...)
		}
	}
}

func (s *HzgoServer) Start() {
	if s.Auth == nil {
		hlog.Fatalf("auth nil")
		return
	}

	s.Hz.Spin()
}
