package main

import (
	"context"
	"solr-indexer/log"

	"github.com/spf13/viper"
)

func main() {
	log.InitLogger()
	defer log.Logger.Sync()

	err := initConfig()
	if err != nil {
		log.Logger.Errorf("Failed to read config file - ", err.Error())
		return
	}

	ctx := context.Background()
	indexer(ctx)
}

func initConfig() error {
	viper.SetConfigName("secrets")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
