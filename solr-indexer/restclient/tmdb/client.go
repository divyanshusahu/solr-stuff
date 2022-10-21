package tmdb

import (
	"context"
	"fmt"
	"net/url"
	"solr-indexer/log"
	"solr-indexer/restclient"
	"time"

	"github.com/spf13/viper"
)

type TmdbClient struct {
	client restclient.RestClient
}

func NewTmdbClient(ctx context.Context) (*TmdbClient, error) {
	rc := restclient.NewRestClient(ctx)

	apiKey := viper.GetString("TMDB_API_V3_KEY")
	authToken := fmt.Sprintf("Bearer %v", apiKey)
	rc.AddRequestHeader("Content-Type", "application/json;charset=utf-8")
	rc.AddRequestHeader("Authorization", authToken)

	tc := TmdbClient{client: rc}
	return &tc, nil
}

func (tc *TmdbClient) GetTopRatedMovies(ctx context.Context) {
	url, err := url.JoinPath(BASEURL, MoviesTopRatedPath)
	if err != nil {
		log.Logger.Errorw("error while joining url", "path", MoviesTopRatedPath, "err", err)
		return
	}
	timeout := time.Duration(timeoutMap[MoviesTopRatedPath]) * time.Millisecond
	tc.client.Method = "GET"
	tc.client.Url = url
	tc.client.Timeout = timeout

	response, err := tc.client.FetchResponse(ctx)
	log.Logger.Info(response, err)
}
