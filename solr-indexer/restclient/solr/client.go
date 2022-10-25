package solr

import (
	"context"
	"net/url"
	"solr-indexer/log"
	"solr-indexer/restclient"
)

const BASEURL = "http://localhost:50000/solr"
const COLLECTION = "moviesdb"

type SolrClient struct {
	client restclient.RestClient
}

func NewSolrClient(ctx context.Context) SolrClient {
	rc := restclient.NewRestClient(ctx)
	rc.AddRequestHeader("Content-Type", "application/json")

	return SolrClient{client: rc}
}

func (sc *SolrClient) IndexDocument(ctx context.Context, doc interface{}) error {
	path := "update/json/docs"
	url, err := url.JoinPath(BASEURL, COLLECTION, path)
	if err != nil {
		log.Logger.Errorw("error while joining url", "baseUrl", BASEURL, "err", err)
	}
	sc.client.Method = "POST"
	sc.client.Url = url
	sc.client.AddRequestBody(doc)

	_, err = sc.client.FetchResponse(ctx)
	return err
}
