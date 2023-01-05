package svcmgr

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
	"time"
)

var (
	svcClients = make(map[string]*SvcClient)
)

type SvcConf struct {
	Name     string
	AddrList []string
}

type ActionResult struct {
	Data resp.BaseResp
	Err  error
}

type BaseAction struct {
	SvcName string
	Name    string
	Ver     int
	Data    interface{}
}

func (a BaseAction) Do() resp.BaseResp {
	uri := fmt.Sprintf("v%d/%s", a.Ver, a.Name)
	if cli, ok := svcClients[a.SvcName]; ok {
		return cli.Do(uri, a.Data)
	}
	logger.Errorf("svc client [%v] not init", a.SvcName)
	return resp.ErrorInternal
}

func Init(confs []SvcConf) {
	for _, c := range confs {
		svcClients[c.Name] = newClient(c.Name, c.AddrList[0])
	}
}

type SvcClient struct {
	Name    string
	Addr    string
	Timeout time.Duration
	cli     *client.Client
}

func newClient(name, addr string) *SvcClient {
	cli, err := client.NewClient()
	if err != nil {
		logger.Errorf("svc client [%v] new error: %v", name, err)
		return nil
	}
	logger.Infof("svc client [%v] init done", name)
	return &SvcClient{
		Name: name,
		Addr: addr,
		cli:  cli,
		Timeout: 5*time.Second,
	}
}

func (s *SvcClient) Do(uri string, data interface{}) (result resp.BaseResp) {
	if s.cli == nil {
		result = resp.ErrorInternal
		logger.Errorf("svc client [%s] not init: %v", s.Name, s.Addr)
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		result = resp.ErrorInternal
		logger.Errorf("svc client [%s] data error: %v, %v", s.Name, err, data)
		return
	}

	req := &protocol.Request{}
	res := &protocol.Response{}
	req.SetBody(bytes)
	req.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	url := fmt.Sprintf("http://%s/%s/%s", s.Addr, s.Name, uri)
	req.SetRequestURI(url)
	err = s.cli.DoTimeout(context.Background(), req, res, s.Timeout)
	if err != nil {
		result = resp.ErrorInternal
		logger.Errorf("svc client [%s] resp error: %v, %v", s.Name, err, url)
		return
	}
	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		result = resp.ErrorInternal
		logger.Errorf("svc client [%s] resp error: %v, %v", s.Name, err, string(res.Body()))
		return
	}
	return
}
