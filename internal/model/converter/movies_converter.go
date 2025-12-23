package converter

import (
	"go-elasticsearch/internal/entity"
	"go-elasticsearch/internal/model"

	"github.com/lib/pq"
)

func requestToEntity(req model.CreateMovieRequest) entity.Movies {
	return entity.Movies{
		Title: req.Title,
		Year:  req.Year,
		// Konversi slice biasa ke pq.StringArray
		Casts:           pq.StringArray(req.Casts),
		Genres:          pq.StringArray(req.Genres),
		Href:            req.Href,
		Extract:         req.Extract,
		Thumbnail:       req.Thumbnail,
		ThumbnailWidth:  req.ThumbnailWidth,
		ThumbnailHeight: req.ThumbnailHeight,
	}
}

func EntityToModel(e *entity.Movies) *model.Movies {
	return &model.Movies{
		ID:              e.ID,
		Title:           e.Title,
		Year:            e.Year,
		Casts:           pq.StringArray(e.Casts),
		Genres:          pq.StringArray(e.Genres),
		Href:            e.Href,
		Extract:         e.Extract,
		Thumbnail:       e.Thumbnail,
		ThumbnailWidth:  e.ThumbnailWidth,
		ThumbnailHeight: e.ThumbnailHeight,
	}
}

func ModelToEntity(m *model.Movies) *entity.Movies {
	return &entity.Movies{
		ID:              m.ID,
		Title:           m.Title,
		Year:            m.Year,
		Casts:           []string(m.Casts),
		Genres:          []string(m.Genres),
		Href:            m.Href,
		Extract:         m.Extract,
		Thumbnail:       m.Thumbnail,
		ThumbnailWidth:  m.ThumbnailWidth,
		ThumbnailHeight: m.ThumbnailHeight,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func RequestToEntity(req *model.CreateMovieRequest) *entity.Movies {
	return &entity.Movies{
		Title:           req.Title,
		Casts:           pq.StringArray(req.Casts),
		Genres:          pq.StringArray(req.Genres),
		Href:            req.Href,
		Extract:         req.Extract,
		Thumbnail:       req.Thumbnail,
		ThumbnailWidth:  req.ThumbnailWidth,
		ThumbnailHeight: req.ThumbnailHeight,
	}
}

func RequestToEntities(reqs []model.CreateMovieRequest) []entity.Movies {
	movies := make([]entity.Movies, 0, len(reqs))

	for _, req := range reqs {
		movies = append(movies, requestToEntity(req))
	}

	return movies
}

func MoviesToResponse(movies *entity.Movies) *model.Movies {
	return &model.Movies{
		ID:              movies.ID,
		Title:           movies.Title,
		Year:            movies.Year,
		Casts:           movies.Genres,
		Href:            movies.Href,
		Extract:         movies.Extract,
		Thumbnail:       movies.Thumbnail,
		ThumbnailWidth:  movies.ThumbnailWidth,
		ThumbnailHeight: movies.ThumbnailHeight,
	}
}
