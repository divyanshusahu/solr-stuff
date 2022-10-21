package main

import (
	"fmt"
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

	fmt.Println(viper.GetString("TMDB_API_V3_KEY"))
}

func initConfig() error {
	viper.SetConfigName("secrets")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
