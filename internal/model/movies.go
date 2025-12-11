package model

import (
	"time"

	"github.com/lib/pq"
)

type Movies struct {
	ID              uint   `gorm:"primaryKey"`
	Title           string `gorm:"size:255;not null"`
	Year            int
	Casts           pq.StringArray `gorm:"type:text[]"`
	Genres          pq.StringArray `gorm:"type:text[]"`
	Href            string
	Extract         string `gorm:"type:text"`
	Thumbnail       string
	ThumbnailWidth  int
	ThumbnailHeight int

	CreatedAt time.Time
	UpdatedAt time.Time
}
