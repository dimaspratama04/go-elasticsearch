package config

import (
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/gofiber/fiber/v2/log"
)

func InitElasticSearch(url, username, password string) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
		Username:  username,
		Password:  password,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: " + err.Error())
	}

	_, err = es.Info()
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: " + err.Error())
	}

	return es
}
