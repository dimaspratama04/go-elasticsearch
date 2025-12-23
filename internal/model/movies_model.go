package model

import (
	"time"
)

type Movies struct {
	ID              uint      `gorm:"column:id;primaryKey"`
	Title           string    `gorm:"column:title"`
	Year            int       `gorm:"column:year"`
	Casts           []string  `gorm:"column:casts"`
	Genres          []string  `gorm:"column:genres[]"`
	Href            string    `gorm:"column:href"`
	Extract         string    `gorm:"column:extract"`
	Thumbnail       string    `gorm:"column:thumbnail"`
	ThumbnailWidth  int       `gorm:"column:thumbnail_width"`
	ThumbnailHeight int       `gorm:"column:thumbnail_height"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

type CreateMovieRequest struct {
	Title           string   `json:"title" validate:"required"`
	Year            int      `json:"year"`
	Casts           []string `json:"cast"`
	Genres          []string `json:"genres"`
	Href            string   `json:"href"`
	Extract         string   `json:"extract"`
	Thumbnail       string   `json:"thumbnail"`
	ThumbnailWidth  int      `json:"thumbnail_width"`
	ThumbnailHeight int      `json:"thumbnail_height"`
}

type MovieResponse struct {
	ID              uint     `json:"id"`
	Title           string   `json:"title"`
	Year            int      `json:"year"`
	Casts           []string `json:"casts"`
	Genres          []string `json:"genres"`
	Href            string   `json:"href"`
	Extract         string   `json:"extract"`
	Thumbnail       string   `json:"thumbnail"`
	ThumbnailWidth  int      `json:"thumbnail_width"`
	ThumbnailHeight int      `json:"thumbnail_height"`
}
