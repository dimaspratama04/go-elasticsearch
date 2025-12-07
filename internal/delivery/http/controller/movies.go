package controller

import (
	"go-elasticsearch/internal/delivery/http/usecase"
	"go-elasticsearch/internal/entity"

	"github.com/gofiber/fiber/v2"
)

type MoviesController struct {
	Usecase *usecase.MoviesUseCase
}

func NewMoviesController(usecase *usecase.MoviesUseCase) *MoviesController {
	return &MoviesController{
		Usecase: usecase,
	}
}

func (mc *MoviesController) GetMovies() {
	// Implementation for getting movies
}

func (mc *MoviesController) GetMovieByID(id string) {
	// Implementation for getting a movie by ID
}

func (mc *MoviesController) InsertMovies(ctx *fiber.Ctx) error {
	// Implementation for creating a new movie
	var movies entity.Movies

	if err := ctx.BodyParser(&movies); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := mc.Usecase.InsertMovies(&movies); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to insert movies",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "movies created successfully",
	})

}
