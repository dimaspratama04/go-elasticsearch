package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-elasticsearch/internal/entity"
	"go-elasticsearch/internal/helper"
	"go-elasticsearch/internal/model"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
)

type MoviesESRepository struct {
	es *elasticsearch.Client
}

func NewMoviesESRepository(es *elasticsearch.Client) *MoviesESRepository {
	return &MoviesESRepository{es: es}
}

func (r *MoviesESRepository) Index(movies *model.Movies) error {
	body, _ := json.Marshal(movies)

	_, err := r.es.Index(
		"movies",
		bytes.NewReader(body),
		r.es.Index.WithDocumentID(
			fmt.Sprintf("%d", movies.ID),
		),
		r.es.Index.WithContext(context.Background()),
	)

	return err
}

func (r *MoviesESRepository) BulkIndex(movies []model.Movies) error {
	if len(movies) == 0 {
		return nil
	}

	var buf bytes.Buffer

	for _, m := range movies {
		meta := map[string]map[string]string{
			"index": {
				"_index": "movies",
				"_id":    fmt.Sprintf("%d", m.ID),
			},
		}

		metaJSON, _ := json.Marshal(meta)
		dataJSON, _ := json.Marshal(m)

		buf.Write(metaJSON)
		buf.WriteByte('\n')
		buf.Write(dataJSON)
		buf.WriteByte('\n')
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.es.Bulk(
		bytes.NewReader(buf.Bytes()),
		r.es.Bulk.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index error: %s", res.String())
	}

	return nil
}

func (r *MoviesESRepository) Search(query string) ([]entity.Movies, error) {
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

	response, err := r.es.Search(
		r.es.Search.WithIndex("movies"),
		r.es.Search.WithBody(&buf),
		r.es.Search.WithTrackTotalHits(true),
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
