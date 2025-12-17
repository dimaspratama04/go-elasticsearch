package controller

import (
	"context"
	"fmt"
	"go-elasticsearch/internal/model"
	"go-elasticsearch/internal/usecase"

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

func (mc *MoviesController) SearchMovies(c *fiber.Ctx) error {
	// Implementation for searching movies
	query := c.Query("q", "")

	fmt.Println(query)

	movies, err := mc.Usecase.SearchMovies(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": movies,
	})
}

func (mc *MoviesController) CreateMovies(c *fiber.Ctx) error {
	// Implementation for creating a new movie
	var movies *model.CreateMovieRequest

	if err := c.BodyParser(&movies); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	ctx := context.Background()
	if err := mc.Usecase.InsertMovies(ctx, movies); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to insert movies",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "movies created successfully",
	})

}

func (mc *MoviesController) BulkInsertMovies(ctx *fiber.Ctx) error {
	var movies []model.Movies

	// url := "https://raw.githubusercontent.com/prust/wikipedia-movie-data/refs/heads/master/movies-1980s.json"

	// resp, err := http.Get(url)

	if err := ctx.BodyParser(&movies); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// if err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "failed to fetch movies data",
	// 	})
	// }

	// if err := json.NewDecoder(resp.Body).Decode(&movies); err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "failed to decode movies data",
	// 	})
	// }

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
