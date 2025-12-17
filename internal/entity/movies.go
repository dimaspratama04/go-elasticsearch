package entity

import (
	"time"

	"github.com/lib/pq"
)

type Movies struct {
	ID              uint           `json:"id"`
	Title           string         `json:"title"`
	Year            int            `json:"year"`
	Casts           pq.StringArray `json:"casts"`
	Genres          pq.StringArray `json:"genres"`
	Href            string         `json:"href"`
	Extract         string         `json:"extract"`
	Thumbnail       string         `json:"thumbnail"`
	ThumbnailWidth  int            `json:"thumbnail_width"`
	ThumbnailHeight int            `json:"thumbnail_height"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
