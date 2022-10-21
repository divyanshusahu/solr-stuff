package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
)

type HttpClient struct {
	Method         string
	Url            string
	Timeout        time.Duration
	requestHeaders map[string]string
	requestCookies map[string]string
	requestBody    *bytes.Reader
}

func NewHttpClient(ctx context.Context, method string, url string, timeout time.Duration) HttpClient {
	if method != "POST" {
		method = "GET"
	}

	hc := HttpClient{
		Method:  method,
		Url:     url,
		Timeout: timeout,
	}

	return hc
}

func (hc *HttpClient) AddRequestHeader(header string, value string) {
	if hc.requestHeaders == nil {
		hc.requestHeaders = map[string]string{}
	}
	hc.requestHeaders[header] = value
}

func (hc *HttpClient) AddRequestCookie(cookie string, value string) {
	if hc.requestCookies == nil {
		hc.requestCookies = map[string]string{}
	}
	hc.requestCookies[cookie] = value
}

func (hc *HttpClient) AddRequestBody(body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		// TODO: add a log here
		data = nil
	}
	hc.requestBody = bytes.NewReader(data)
}

func (hc *HttpClient) FetchResponse(ctx context.Context) ([]byte, error) {
	return hc.getResponse(ctx)
}

func (hc *HttpClient) getResponse(ctx context.Context) ([]byte, error) {
	if hc.Timeout == 0 {
		hc.Timeout = 1000 * time.Millisecond
	}
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(hc.Timeout))

	var req *http.Request
	var err error

	switch hc.Method {
	case "GET":
		req, err = http.NewRequest(hc.Method, hc.Url, nil)
	case "POST":
		req, err = http.NewRequest(hc.Method, hc.Url, hc.requestBody)
	default:
		req, err = http.NewRequest(hc.Method, hc.Url, nil)
	}

	if err != nil {
		// TODO: add a log here
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		// TODO: add a log here
		return nil, err
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		// TODO: add a log here
		return nil, err
	}

	return response, nil
}
