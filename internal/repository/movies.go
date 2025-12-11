package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-elasticsearch/internal/delivery/messaging"
	"go-elasticsearch/internal/entity"
	"go-elasticsearch/internal/helper"
	"go-elasticsearch/internal/model"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type MoviesRepository struct {
	ES              *elasticsearch.Client
	DB              *gorm.DB
	RabbitMQ        *amqp091.Connection
	MoviesPublisher *messaging.MoviesPublisher
}

func NewMoviesRepository(db *gorm.DB, es *elasticsearch.Client, rabbitmq *amqp091.Connection, moviesPublisher *messaging.MoviesPublisher) *MoviesRepository {
	return &MoviesRepository{ES: es, DB: db, RabbitMQ: rabbitmq, MoviesPublisher: moviesPublisher}
}

func (r *MoviesRepository) Insert(movie *model.Movies) error {
	// 1) Insert ke PostgreSQL
	if err := r.DB.Create(movie).Error; err != nil {
		return err
	}

	// 2) Publish ke RabbitMQ
	if err := r.MoviesPublisher.Publish(movie); err != nil {
		return err
	}

	return nil

	// return r.DB.Create(movie).Error
}

func (r *MoviesRepository) BulkInsert(movies []model.Movies) error {
	return r.DB.CreateInBatches(movies, 1000).Error
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

func (r *MoviesRepository) Search(query string) ([]entity.Movies, error) {
	var buf bytes.Buffer

	searchQuery := helper.ESQuery{
		Query: helper.BoolQuery{
			Bool: helper.BoolShould{
				Should: []any{
					helper.MatchPhrase{
						MatchPhrase: map[string]helper.MatchPhraseField{
							"title": {
								Query: query,
								Slop:  0,
								Boost: 5,
							},
						},
					},
					helper.MatchPhrase{
						MatchPhrase: map[string]helper.MatchPhraseField{
							"cast": {
								Query: query,
								Slop:  0,
								Boost: 5,
							},
						},
					},
					helper.MultiMatch{
						MultiMatch: helper.MultiMatchField{
							Query:    query,
							Fields:   []string{"title", "cast"},
							Fuzzines: "AUTO",
						},
					},
				},
				MinimalShouldMatch: 1,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	response, err := r.ES.Search(
		r.ES.Search.WithIndex("movies"),
		r.ES.Search.WithBody(&buf),
		r.ES.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var esResp struct {
		Hits struct {
			Hits []struct {
				Source entity.Movies `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(response.Body).Decode(&esResp); err != nil {
		return nil, err
	}

	var movies []entity.Movies
	for _, hit := range esResp.Hits.Hits {
		movies = append(movies, hit.Source)
	}

	return movies, nil
}
