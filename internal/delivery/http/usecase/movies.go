package usecase

import (
	"go-elasticsearch/internal/delivery/messaging"
	"go-elasticsearch/internal/entity"
	"go-elasticsearch/internal/model"
	"go-elasticsearch/internal/repository"
	"log"
)

type MoviesUseCase struct {
	pgRepository *repository.MoviesDBRepository
	esRepository *repository.MoviesESRepository
	publisher    *messaging.MoviesPublisher
}

type MoviesIndexUseCase struct {
	esRepository *repository.MoviesESRepository
}

func NewMoviesIndexUseCase(es *repository.MoviesESRepository) *MoviesIndexUseCase {
	return &MoviesIndexUseCase{esRepository: es}
}

func NewMoviesUseCase(pg *repository.MoviesDBRepository, es *repository.MoviesESRepository, pub *messaging.MoviesPublisher) *MoviesUseCase {
	return &MoviesUseCase{pgRepository: pg, esRepository: es, publisher: pub}
}

func (uc *MoviesUseCase) InsertMovies(movies *model.Movies) error {
	// 1. insert ke Postgres
	if err := uc.pgRepository.Create(movies); err != nil {
		log.Println("[ERROR] Failed insert movie to Postgres:", err)
		return err
	}

	// 1. pub ke msgbroker
	if err := uc.publisher.Publish("movies.created", movies); err != nil {
		log.Println("[ERROR] Failed publish movie event:", err)
		return err
	}

	return nil
}

func (uc *MoviesUseCase) BulkInsertMovies(movies []model.Movies) error {
	// insert ke postgress
	if err := uc.pgRepository.BulkInsert(movies); err != nil {
		log.Println("[ERROR] Failed bulk insert to postgres", err)
	}

	// pass ke msgbroker
	for _, m := range movies {
		if err := uc.publisher.Publish("movies.created", &m); err != nil {
			log.Printf(
				"[ERROR] Failed publish movie id=%d title=%s : %v\n",
				m.ID,
				m.Title,
				err,
			)
		}
	}
	return nil
}

func (uc *MoviesUseCase) SearchMovies(query string) ([]entity.Movies, error) {
	return uc.esRepository.Search(query)
}

func (u *MoviesIndexUseCase) IndexMovies(movies *model.Movies) error {
	return u.esRepository.Index(movies)
}
