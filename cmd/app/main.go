package main

import (
	"go-elasticsearch/internal/config"
	"go-elasticsearch/internal/delivery/http"
	"go-elasticsearch/internal/delivery/messaging"
	"go-elasticsearch/internal/repository"
	"go-elasticsearch/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New()

	// infrastructure
	elasticsearch := config.InitElasticSearch(cfg.ElasticURL, cfg.ElasticUsername, cfg.ElasticPassword)
	db := config.InitPostgresConnection(cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDBName, cfg.PostgresPort, cfg.PostgresSSLMode)
	rabbitmq, err := config.InitRabbitMQConnection(cfg.RabbitMQHost, cfg.RabbitMQPort, cfg.RabbitMQUsername, cfg.RabbitMQPassword, cfg.RabbitMQVhost)
	if err != nil {
		log.Info("[FATAL] Failed connect to rabbitmq:", err)
	}

	// Init Movies Messaging Delivery (Pub)
	moviesPublisher, err := messaging.NewMoviesPublisher(rabbitmq, "movies")
	if err != nil {
		log.Info("[FATAL] Failed to init publisher:", err)
	}

	// repository (postgres)
	PgMoviesRepository := repository.NewMoviesDBRepository(db)

	// repository (elasticsearch)
	ESMoviesRepository := repository.NewMoviesESRepository(elasticsearch)

	// usecase
	moviesUseCase := usecase.NewMoviesUseCase(PgMoviesRepository, ESMoviesRepository, moviesPublisher)

	// router
	http.InitializeRoute(&http.RouteConfig{
		App:           app,
		MoviesUseCase: moviesUseCase,
	})

	log.Info("✓ Starting server on port " + cfg.AppPort)
	log.Info("✓ Connected to Elasticsearch at " + cfg.ElasticURL)
	log.Info("✓ Connected to PostgreSQL at " + cfg.PostgresHost + ":" + cfg.PostgresPort)

	app.Listen(":" + cfg.AppPort)

}
