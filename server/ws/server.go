package ws

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/websocket"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/utils"
	"log"
)

type HzgoServer struct {
	Conf *config.ServerConfig
	Hz   *server.Hertz
	hub *Hub
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
		Hz: hz,
		hub: newHub(handler),
		upgrader: websocket.HertzUpgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	go svr.hub.run()

	return svr
}

func (s *HzgoServer) Start() {
	s.Hz.GET("/"+s.Conf.Name, s.handle)
	s.Hz.Spin()
}

func (s *HzgoServer) handle(_ context.Context, c *app.RequestContext) {
	err := s.upgrader.Upgrade(c, func(conn *websocket.Conn) {
		pipe := &Pipe{hub: s.hub, conn: conn, send: make(chan []byte, 256)}
		pipe.hub.register <- pipe

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
