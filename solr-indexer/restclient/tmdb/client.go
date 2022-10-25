package tmdb

import (
	"context"
	"encoding/json"
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

func NewTmdbClient(ctx context.Context) *TmdbClient {
	rc := restclient.NewRestClient(ctx)

	apiKey := viper.GetString("TMDB_API_V3_KEY")
	authToken := fmt.Sprintf("Bearer %v", apiKey)
	rc.AddRequestHeader("Content-Type", "application/json;charset=utf-8")
	rc.AddRequestHeader("Authorization", authToken)

	tc := TmdbClient{client: rc}
	return &tc
}

func (tc *TmdbClient) GetTopRatedMovies(ctx context.Context, page string) (TopRatedMoviesResponse, error) {
	var topRatedMoviesResponse TopRatedMoviesResponse
	url, err := url.JoinPath(BASEURL, MoviesTopRatedPath)
	if err != nil {
		log.Logger.Errorw("error while joining url", "path", MoviesTopRatedPath, "err", err)
		return topRatedMoviesResponse, fmt.Errorf("error while joining url")
	}
	timeout := time.Duration(timeoutMap[MoviesTopRatedPath]) * time.Millisecond
	tc.client.Method = "GET"
	tc.client.Url = url
	tc.client.AddRequestParam("api_key", viper.GetString("TMDB_API_V3_KEY"))
	tc.client.AddRequestParam("page", page)
	tc.client.Timeout = timeout

	data, err := tc.client.FetchResponse(ctx)
	if err != nil {
		return topRatedMoviesResponse, err
	}

	err = json.Unmarshal(data, &topRatedMoviesResponse)
	if err != nil {
		log.Logger.Error("error while unmarshalling response - ", err)
	}

	return topRatedMoviesResponse, nil
}
