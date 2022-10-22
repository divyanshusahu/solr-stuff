package main

import (
	"context"
	"solr-indexer/log"
	"solr-indexer/restclient/solr"
	"solr-indexer/restclient/tmdb"
)

func indexer(ctx context.Context) {
	tmdbClient := tmdb.NewTmdbClient(ctx)
	topRatedMoviesResponse, err := tmdbClient.GetTopRatedMovies(ctx)
	log.Logger.Info(topRatedMoviesResponse.TotalResults, err)
	solrClient := solr.NewSolrClient(ctx)
	solrClient.IndexDocument(ctx)
}
