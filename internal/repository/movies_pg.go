package repository

import (
	"go-elasticsearch/internal/entity"
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

func (r *MoviesDBRepository) Create(movies *entity.Movies) error {
	if err := r.db.Table("movies").Create(movies).Error; err != nil {
		log.Println("[ERROR] Failed insert movie to Postgres:", err)
		return err
	}

	return nil
}

func (r *MoviesDBRepository) BulkInsert(movies []model.Movies) error {
	if err := r.db.CreateInBatches(movies, 1000).Error; err != nil {
		log.Println("[ERROR] Failed Bulk insert movie to Postgres:", err)
	}

	return nil
}
