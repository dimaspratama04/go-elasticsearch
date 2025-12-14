package postgresdb

import (
	"go-elasticsearch/internal/model"
	"log"

	"gorm.io/gorm"
)

type MoviesDBRepository struct {
	db *gorm.DB
}

func NewMoviesDBRepository(db *gorm.DB) *MoviesDBRepository {
	return &MoviesDBRepository{db: db}
}

func (r *MoviesDBRepository) Insert(movies *model.Movies) error {
	err := r.db.Create(movies).Error
	log.Println("[SUCESS] insert movie to Postgres")

	if err != nil {
		log.Println("[ERROR] Failed insert movie to Postgres:", err)
	}

	return nil
}

func (r *MoviesDBRepository) BulkInsert(movies []model.Movies) error {
	if err := r.db.CreateInBatches(movies, 1000).Error; err != nil {
		return err
	}

	return nil
}
