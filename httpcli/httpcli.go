package httpcli

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	urllib "net/url"

	"github.com/xhigher/hzgo/logger"
	"time"
)

const (
	contentTypeJSON = "application/json"
	contentTypeForm = "application/x-www-form-urlencoded"
	contentTypeXML = "application/xml"
)

type HttpCli struct {
	client     *client.Client
	timeout time.Duration
}


func New(timeout time.Duration) *HttpCli{
	c, err := client.NewClient()
	if err != nil {
		logger.Errorf("http cli new error: %v", err)
		return nil
	}
	return &HttpCli{
		client: c,
		timeout: timeout,
	}
}

func (cli *HttpCli) GetJSON(url string, result interface{}) (err error){
	_, err = cli.BaseGetJSON(url, nil, nil, result)
	return
}

func (cli *HttpCli) GetJSON2(url string, params map[string]string, result interface{}) (err error){
	_, err = cli.BaseGetJSON(url, nil, params, result)
	return
}

func (cli *HttpCli) PostJSON(url string, data, result interface{}) (err error){
	_, err = cli.BasePostJSON(url, nil, data, result)
	return
}

func (cli *HttpCli) PostForm(url string, data map[string]string, result interface{}) (err error){
	_, err = cli.BasePostForm(url, nil, data, result)
	return
}

func (cli *HttpCli) BaseGetJSON(url string, headers, params map[string]string, result interface{}) (bs string, err error){
	if cli.client == nil {
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
		values.Set(k,v)
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
	err = cli.client.DoTimeout(context.Background(), req, resp, cli.timeout)
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

func (cli *HttpCli) BasePostJSON(url string, headers map[string]string, data, result interface{}) (bs string, err error){
	if cli.client == nil {
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
	err = cli.client.DoTimeout(context.Background(), req, resp, cli.timeout)
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

func (cli *HttpCli) BasePostForm(url string, headers, data map[string]string, result interface{}) (bs string, err error){
	if cli.client == nil {
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
	err = cli.client.DoTimeout(context.Background(), req, resp, cli.timeout)
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