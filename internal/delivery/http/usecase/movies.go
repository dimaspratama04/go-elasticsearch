package usecase

import (
	"go-elasticsearch/internal/entity"
	"go-elasticsearch/internal/repository"
)

type MoviesUseCase struct {
	Repository *repository.MoviesRepository
}

func NewMoviesUseCase(repository *repository.MoviesRepository) *MoviesUseCase {
	return &MoviesUseCase{Repository: repository}
}

func (uc *MoviesUseCase) InsertMovies(movies *entity.Movies) error {
	return uc.Repository.Insert(movies)
}

func (uc *MoviesUseCase) GetMovieByID(id string) (*entity.Movies, error) {
	return uc.Repository.GetByID(id)
}
