package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-elasticsearch/internal/entity"

	"github.com/elastic/go-elasticsearch/v9"
)

type MoviesRepository struct {
	ES *elasticsearch.Client
}

func NewMoviesRepository(es *elasticsearch.Client) *MoviesRepository {
	return &MoviesRepository{ES: es}
}

func (r *MoviesRepository) Insert(movie *entity.Movies) error {
	movieJSON, err := json.Marshal(movie)
	if err != nil {
		return err
	}

	_, err = r.ES.Index(
		"movies",
		bytes.NewReader(movieJSON),
		r.ES.Index.WithDocumentID(movie.ID),
		r.ES.Index.WithRefresh("true"),
	)

	return err
}

func (r *MoviesRepository) GetByID(id string) (*entity.Movies, error) {
	response, err := r.ES.Get("movies", id)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == 404 {
		return nil, fmt.Errorf("movie with ID %s not found", id)
	}

	var esResp struct {
		Source entity.Movies `json:"_source"`
	}

	if err := json.NewDecoder(response.Body).Decode(&esResp); err != nil {
		return nil, err
	}

	return &esResp.Source, nil

}
