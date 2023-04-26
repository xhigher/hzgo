package httpcli

import (
	"encoding/json"
	"github.com/ddliu/go-httpclient"
	"time"
	"github.com/xhigher/hzgo/logger"
)

func init() {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "hzgo http client",
	})
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
	client := httpclient.WithHeaders(headers)
	if timeout > 0 {
		client.WithOption(httpclient.OPT_TIMEOUT, timeout)
	}
	var _resp *httpclient.Response
	_resp, err = client.Get(url, data)
	if err != nil {
		logger.Errorf("http cli url: %v, error: %v", url, err)
		return
	}
	bb, err := _resp.ReadAll()
	if err != nil {
		logger.Errorf("http cli url: %v, error: %v", url, err)
		return
	}
	if resp != nil {
		err = json.Unmarshal(bb, &resp)
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
	client := httpclient.WithHeaders(headers)
	if timeout > 0 {
		client.WithOption(httpclient.OPT_TIMEOUT, timeout)
	}
	var _resp *httpclient.Response
	_resp, err = client.PostJson(url, data)
	if err != nil {
		logger.Errorf("http cli url: %v, error %s", url, err)
		return
	}
	bb, err := _resp.ReadAll()
	if err != nil {
		logger.Errorf("http cli url: %v, error: %v", url, err)
		return
	}
	if resp != nil {
		err = json.Unmarshal(bb, &resp)
		if err != nil {
			logger.Errorf("http cli url: %v, error: %v", url, err)
			return
		}
	}
	return
}
