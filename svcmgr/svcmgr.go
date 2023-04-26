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
	"github.com/xhigher/hzgo/srd"
	urllib "net/url"
	"time"
)

const (
	contentTypeJSON = "application/json"
	contentTypeForm = "application/x-www-form-urlencoded"
	contentTypeXML  = "application/xml"
)

var (
	svcClients = make(map[string]*SvcClient)
)

type SvcConf struct {
	Name     string
	AddrList []string
	Timeout  time.Duration
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
		svcClients[c.Name] = newClient(c.Name, c.AddrList[0], c.Timeout)
	}
}

func GetClient(name string) *SvcClient {
	return svcClients[name]
}

type SvcClient struct {
	Name    string
	Addr    string
	timeout time.Duration
	cli     *client.Client
}

func newClient(name, addr string, timeout time.Duration) *SvcClient {
	cli, err := client.NewClient()
	if err != nil {
		logger.Errorf("svc client [%v] new error: %v", name, err)
		return nil
	}
	if ok, discovery := srd.GetDiscovery(); ok {
		cli.Use(discovery)
	}

	logger.Infof("svc client [%v] init done", name)
	return &SvcClient{
		Name:    name,
		Addr:    addr,
		cli:     cli,
		timeout: timeout,
	}
}

func (s *SvcClient) Do(uri string, data interface{}) (result resp.BaseResp) {
	if s.cli == nil {
		result = resp.ErrorInternal
		logger.Errorf("svc client [%s][%s] not init: %v", s.Name, uri, s.Addr)
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		result = resp.ErrorInternal
		logger.Errorf("svc client [%s][%s] data error: %v, %v", s.Name, uri, err, data)
		return
	}

	req := &protocol.Request{}
	res := &protocol.Response{}
	req.SetBody(bytes)
	req.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	url := fmt.Sprintf("http://%s/%s/%s", s.Addr, s.Name, uri)
	req.SetRequestURI(url)
	err = s.cli.DoTimeout(context.Background(), req, res, s.timeout)
	if err != nil {
		result = resp.ErrorInternal
		logger.Errorf("svc client [%s][%s] resp error: %v, %v", s.Name, uri, err, url)
		return
	}
	logger.Infof("svc client [%s][%s] resp: %v", s.Name, uri, string(res.Body()))
	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		result = resp.ErrorInternal
		logger.Errorf("svc client [%s][%s] resp error: %v, %v", s.Name, uri, err, string(res.Body()))
		return
	}
	return
}


func (s *SvcClient) GetJSON(url string, result interface{}) (err error) {
	_, err = s.BaseGetJSON(url, nil, nil, result)
	return
}

func (s *SvcClient) GetJSON2(url string, params map[string]string, result interface{}) (err error) {
	_, err = s.BaseGetJSON(url, nil, params, result)
	return
}

func (s *SvcClient) PostJSON(url string, data, result interface{}) (err error) {
	_, err = s.BasePostJSON(url, nil, data, result)
	return
}

func (s *SvcClient) PostForm(url string, data map[string]string, result interface{}) (err error) {
	_, err = s.BasePostForm(url, nil, data, result)
	return
}

func (s *SvcClient) BaseGetJSON(url string, headers, params map[string]string, result interface{}) (bs string, err error) {
	if s.cli == nil {
		logger.Errorf("http cli get error: client nil")
		return
	}

	urlSt, err := urllib.Parse(url)
	if err != nil {
		logger.Errorf("http cli get url :%v, error: %v", url, err)
		return
	}
	values := urlSt.Query()
	for k, v := range params {
		values.Set(k, v)
	}
	urlSt.RawQuery = values.Encode()
	url = urlSt.String()
	req := &protocol.Request{}
	resp := &protocol.Response{}
	req.SetMethod(consts.MethodGet)
	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
	req.Header.SetContentTypeBytes([]byte(contentTypeJSON))
	req.SetRequestURI(url)
	err = s.cli.DoTimeout(context.Background(), req, resp, s.timeout)
	if err != nil {
		logger.Errorf("http cli url: %v, error: %v", url, err)
		return
	}
	bb := resp.Body()
	bs = string(bb)
	logger.Infof("http cli url: %v, resp: %v", url, bs)
	if result != nil {
		err = json.Unmarshal(bb, &result)
		if err != nil {
			logger.Errorf("http cli url: %v, resp: %v, error: %v", url, bs, err)
			return
		}
	}
	return
}

func (s *SvcClient) BasePostJSON(url string, headers map[string]string, data, result interface{}) (bs string, err error) {
	if s.cli == nil {
		logger.Errorf("http cli post error: client nil")
		return
	}

	body, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("http cli post url :%v, data: %v, error: %v", url, data, err)
		return
	}

	req := &protocol.Request{}
	resp := &protocol.Response{}
	req.SetMethod(consts.MethodPost)
	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
	req.Header.SetContentTypeBytes([]byte(contentTypeJSON))
	req.SetBody(body)
	req.SetRequestURI(url)
	err = s.cli.DoTimeout(context.Background(), req, resp, s.timeout)
	if err != nil {
		logger.Errorf("http cli post url: %v, error: %v", url, err)
		return
	}
	bb := resp.Body()
	bs = string(bb)
	logger.Infof("http cli post url: %v, resp: %v", url, bs)
	if result != nil {
		err = json.Unmarshal(bb, &result)
		if err != nil {
			logger.Errorf("http cli post url: %v, resp: %v, error: %v", url, bs, err)
			return
		}
	}
	return
}

func (s *SvcClient) BasePostForm(url string, headers, data map[string]string, result interface{}) (bs string, err error) {
	if s.cli == nil {
		logger.Errorf("http cli post error: client nil")
		return
	}

	req := &protocol.Request{}
	resp := &protocol.Response{}
	req.SetMethod(consts.MethodPost)
	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
	req.Header.SetContentTypeBytes([]byte(contentTypeForm))
	req.SetFormData(data)
	req.SetRequestURI(url)
	err = s.cli.DoTimeout(context.Background(), req, resp, s.timeout)
	if err != nil {
		logger.Errorf("http cli post url: %v, error: %v", url, err)
		return
	}
	bb := resp.Body()
	bs = string(bb)
	logger.Infof("http cli post url: %v, resp: %v", url, bs)
	if result != nil {
		err = json.Unmarshal(bb, &result)
		if err != nil {
			logger.Errorf("http cli post url: %v, resp: %v, error: %v", url, bs, err)
			return
		}
	}
	return
}

