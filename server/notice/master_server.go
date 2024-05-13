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
	"github.com/xhigher/hzgo/resp"
	"net/http"
	"time"

	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/utils"
)

type HzgoMasterServer struct {
	Conf              *config.ServerConfig
	Hz                *server.Hertz
	BroadcastMessageC chan Message
	// DirectMessageC direct messages are pushed to this channel
	DirectMessageC chan Message
	// Receive keeps a list of chan ChatMessage, one per user
	Receive     map[string]chan Message
	tokenHelper TokenHelper
	Clients     map[string]string
}

func NewMasterServer(conf *config.ServerConfig) *HzgoMasterServer {
	logger.Init(conf.Logger)
	mysql.Init(conf.Mysql)
	fmt.Println("server config: ", utils.JSONString(conf))

	hz := server.Default(server.WithHostPorts(conf.Addr),
		server.WithExitWaitTime(consts.TimeSecond1),
		server.WithMaxRequestBodySize(conf.MaxReqSize))

	svr := &HzgoMasterServer{
		Conf:              conf,
		Hz:                hz,
		BroadcastMessageC: make(chan Message),
		DirectMessageC:    make(chan Message),
		Receive:           make(map[string]chan Message),
		tokenHelper:       newTokenHelper(conf.JWT),
		Clients:           map[string]string{},
	}

	go svr.relay()

	hz.Use(svr.CreateReceiveChannel())
	hz.Use(svr.CorsHandle())
	hz.GET("/master", svr.ServerSentEvent)
	hz.POST("/token", svr.Token)
	hz.POST("/register", svr.Register)
	hz.POST("/message/broadcast", svr.MessageBroadcast)
	hz.POST("/message/direct", svr.MessageDirect)

	return svr
}

func (s *HzgoMasterServer) ServerSentEvent(ctx context.Context, c *app.RequestContext) {
	nid := c.Query("nid")
	stream := sse.NewStream(c)
	go func() {
		s.Heartbeat(stream)
	}()
	for msg := range s.Receive[nid] {
		c.SetStatusCode(http.StatusOK)
		payload, err := json.Marshal(msg)
		if err != nil {
			hlog.CtxInfof(ctx, "message error: %+v, %v", msg, err)
			event := &sse.Event{
				Event: "error",
				Data:  []byte(err.Error()),
			}
			err = stream.Publish(event)
			if err != nil {
				return
			}
		}
		hlog.CtxInfof(ctx, "message received: %+v", msg)
		event := &sse.Event{
			Event: msg.Type,
			Data:  payload,
		}
		err = stream.Publish(event)
		if err != nil {
			return
		}
	}
}

func (s *HzgoMasterServer) Heartbeat(stream *sse.Stream) {
	for t := range time.NewTicker(3 * time.Second).C {
		event := &sse.Event{
			Event: "heartbeat",
			Data:  []byte(t.Format(time.RFC3339)),
		}
		err := stream.Publish(event)
		if err != nil {
			return
		}
	}
}

func (s *HzgoMasterServer) Register(ctx context.Context, c *app.RequestContext) {
	responder := resp.Responder{
		Ctx: c,
	}
	params := RegisterReq{}
	if err := c.Bind(&params); err != nil {
		responder.ReplyErrorParam()
		return
	}

	s.Clients[params.Uid] = params.Nid
	hlog.CtxInfof(ctx, "register uid: %v, nid: %v", params.Uid, params.Nid)
	responder.ReplyOK()
}

func (s *HzgoMasterServer) MessageBroadcast(ctx context.Context, c *app.RequestContext) {
	responder := resp.Responder{
		Ctx: c,
	}
	params := MessageBroadcastReq{}
	if err := c.Bind(&params); err != nil {
		responder.ReplyErrorParam()
		return
	}

	msg := Message{
		From:      params.From,
		To:        params.To,
		Data:      params.Data,
		Type:      "broadcast",
		Timestamp: time.Now(),
	}
	// deliver message to DirectMessageC.
	s.BroadcastMessageC <- msg

	hlog.CtxInfof(ctx, "message sent: %+v", msg)
	c.AbortWithStatus(http.StatusOK)
}

func (s *HzgoMasterServer) MessageChannel(ctx context.Context, c *app.RequestContext) {
	responder := resp.Responder{
		Ctx: c,
	}
	params := MessageBroadcastReq{}
	if err := c.Bind(&params); err != nil {
		responder.ReplyErrorParam()
		return
	}
	msg := Message{
		From:      params.From,
		To:        params.To,
		Data:      params.Data,
		Type:      "channel",
		Timestamp: time.Now(),
	}
	// deliver message to DirectMessageC.
	s.DirectMessageC <- msg

	hlog.CtxInfof(ctx, "message sent: %+v", msg)
	c.AbortWithStatus(http.StatusOK)
}

func (s *HzgoMasterServer) MessageDirect(ctx context.Context, c *app.RequestContext) {
	responder := resp.Responder{
		Ctx: c,
	}
	params := MessageBroadcastReq{}
	if err := c.Bind(&params); err != nil {
		responder.ReplyErrorParam()
		return
	}
	msg := Message{
		From:      params.From,
		To:        params.To,
		Data:      params.Data,
		Type:      "direct",
		Timestamp: time.Now(),
	}
	// deliver message to DirectMessageC.
	s.DirectMessageC <- msg

	hlog.CtxInfof(ctx, "message sent: %+v", msg)
	c.AbortWithStatus(http.StatusOK)
}

// relay handles messages sent to BroadcastMessageC and DirectMessageC and
// relay messages to receive channels depends on message type.
func (s *HzgoMasterServer) relay() {
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

type TokenInfo struct {
	Uid   string `json:"uid"`
	Did   string `json:"did"`
	Token string `json:"token"`
}

func (s *HzgoMasterServer) Token(ctx context.Context, c *app.RequestContext) {
	// in production, you would get user's identity in other ways e.g. Authorization
	uid := c.Query("uid")
	did := c.Query("did")

	responder := resp.Responder{
		Ctx: c,
	}

	token, err := s.tokenHelper.CreateToken(uid, did)
	if err != nil {
		responder.ReplyErrorInternal()
		return
	}

	data := &TokenInfo{
		Uid:   uid,
		Did:   did,
		Token: token,
	}
	hlog.CtxInfof(ctx, "token: %+v", data)
	responder.ReplyData(data)
}

func (s *HzgoMasterServer) CorsHandle() app.HandlerFunc {
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
func (s *HzgoMasterServer) CreateReceiveChannel() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		nid := c.Query("nid")
		if _, found := s.Receive[nid]; !found {
			receive := make(chan Message, 1000)
			s.Receive[nid] = receive
		}
		c.Next(ctx)
	}
}

func (s *HzgoMasterServer) Start() {

	s.Hz.Spin()
}
