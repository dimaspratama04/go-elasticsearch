package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-elasticsearch/internal/entity"
	"strconv"

	"github.com/elastic/go-elasticsearch/v9"
)

type MoviesRepository struct {
	ES *elasticsearch.Client
}

func NewMoviesRepository(es *elasticsearch.Client) *MoviesRepository {
	return &MoviesRepository{ES: es}
}

func (r *MoviesRepository) BulkInsert(movies []entity.Movies) error {
	var buf bytes.Buffer

	for _, movie := range movies {

		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "movies" } }%s`, "\n"))

		data, err := json.Marshal(movie)
		if err != nil {
			return err
		}
		data = append(data, byte('\n'))

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	resp, err := r.ES.Bulk(bytes.NewReader(buf.Bytes()), r.ES.Bulk.WithIndex("movies"))

	fmt.Println("debug", resp)
	return err
}

func (r *MoviesRepository) Insert(movie *entity.Movies) error {
	movieJSON, err := json.Marshal(movie)
	var movieID = strconv.Itoa(movie.ID)

	if err != nil {
		return err
	}

	_, err = r.ES.Index(
		"movies",
		bytes.NewReader(movieJSON),
		r.ES.Index.WithDocumentID(movieID),
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

func (r *MoviesRepository) GetAllAutoCompleteSuggestions(field, indexName, query string) ([]string, error) {
	bodyQuery := map[string]interface{}{
		"suggest": map[string]interface{}{
			"movie-suggest": map[string]interface{}{
				"prefix": query,
				"completion": map[string]interface{}{
					"field": field, // title_autocomplete, director_autocomplete, etc.
					"fuzzy": true,
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(bodyQuery); err != nil {
		return nil, err
	}

	// Perform the search request.
	resp, err := r.ES.Search(
		r.ES.Search.WithIndex(indexName),
		r.ES.Search.WithBody(&buf),
	)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var esResp struct {
		Suggest map[string][]struct {
			Options []struct {
				Text string `json:"text"`
			} `json:"options"`
		} `json:"suggest"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&esResp); err != nil {
		return nil, err
	}

	suggestions := []string{}
	for _, option := range esResp.Suggest["movie-suggest"][0].Options {
		suggestions = append(suggestions, option.Text)
	}

	return suggestions, nil
}
