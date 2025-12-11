package usecase

import (
	"go-elasticsearch/internal/entity"
	models "go-elasticsearch/internal/model"
	"go-elasticsearch/internal/repository"
)

type MoviesUseCase struct {
	Repository *repository.MoviesRepository
}

func NewMoviesUseCase(repository *repository.MoviesRepository) *MoviesUseCase {
	return &MoviesUseCase{Repository: repository}
}

func (uc *MoviesUseCase) BulkInsertMovies(movies []models.Movies) error {
	return uc.Repository.BulkInsert(movies)
}

func (uc *MoviesUseCase) InsertMovies(movies *models.Movies) error {
	return uc.Repository.Insert(movies)
}

func (uc *MoviesUseCase) GetMovieByID(id string) (*entity.Movies, error) {
	return uc.Repository.GetByID(id)
}

func (uc *MoviesUseCase) SearchMovies(query string) ([]entity.Movies, error) {
	return uc.Repository.Search(query)
}
