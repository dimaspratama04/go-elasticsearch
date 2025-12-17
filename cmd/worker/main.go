package main

import (
	"go-elasticsearch/internal/config"
	"go-elasticsearch/internal/delivery/messaging"
	"go-elasticsearch/internal/repository"

	"github.com/gofiber/fiber/v2/log"
)

func main() {
	cfg := config.LoadConfig()
	// global vars
	exchangeName := "movies"
	routingKeyName := "movies.created"

	// infrastrcucture
	rabbitmq, err := config.InitRabbitMQConnection(cfg.RabbitMQHost, cfg.RabbitMQPort, cfg.RabbitMQUsername, cfg.RabbitMQPassword, cfg.RabbitMQVhost)
	if err != nil {
		log.Info("[FATAL] Failed connect to rabbitmq:", err)
	}

	elasticsearch := config.InitElasticSearch(cfg.ElasticURL, cfg.ElasticUsername, cfg.ElasticPassword)

	// ES Repository
	ESMoviesRepository := repository.NewMoviesESRepository(elasticsearch)

	// Init worker
	worker, _ := messaging.NewMoviesConsumer(rabbitmq, exchangeName, ESMoviesRepository)
	if err := worker.Consumer(routingKeyName); err != nil {
		log.Info("[WORKER ERROR] failed connect to broker")
	}

}
