package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"solr-indexer/log"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
)

type RestClient struct {
	Method         string
	Url            string
	Timeout        time.Duration
	requestHeaders map[string]string
	requestParams  map[string]string
	requestBody    interface{}
}

func NewRestClient(ctx context.Context) RestClient {
	hc := RestClient{}

	return hc
}

func (rc *RestClient) AddRequestHeader(header string, value string) {
	if rc.requestHeaders == nil {
		rc.requestHeaders = map[string]string{}
	}
	rc.requestHeaders[header] = value
}

func (rc *RestClient) AddRequestParam(key string, value string) {
	if rc.requestParams == nil {
		rc.requestParams = map[string]string{}
	}
	rc.requestParams[key] = value
}

func (rc *RestClient) AddRequestBody(body interface{}) {
	rc.requestBody = body
}

func (rc *RestClient) FetchResponse(ctx context.Context) ([]byte, error) {
	return rc.getResponse(ctx)
}

func (rc *RestClient) getResponse(ctx context.Context) ([]byte, error) {
	if rc.Url == "" {
		log.Logger.Error("no url supplied to the rest client")
		return nil, fmt.Errorf("no url supplied to the rest client")
	}

	if rc.Timeout == 0 {
		rc.Timeout = 1000 * time.Millisecond
	}
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(rc.Timeout))

	var req *http.Request
	var err error

	switch rc.Method {
	case "GET":
		req, err = http.NewRequest(rc.Method, rc.Url, nil)
		if req != nil && err == nil {
			rc.addQueryParams(req)
		}
	case "POST":
		req, err = http.NewRequest(rc.Method, rc.Url, rc.addPostBody())
	default:
		req, err = http.NewRequest(rc.Method, rc.Url, nil)
	}
	if err != nil {
		log.Logger.Error("error creating http request - ", err)
		return nil, err
	}

	for k, v := range rc.requestHeaders {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Logger.Errorw("error in http call", "rc", rc, "err", err)
		return nil, err
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		log.Logger.Errorw("error while reading response", "rc", rc, "err", err)
		return nil, err
	}

	if res.StatusCode >= 400 {
		log.Logger.Warnf("non success response code returned from the server. StatusCode: %v", res.StatusCode)
		log.Logger.Infof("rc: %v and response: %v", rc, string(response))
	}

	return response, nil
}

func (rc *RestClient) addQueryParams(req *http.Request) {
	q := req.URL.Query()
	for k, v := range rc.requestParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func (rc *RestClient) addPostBody() *bytes.Buffer {
	var data []byte
	var err error
	if rc.requestBody != nil {
		data, err = json.Marshal(rc.requestBody)
	} else {
		data, err = json.Marshal(rc.requestParams)
	}
	if err != nil {
		log.Logger.Error("error while marshalling post body - ", err)
	}

	return bytes.NewBuffer(data)
}
