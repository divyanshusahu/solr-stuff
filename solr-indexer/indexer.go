package main

import (
	"context"
	"solr-indexer/log"
	"solr-indexer/restclient/solr"
	"solr-indexer/restclient/tmdb"
	"strconv"
)

func indexer(ctx context.Context) {
	tmdbClient := tmdb.NewTmdbClient(ctx)
	solrClient := solr.NewSolrClient(ctx)
	totalPages := 500
	for page := 1; page <= totalPages; page++ {
		topRatedMoviesResponse, err := tmdbClient.GetTopRatedMovies(ctx, strconv.Itoa(page))
		if err != nil {
			log.Logger.Errorf("Error while fetching top rated movies for page:%v err:%v", page, err)
			continue
		}
		totalPages = topRatedMoviesResponse.TotalPages

		for i, result := range topRatedMoviesResponse.Results {
			log.Logger.Infof("indexing movie %v of page %v", i, page)

			doc := solr.TopRatedMoviesDoc{
				DocType:          "top_rated",
				PosterPath:       result.PosterPath,
				Adult:            result.Adult,
				Overview:         result.Overview,
				ReleaseDate:      result.ReleaseDate,
				GenreIds:         result.GenreIds,
				Id:               result.Id,
				OriginalTitle:    result.OriginalTitle,
				OriginalLanguage: result.OriginalLanguage,
				Title:            result.Title,
				BackdropPath:     result.BackdropPath,
				Popularity:       result.Popularity,
				VoteCount:        result.VoteCount,
				Video:            result.Video,
				VoteAverage:      result.VoteAverage,
			}

			err = solrClient.IndexDocument(ctx, doc)
			if err != nil {
				log.Logger.Errorf("error while indexing movie %v of page %v movie_id %v and movie_title %v", i, page, result.Id, result.Title)
				continue
			}

			log.Logger.Infof("successfully indexed movie %v of page %v movie_id %v and movie_title %v", i, page, result.Id, result.Title)
		}
	}
}
