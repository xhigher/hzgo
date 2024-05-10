package notice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/sse"
	"net/http"
	"time"

	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/server/gateway/middlewares"
	"github.com/xhigher/hzgo/utils"
)

type HzgoServer struct {
	Conf              *config.ServerConfig
	OuterHz           *server.Hertz
	InnerHz           *server.Hertz
	Auth              *middlewares.JWTAuth
	Sign              *middlewares.SecSign
	BroadcastMessageC chan Message
	// DirectMessageC direct messages are pushed to this channel
	DirectMessageC chan Message
	// Receive keeps a list of chan ChatMessage, one per user
	Receive map[string]chan Message
}

func NewServer(conf *config.ServerConfig) *HzgoServer {
	logger.Init(conf.Logger)
	mysql.Init(conf.Mysql)
	fmt.Println("server config: ", utils.JSONString(conf))

	ohz := server.Default(server.WithHostPorts(conf.OuterAddr),
		server.WithExitWaitTime(consts.TimeSecond1),
		server.WithMaxRequestBodySize(conf.MaxReqSize))

	ihz := server.Default(server.WithHostPorts(conf.InnerAddr),
		server.WithExitWaitTime(consts.TimeSecond1),
		server.WithMaxRequestBodySize(conf.MaxReqSize))

	svr := &HzgoServer{
		Conf:    conf,
		OuterHz: ohz,
		InnerHz: ihz,
		//Auth:              middlewares.NewJWTAuth(conf.JWT),
		//Sign:              middlewares.NewSecSign(conf.Sec),
		BroadcastMessageC: make(chan Message),
		DirectMessageC:    make(chan Message),
		Receive:           make(map[string]chan Message),
	}

	go svr.relay()

	ohz.Use(svr.CreateReceiveChannel())
	ohz.Use(svr.CorsHandle())

	ohz.GET("/sse", svr.ServerSentEvent)

	ihz.POST("/chat/broadcast", svr.Broadcast)

	ihz.POST("/chat/direct", svr.Direct)

	return svr
}

func (s *HzgoServer) ServerSentEvent(ctx context.Context, c *app.RequestContext) {
	// in production, you would get user's identity in other ways e.g. Authorization
	username := c.Query("username")

	stream := sse.NewStream(c)
	go func() {
		s.Heartbeat(stream)
	}()
	// get messages from user's receive channel
	for msg := range s.Receive[username] {

		payload, err := json.Marshal(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		hlog.CtxInfof(ctx, "message received: %+v", msg)
		event := &sse.Event{
			Event: msg.Type,
			Data:  payload,
		}
		c.SetStatusCode(http.StatusOK)
		err = stream.Publish(event)
		if err != nil {
			return
		}
	}
}

func (s *HzgoServer) Heartbeat(stream *sse.Stream) {
	for t := range time.NewTicker(3 * time.Second).C {
		event := &sse.Event{
			Event: "timestamp",
			Data:  []byte(t.Format(time.RFC3339)),
		}
		err := stream.Publish(event)
		if err != nil {
			return
		}
	}
}

func (s *HzgoServer) Direct(ctx context.Context, c *app.RequestContext) {
	// in production, you would get user's identity in other ways e.g. Authorization
	from := c.Query("from")
	to := c.Query("to")
	message := c.Query("message")

	msg := Message{
		From:      from,
		To:        to,
		Message:   message,
		Type:      "direct",
		Timestamp: time.Now(),
	}
	// deliver message to DirectMessageC.
	s.DirectMessageC <- msg

	hlog.CtxInfof(ctx, "message sent: %+v", msg)
	c.AbortWithStatus(http.StatusOK)
}

func (s *HzgoServer) Broadcast(ctx context.Context, c *app.RequestContext) {
	// in production, you would get user's identity in other ways e.g. Authorization
	from := c.Query("from")
	message := c.Query("message")

	msg := Message{
		From:      from,
		Message:   message,
		Type:      "broadcast",
		Timestamp: time.Now(),
	}
	// deliver message to BroadcastMessageC.
	s.BroadcastMessageC <- msg

	hlog.CtxInfof(ctx, "message sent: %+v", msg)
	c.AbortWithStatus(http.StatusOK)
}

// relay handles messages sent to BroadcastMessageC and DirectMessageC and
// relay messages to receive channels depends on message type.
func (s *HzgoServer) relay() {
	for {
		select {
		// broadcast message to all users
		case msg := <-s.BroadcastMessageC:
			for _, r := range s.Receive {
				r <- msg
			}

		// deliver message to user specified in To
		case msg := <-s.DirectMessageC:
			s.Receive[msg.To] <- msg
		}
	}
}

func (s *HzgoServer) CorsHandle() app.HandlerFunc {
	if s.Conf.Cors == nil {
		s.Conf.Cors = &config.CorsConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET"},
			AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
		}
	}
	return cors.New(cors.Config{
		//准许跨域请求网站,多个使用,分开,限制使用*
		AllowOrigins: s.Conf.Cors.AllowOrigins,
		//准许使用的请求方式
		AllowMethods: s.Conf.Cors.AllowMethods,
		//准许使用的请求表头
		AllowHeaders: s.Conf.Cors.AllowHeaders,
		//显示的请求表头
		ExposeHeaders: []string{"Content-Type"},
		//凭证共享,确定共享
		AllowCredentials: true,
		//容许跨域的原点网站,可以直接return true就万事大吉了
		//AllowOriginFunc: func(origin string) bool {
		//	return true
		//},
		//超时时间设定
		MaxAge: 24 * time.Hour,
	})
}

// CreateReceiveChannel creates a buffered receive channel for each user.
func (s *HzgoServer) CreateReceiveChannel() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		username := c.Query("username")
		// if user doesn't have a channel yet, create a new one.
		if _, found := s.Receive[username]; !found {
			receive := make(chan Message, 1000)
			s.Receive[username] = receive
		}
		c.Next(ctx)
	}
}

func (s *HzgoServer) Start() {

	go func() {
		s.InnerHz.Spin()
	}()

	s.OuterHz.Spin()
}
