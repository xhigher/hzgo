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
	"github.com/xhigher/hzgo/httpcli"
	"net/http"
	"time"

	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/server/gateway/middlewares"
	"github.com/xhigher/hzgo/utils"
)

type HzgoNodeServer struct {
	Conf              *config.ServerConfig
	Hz                *server.Hertz
	Auth              *middlewares.JWTAuth
	Sign              *middlewares.SecSign
	BroadcastMessageC chan Message
	// DirectMessageC direct messages are pushed to this channel
	DirectMessageC chan Message
	// Receive keeps a list of chan ChatMessage, one per user
	Receive     map[string]chan Message
	tokenHelper TokenHelper
	Id          string
}

func NewNodeServer(conf *config.ServerConfig) *HzgoNodeServer {
	logger.Init(conf.Logger)
	mysql.Init(conf.Mysql)
	fmt.Println("server config: ", utils.JSONString(conf))

	hz := server.Default(server.WithHostPorts(conf.Addr),
		server.WithExitWaitTime(consts.TimeSecond1),
		server.WithMaxRequestBodySize(conf.MaxReqSize))

	svr := &HzgoNodeServer{
		Conf: conf,
		Hz:   hz,
		//Auth:              middlewares.NewJWTAuth(conf.JWT),
		//Sign:              middlewares.NewSecSign(conf.Sec),
		BroadcastMessageC: make(chan Message),
		DirectMessageC:    make(chan Message),
		Receive:           make(map[string]chan Message),
		tokenHelper:       newTokenHelper(conf.JWT),
		Id:                utils.TimeUUID(),
	}

	go svr.relay()

	hz.Use(svr.CorsHandle(), svr.CreateReceiveChannel())
	hz.GET("/sse", svr.ServerSentEvent)

	return svr
}

func (s *HzgoNodeServer) ServerSentEvent(ctx context.Context, c *app.RequestContext) {
	// in production, you would get user's identity in other ways e.g. Authorization
	uid, ok := ctx.Value("uid").(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	stream := sse.NewStream(c)
	go func() {
		s.Heartbeat(stream)
	}()

	s.register(uid)

	// get messages from user's receive channel
	for msg := range s.Receive[uid] {

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

func (s *HzgoNodeServer) Heartbeat(stream *sse.Stream) {
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

// relay handles messages sent to BroadcastMessageC and DirectMessageC and
// relay messages to receive channels depends on message type.
func (s *HzgoNodeServer) relay() {
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

func (s *HzgoNodeServer) register(uid string) (err error) {
	params := RegisterReq{
		Uid: uid,
		Nid: s.Id,
	}
	url := fmt.Sprintf("http://%s/register", s.Conf.MasterAddr)
	err = httpcli.PostJSON(url, params, nil)
	return
}

func (s *HzgoNodeServer) CorsHandle() app.HandlerFunc {
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
func (s *HzgoNodeServer) CreateReceiveChannel() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.Query("token")
		if len(token) < 30 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		hlog.CtxInfof(ctx, "token: %v", token)
		uid, did, err := s.tokenHelper.ParseTokenInfo(token)
		if err != nil {
			hlog.Errorf("token: %v, error: %v", token, err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "uid", uid)
		ctx = context.WithValue(ctx, "did", did)
		// if user doesn't have a channel yet, create a new one.
		if _, found := s.Receive[uid]; !found {
			receive := make(chan Message, 1000)
			s.Receive[uid] = receive
		}
		c.Next(ctx)
	}
}

func (s *HzgoNodeServer) ConnectMaster() {
	c := sse.NewClient(fmt.Sprintf("http://%s/master?nid=%s", s.Conf.MasterAddr, s.Id))

	// touch off when connected to the server
	c.SetOnConnectCallback(func(ctx context.Context, client *sse.Client) {
		hlog.Infof("node connect to master server %s success", c.GetURL())
	})

	// touch off when the connection is shutdown
	c.SetDisconnectCallback(func(ctx context.Context, client *sse.Client) {
		hlog.Infof("node disconnect to master server %s success", c.GetURL())
	})

	err := c.Subscribe(func(e *sse.Event) {
		if e.Data != nil {
			msg := Message{}
			json.Unmarshal(e.Data, &msg)
			if e.Event == "broadcast" {
				s.BroadcastMessageC <- msg
			} else if e.Event == "direct" {
				s.DirectMessageC <- msg
			}
			return
		}
	})
	if err != nil {
		hlog.Errorf("node subscribe to master server failed: nid: %v, error: %v", s.Id, err)
		time.Sleep(3 * time.Second)
		s.ConnectMaster()
	}
}

func (s *HzgoNodeServer) Start() {

	go func() {
		s.ConnectMaster()
	}()

	s.Hz.Spin()
}
