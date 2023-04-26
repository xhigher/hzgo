package httpcli

import (
	"github.com/ddliu/go-httpclient"
	"time"
	"github.com/xhigher/hzgo/logger"
)

func init() {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "hzgo http client",
	})
}

func PostJSON(url string, data interface{}) (body string, err error) {
	return PostJSONWithHeaders(url, map[string]string{}, data)
}

func PostJSONWithHeaders(url string, headers map[string]string, data interface{}) (body string, err error) {
	client := httpclient.WithHeaders(headers)
	var resp *httpclient.Response
	resp, err = client.PostJson(url, data)
	if err != nil {
		logger.Errorf("error %s", err)
		return
	}
	body, err = resp.ToString()
	return
}

func Get(url string, data map[string]string) (body string, err error) {
	return GetWithHeaders(url, map[string]string{}, data)
}

func GetWithHeaders(url string, headers, data map[string]string) (body string, err error) {
	client := httpclient.WithHeaders(headers)
	var resp *httpclient.Response
	resp, err = client.Get(url, data)
	if err != nil {
		logger.Errorf("error %s", err)
		return
	}
	body, err = resp.ToString()
	return
}

func PostJSONWithTimeout(url string, headers map[string]string, timeout time.Duration, data interface{}) (body string, err error) {
	return PostJSONWithHeadersAndTimeout(url, map[string]string{}, timeout, data)
}

func PostJSONWithHeadersAndTimeout(url string, headers map[string]string, timeout time.Duration, data interface{}) (body string, err error) {
	client := httpclient.WithHeaders(headers)
	if timeout > 0 {
		client.WithOption(httpclient.OPT_TIMEOUT, timeout)
	}
	var resp *httpclient.Response
	resp, err = client.PostJson(url, data)
	if err != nil {
		logger.Errorf("error %s", err)
		return
	}
	body, err = resp.ToString()
	return
}
