package http

import (
	"go-elasticsearch/internal/delivery/http/controller"
	"go-elasticsearch/internal/delivery/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App           *fiber.App
	MoviesUseCase *usecase.MoviesUseCase
}

func InitializeRoute(rc *RouteConfig) {
	moviesController := controller.NewMoviesController(rc.MoviesUseCase)

	rc.App.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "ok, api is works!",
		})
	})

	rc.App.Get("/api/v1/movies", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "oops, endpoint not implemented yet",
		})
	})

	rc.App.Post("/api/v1/movies", moviesController.InsertMovies)
}
