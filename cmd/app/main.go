package main

import (
	"go-elasticsearch/internal/config"
	"go-elasticsearch/internal/delivery/http"
	"go-elasticsearch/internal/delivery/http/usecase"
	"go-elasticsearch/internal/delivery/messaging/publisher"
	"go-elasticsearch/internal/repository/elasticsearchdb"
	postgresdb "go-elasticsearch/internal/repository/postgresdb"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New()

	// infrastructure
	elasticsearch := config.InitElasticSearch(cfg.ElasticURL, cfg.ElasticUsername, cfg.ElasticPassword)
	db := config.InitPostgresConnection(cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDBName, cfg.PostgresPort, cfg.PostgresSSLMode)
	rabbitmq := config.InitRabbitMQConnection(cfg.RabbitMQHost, cfg.RabbitMQPort, cfg.RabbitMQUsername, cfg.RabbitMQPassword, cfg.RabbitMQVhost)

	// Init Movies Messaging Delivery (Pub)
	moviesPublisher, _ := publisher.NewMoviesPublisher(rabbitmq, "movies")

	// repository (postgres)
	PgMoviesRepository := postgresdb.NewMoviesDBRepository(db)

	// repository (elasticsearch)
	ESMoviesRepository := elasticsearchdb.NewMoviesESRepository(elasticsearch)

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
