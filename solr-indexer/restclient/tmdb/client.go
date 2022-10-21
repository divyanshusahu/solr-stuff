package tmdb

import (
	"context"
	"fmt"
	"net/url"
	"solr-indexer/restclient"
	"time"
)

type TmdbClient struct {
	client restclient.RestClient
}

func NewTmdbClient(ctx context.Context, path string) (*TmdbClient, error) {
	url, err := url.JoinPath(BASEURL, path)
	if err != nil {
		// TODO: add log here
		return nil, err
	}
	timeout := time.Duration(timeoutMap[path]) * time.Millisecond
	rc := restclient.NewRestClient(ctx, "GET", url, timeout)

	apiKey := ""
	authToken := fmt.Sprintf("Bearer %v", apiKey)
	rc.AddRequestHeader("Content-Type", "application/json;charset=utf-8")
	rc.AddRequestHeader("Authorization", authToken)

	tc := TmdbClient{client: rc}
	return &tc, nil
}
