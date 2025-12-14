package usecase

import (
	"go-elasticsearch/internal/delivery/messaging/publisher"
	"go-elasticsearch/internal/entity"
	"go-elasticsearch/internal/model"
	"go-elasticsearch/internal/repository/elasticsearchdb"
	"go-elasticsearch/internal/repository/postgresdb"
	"log"
)

type MoviesUseCase struct {
	pgRepository *postgresdb.MoviesDBRepository
	esRepository *elasticsearchdb.MoviesESRepository
	publisher    *publisher.MoviesPublisher
}

type MoviesIndexUseCase struct {
	esRepository *elasticsearchdb.MoviesESRepository
}

func NewMoviesIndexUseCase(es *elasticsearchdb.MoviesESRepository) *MoviesIndexUseCase {
	return &MoviesIndexUseCase{esRepository: es}
}

func NewMoviesUseCase(pg *postgresdb.MoviesDBRepository, es *elasticsearchdb.MoviesESRepository, pub *publisher.MoviesPublisher) *MoviesUseCase {
	return &MoviesUseCase{pgRepository: pg, esRepository: es, publisher: pub}
}

func (uc *MoviesUseCase) InsertMovies(movies *model.Movies) error {
	// 1. insert ke Postgres
	if err := uc.pgRepository.Insert(movies); err != nil {
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

// todo add bulk insert logic
// func (uc *MoviesUseCase) BulkInsertMovies(movies []model.Movies) error {
// 	// 1. insert ke Postgres
// 	if err := uc.pgRepository.BulkInsert(movies); err != nil {
// 		log.Println("[ERROR] Failed insert movie to Postgres:", err)
// 		return err
// 	}

// 	event := struct {
// 		Event string         `json:"event"`
// 		Data  []model.Movies `json:"data"`
// 	}{
// 		Event: "movies.created",
// 		Data:  movies,
// 	}

// 	// 2. pub ke msgbroker
// 	if err := uc.publisher.Publish("movies.created", event); err != nil {
// 		log.Println("[ERROR] Failed publish movie bulk event:", err)
// 		return err
// 	}

// 	log.Printf("[INFO] Bulk movies inserted & published (%d items)", len(movies))

// 	return nil
// }

func (uc *MoviesUseCase) SearchMovies(query string) ([]entity.Movies, error) {
	return uc.esRepository.Search(query)
}

func (u *MoviesIndexUseCase) IndexMovies(movies *model.Movies) error {
	return u.esRepository.Index(movies)
}
