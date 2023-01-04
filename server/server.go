package server

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/xhigher/hzgo/auth"
	"github.com/xhigher/hzgo/log"
	"github.com/xhigher/hzgo/router"
)

type HzgoServer struct {
	Hz *server.Hertz
	Auth *auth.HzgoJWTMiddleware
}

func (s *HzgoServer) InitRouter(mgr router.Manager) {
	routers := mgr.Routers()
	for _, r := range routers {
		if r.Version == 0 {
			r.Version = 1
		}
		path := fmt.Sprintf("/v%d/%s", r.Version, r.Path)
		handlers := make([]app.HandlerFunc, 0)
		if r.Auth {
			handlers = append(handlers, s.Auth.MiddlewareFunc())
		}
		handlers = append(handlers, r.Handler)
		switch r.Method {
		case router.MethodPost:
			s.Hz.POST(path, handlers...)
		case router.MethodGet:
			s.Hz.GET(path, handlers...)
		case router.MethodPut:
			s.Hz.PUT(path, handlers...)
		case router.MethodDelete:
			s.Hz.DELETE(path, handlers...)
		}
	}
}

func (s *HzgoServer) InitAuth(mw *auth.HzgoJWTMiddleware) {
	var err error
	s.Auth, err = auth.New(mw)
	if err != nil {
		panic(err)
	}
	s.InitRouter(mw)
}

func New() *HzgoServer{
	log.Init("log")
	return &HzgoServer{
		Hz: server.Default(),
	}
}

func  (s *HzgoServer) Start(){
	if s.Auth == nil {
		hlog.Fatalf("auth nil")
		return
	}

	s.Hz.Spin()
}