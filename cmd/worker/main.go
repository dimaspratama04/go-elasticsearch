package main

import (
	"go-elasticsearch/internal/config"
	"go-elasticsearch/internal/delivery/messaging"
	"go-elasticsearch/internal/repository"
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
	ESMoviesRepository := repository.NewMoviesESRepository(elasticsearch)

	// Init worker
	worker, _ := messaging.NewMoviesConsumer(rabbitmq, exchangeName, ESMoviesRepository)
	if err := worker.Consumer(routingKeyName); err != nil {
		log.Println("[WORKER ERROR] failed connect to broker")
	}

}
