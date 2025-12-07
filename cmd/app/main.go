package main

import (
	"go-elasticsearch/internal/config"
	"go-elasticsearch/internal/delivery/http"
	"go-elasticsearch/internal/delivery/http/usecase"
	"go-elasticsearch/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New()

	elasticsearch := config.InitElasticSearch(cfg.ElasticURL, cfg.ElasticUsername, cfg.ElasticPassword)

	// repository
	moviesRepository := repository.NewMoviesRepository(elasticsearch)

	// usecase
	moviesUseCase := usecase.NewMoviesUseCase(moviesRepository)

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
