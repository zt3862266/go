package util

import (
	. "github.com/zt3862266/go/log"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewHttpClient(maxIdleConns, maxIdleConnsPerHost, idleConnTimeout int) *http.Client {

	transport := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		IdleConnTimeout:     time.Duration(idleConnTimeout) * time.Second,
		Dial: (&net.Dialer{
			Timeout:   1 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := http.Client{
		Transport: transport,
		Timeout:   1 * time.Second,
	}
	return &client
}

func SendPost(client *http.Client, postUrl string, param url.Values, header map[string]string) (ret []byte, err error) {
	req, _ := http.NewRequest("POST", postUrl, strings.NewReader(param.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded ")
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	response, err := client.Do(req)

	if err != nil {
		Warn("request failed,url:%s ,err:%s", postUrl, err)
		return nil, err
	}
	defer response.Body.Close()
	statusCode := response.StatusCode
	retStr, err := ioutil.ReadAll(response.Body)
	Info("send post,url:%s,msg:%s,ret:%s", postUrl, param, retStr)
	if err != nil {
		Warn("get response failed,url:%s,err:%s,statusCode:%d", postUrl, err, statusCode)
		return nil, err
	}
	return retStr, nil
}
