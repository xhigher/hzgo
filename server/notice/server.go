package notice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	Hz                *server.Hertz
	Auth              *middlewares.JWTAuth
	Sign              *middlewares.SecSign
	BroadcastMessageC chan Message
	// DirectMessageC direct messages are pushed to this channel
	DirectMessageC chan Message
	// Receive keeps a list of chan ChatMessage, one per user
	Receive map[string]chan Message
}

func NewSSEServer(conf *config.ServerConfig) *HzgoServer {
	logger.Init(conf.Logger)
	mysql.Init(conf.Mysql)
	fmt.Println("server config: ", utils.JSONString(conf))

	hz := server.Default(server.WithHostPorts(conf.Addr),
		server.WithExitWaitTime(consts.TimeSecond1),
		server.WithMaxRequestBodySize(conf.MaxReqSize))

	svr := &HzgoServer{
		Conf:              conf,
		Hz:                hz,
		Auth:              middlewares.NewJWTAuth(conf.JWT),
		Sign:              middlewares.NewSecSign(conf.Sec),
		BroadcastMessageC: make(chan Message),
		DirectMessageC:    make(chan Message),
		Receive:           make(map[string]chan Message),
	}

	go svr.relay()

	hz.Use(svr.CreateReceiveChannel())

	hz.GET("/sse", svr.ServerSentEvent)
	hz.POST("/chat/broadcast", svr.Broadcast)
	hz.POST("/chat/direct", svr.Direct)

	return svr
}

func (srv *HzgoServer) ServerSentEvent(ctx context.Context, c *app.RequestContext) {
	// in production, you would get user's identity in other ways e.g. Authorization
	username := c.Query("username")

	stream := sse.NewStream(c)
	// get messages from user's receive channel
	for msg := range srv.Receive[username] {

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

func (srv *HzgoServer) Direct(ctx context.Context, c *app.RequestContext) {
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
	srv.DirectMessageC <- msg

	hlog.CtxInfof(ctx, "message sent: %+v", msg)
	c.AbortWithStatus(http.StatusOK)
}

func (srv *HzgoServer) Broadcast(ctx context.Context, c *app.RequestContext) {
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
	srv.BroadcastMessageC <- msg

	hlog.CtxInfof(ctx, "message sent: %+v", msg)
	c.AbortWithStatus(http.StatusOK)
}

// relay handles messages sent to BroadcastMessageC and DirectMessageC and
// relay messages to receive channels depends on message type.
func (srv *HzgoServer) relay() {
	for {
		select {
		// broadcast message to all users
		case msg := <-srv.BroadcastMessageC:
			for _, r := range srv.Receive {
				r <- msg
			}

		// deliver message to user specified in To
		case msg := <-srv.DirectMessageC:
			srv.Receive[msg.To] <- msg
		}
	}
}

// CreateReceiveChannel creates a buffered receive channel for each user.
func (srv *HzgoServer) CreateReceiveChannel() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		username := c.Query("username")
		// if user doesn't have a channel yet, create a new one.
		if _, found := srv.Receive[username]; !found {
			receive := make(chan Message, 1000)
			srv.Receive[username] = receive
		}
		c.Next(ctx)
	}
}

func (s *HzgoServer) Start() {
	if s.Auth == nil {
		hlog.Fatalf("auth nil")
		return
	}

	s.Hz.Spin()
}
