package httpcli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/xhigher/hzgo/logger"
	"net/http"
	"time"
)

var defaultClient *client.Client

func init() {
	defaultClient, _ = client.NewClient(
		client.WithDialTimeout(1*time.Second),
		client.WithMaxConnsPerHost(1024),
		client.WithMaxIdleConnDuration(10*time.Second),
		client.WithMaxConnDuration(10*time.Second),
		client.WithMaxConnWaitTimeout(10*time.Second),
		client.WithClientReadTimeout(60*time.Second),
		client.WithWriteTimeout(60*time.Second),
		client.WithDialer(standard.NewDialer()),
		client.WithName("hzgo-http-client"))
}

func GetJSON(url string, data map[string]string, resp interface{}) (err error) {
	return GetJSONWithHeaders(url, map[string]string{}, data, resp)
}

func GetJSONWithTimeout(url string, data map[string]string, timeout time.Duration, resp interface{}) (err error) {
	return GetJSONWithHeadersAndTimeout(url, map[string]string{}, timeout, data, resp)
}

func GetJSONWithHeaders(url string, headers, data map[string]string, resp interface{}) (err error) {
	return GetJSONWithHeadersAndTimeout(url, headers, 0, data, resp)
}

func GetJSONWithHeadersAndTimeout(url string, headers map[string]string, timeout time.Duration, data map[string]string, resp interface{}) (err error) {
	req, res := &protocol.Request{}, &protocol.Response{}
	req.SetMethod(http.MethodGet)
	req.SetRequestURI(url)
	req.SetHeaders(headers)
	if timeout > 0 {
		req.SetOptions(config.WithRequestTimeout(5 * time.Second))
	}
	err = defaultClient.Do(context.Background(), req, res)
	if err != nil {
		logger.Errorf("http cli url: %v, error: %v", url, err)
		return
	}
	statusCode := res.StatusCode()
	if statusCode != http.StatusOK {
		logger.Errorf("http cli url: %v, error: %v", url, statusCode)
		err = fmt.Errorf("StatusCode=%d", statusCode)
		return
	}
	if resp != nil {
		err = json.Unmarshal(res.Body(), &resp)
		if err != nil {
			logger.Errorf("http cli url: %v, error: %v", url, err)
			return
		}
	}
	return
}

func PostJSON(url string, data interface{}, resp interface{}) (err error) {
	return PostJSONWithHeaders(url, map[string]string{}, data, resp)
}

func PostJSONWithHeaders(url string, headers map[string]string, data interface{}, resp interface{}) (err error) {
	return PostJSONWithHeadersAndTimeout(url, headers, 0, data, resp)
}

func PostJSONWithTimeout(url string, timeout time.Duration, data interface{}, resp interface{}) (err error) {
	return PostJSONWithHeadersAndTimeout(url, map[string]string{}, timeout, data, resp)
}

func PostJSONWithHeadersAndTimeout(url string, headers map[string]string, timeout time.Duration, data interface{}, resp interface{}) (err error) {
	req, res := &protocol.Request{}, &protocol.Response{}
	bs, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("http cli url: %v, error: %v", url, err)
		return
	}
	req.SetBody(bs)
	req.SetMethod(http.MethodPost)
	req.SetRequestURI(url)
	headers["Content-Type"] = "application/json"
	req.SetHeaders(headers)
	if timeout > 0 {
		req.SetOptions(config.WithRequestTimeout(5 * time.Second))
	}
	err = defaultClient.Do(context.Background(), req, res)
	if err != nil {
		logger.Errorf("http cli url: %v, error: %v", url, err)
		return
	}
	statusCode := res.StatusCode()
	if statusCode != http.StatusOK {
		logger.Errorf("http cli url: %v, error: %v", url, statusCode)
		err = fmt.Errorf("StatusCode=%d", statusCode)
		return
	}
	if resp != nil {
		err = json.Unmarshal(res.Body(), &resp)
		if err != nil {
			logger.Errorf("http cli url: %v, error: %v", url, err)
			return
		}
	}
	return
}
