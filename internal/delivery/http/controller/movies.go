package controller

import (
	"encoding/json"
	"fmt"
	"go-elasticsearch/internal/delivery/http/usecase"
	"go-elasticsearch/internal/model"
	"net/http"

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

func (mc *MoviesController) SearchMovies(ctx *fiber.Ctx) error {
	// Implementation for searching movies
	query := ctx.Query("q", "")

	fmt.Println(query)

	movies, err := mc.Usecase.SearchMovies(query)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": movies,
	})
}

func (mc *MoviesController) CreateMovies(ctx *fiber.Ctx) error {
	// Implementation for creating a new movie
	var movies model.Movies

	if err := ctx.BodyParser(&movies); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
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

func (mc *MoviesController) BulkInsertMovies(ctx *fiber.Ctx) error {
	var movies []model.Movies

	url := "https://raw.githubusercontent.com/prust/wikipedia-movie-data/refs/heads/master/movies-1980s.json"

	resp, err := http.Get(url)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to fetch movies data",
		})
	}

	if err := json.NewDecoder(resp.Body).Decode(&movies); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to decode movies data",
		})
	}

	if err := mc.Usecase.BulkInsertMovies(movies); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to insert movies",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "movies created successfully",
		"count":   len(movies),
	})
}
