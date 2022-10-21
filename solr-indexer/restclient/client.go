package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
)

type RestClient struct {
	Method         string
	Url            string
	Timeout        time.Duration
	requestHeaders map[string]string
	requestCookies map[string]string
	requestBody    *bytes.Reader
}

func NewRestClient(ctx context.Context, method string, url string, timeout time.Duration) RestClient {
	if method != "POST" {
		method = "GET"
	}

	hc := RestClient{
		Method:  method,
		Url:     url,
		Timeout: timeout,
	}

	return hc
}

func (rc *RestClient) AddRequestHeader(header string, value string) {
	if rc.requestHeaders == nil {
		rc.requestHeaders = map[string]string{}
	}
	rc.requestHeaders[header] = value
}

func (rc *RestClient) AddRequestCookie(cookie string, value string) {
	if rc.requestCookies == nil {
		rc.requestCookies = map[string]string{}
	}
	rc.requestCookies[cookie] = value
}

func (rc *RestClient) AddRequestBody(body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		// TODO: add a log here
		data = nil
	}
	rc.requestBody = bytes.NewReader(data)
}

func (rc *RestClient) FetchResponse(ctx context.Context) ([]byte, error) {
	return rc.getResponse(ctx)
}

func (rc *RestClient) getResponse(ctx context.Context) ([]byte, error) {
	if rc.Timeout == 0 {
		rc.Timeout = 1000 * time.Millisecond
	}
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(rc.Timeout))

	var req *http.Request
	var err error

	switch rc.Method {
	case "GET":
		req, err = http.NewRequest(rc.Method, rc.Url, nil)
	case "POST":
		req, err = http.NewRequest(rc.Method, rc.Url, rc.requestBody)
	default:
		req, err = http.NewRequest(rc.Method, rc.Url, nil)
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
