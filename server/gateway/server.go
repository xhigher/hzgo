package gateway

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/server/gateway/middlewares"
	"github.com/xhigher/hzgo/utils"
)

type HzgoServer struct {
	Conf *config.ServerConfig
	Hz   *server.Hertz
	Auth *middlewares.JWTAuth
	Sign *middlewares.SecSign
}

func NewServer(conf *config.ServerConfig) *HzgoServer {
	logger.Init(conf.Logger)
	mysql.Init(conf.Mysql)
	fmt.Println("server config: ", utils.JSONString(conf))
	return &HzgoServer{
		Conf: conf,
		Hz: server.Default(server.WithHostPorts(conf.Addr),
			server.WithMaxRequestBodySize(conf.MaxReqSize)),
		Auth: middlewares.NewJWTAuth(conf.JWT),
		Sign: middlewares.NewSecSign(conf.Sec),
	}
}

func (s *HzgoServer) InitRouter(mgr RouterManager) {
	routers := mgr.Routers()
	for _, r := range routers {
		if r.Version == 0 {
			r.Version = 1
		}
		path := fmt.Sprintf("/%s/v%d/%s", mgr.Name(), r.Version, r.Path)
		handlers := make([]app.HandlerFunc, 0)
		if r.Sign {
			handlers = append(handlers, s.Sign.Verify())
		}
		if r.Auth {
			handlers = append(handlers, s.Auth.Authenticate())
		}
		handlers = append(handlers, r.Handler)
		switch r.Method {
		case consts.MethodPost:
			s.Hz.POST(path, handlers...)
		case consts.MethodGet:
			s.Hz.GET(path, handlers...)
		case consts.MethodPut:
			s.Hz.PUT(path, handlers...)
		case consts.MethodDelete:
			s.Hz.DELETE(path, handlers...)
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
