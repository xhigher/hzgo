package game

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/websocket"
	"log"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/utils"
)

type HzgoServer struct {
	Conf     *config.ServerConfig
	Hz       *server.Hertz
	hub      *Hub
	upgrader websocket.HertzUpgrader
}

func NewServer(conf *config.ServerConfig, handler Handler) *HzgoServer {
	logger.Init(conf.Logger)
	mysql.Init(conf.Mysql)
	fmt.Println("server config: ", utils.JSONString(conf))
	hz := server.Default(server.WithHostPorts(conf.Addr),
		server.WithExitWaitTime(consts.TimeSecond1),
		server.WithMaxRequestBodySize(conf.MaxReqSize))
	hz.NoHijackConnPool = true
	svr := &HzgoServer{
		Conf: conf,
		Hz:   hz,
		hub:  newHub(handler),
		upgrader: websocket.HertzUpgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	go svr.hub.run()

	return svr
}

func (s *HzgoServer) Start(api API) {
	for _, r := range api.Routers() {
		switch r.Method {
		case consts.MethodPost:
			s.Hz.POST("/api/"+r.Path, r.Handler)
		case consts.MethodGet:
			s.Hz.GET("/api/"+r.Path, r.Handler)
		}
	}

	s.Hz.GET("/ws", s.wsHandler)
	s.Hz.Spin()
}

func (s *HzgoServer) wsHandler(_ context.Context, c *app.RequestContext) {
	err := s.upgrader.Upgrade(c, func(conn *websocket.Conn) {
		pipe := NewPipe(s.hub, conn)
		go pipe.writePump()
		pipe.readPump()
	})
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
}
