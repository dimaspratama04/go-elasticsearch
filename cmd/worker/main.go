package main

import (
	"go-elasticsearch/internal/config"
	"go-elasticsearch/internal/delivery/http/usecase"
	"go-elasticsearch/internal/delivery/messaging/consumer"
	"go-elasticsearch/internal/repository/elasticsearchdb"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	// global vars
	exchangeName := "movies"
	routingKeyName := "movies.created"

	// infrastrcucture
	rabbitmq := config.InitRabbitMQConnection(cfg.RabbitMQHost, cfg.RabbitMQPort, cfg.RabbitMQUsername, cfg.RabbitMQPassword, cfg.RabbitMQVhost)
	elasticsearch := config.InitElasticSearch(cfg.ElasticURL, cfg.ElasticUsername, cfg.ElasticPassword)

	// ES Repository
	ESMoviesRepository := elasticsearchdb.NewMoviesESRepository(elasticsearch)

	// Index Usecase
	usecase := usecase.NewMoviesIndexUseCase(ESMoviesRepository)

	// Init worker
	worker, err := consumer.NewMoviesConsumerWorker(rabbitmq, usecase, exchangeName)

	if err != nil {
		log.Println("[WORKER ERROR] failed connect to broker")
	}

	worker.Consumer(exchangeName, routingKeyName)
}
