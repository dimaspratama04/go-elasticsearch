package entity

type Movies struct {
	ID              int      `json:"id"`
	Title           string   `json:"title"`
	Year            int      `json:"year"`
	Cast            []string `json:"casts"`
	Genres          []string `json:"genres"`
	Href            string   `json:"href"`
	Extract         string   `json:"extract"`
	Thumbnail       string   `json:"thumbnail"`
	ThumbnailWidth  int      `json:"thumbnail_width"`
	ThumbnailHeight int      `json:"thumbnail_height"`
}
