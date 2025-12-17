package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-elasticsearch/internal/entity"
	"go-elasticsearch/internal/helper"
	"go-elasticsearch/internal/model"
	"strings"
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
								Boost: 3,
							},
						},
					},
					helper.MatchPhrase{
						MatchPhrase: map[string]helper.MatchPhraseField{
							"cast": {
								Query: query,
								Slop:  0,
								Boost: 1,
							},
						},
					},
				},
				MinimalShouldMatch: 1,
			},
		},
	}

	reqBody := fmt.Sprintf(`{
  "query": {
    "bool": {
      "should": [
        {
          "match_phrase": {
            "Title": {
              "query": "%s",
              "slop": 0,
              "boost": 5
            }
          }
        },
        {
          "match_phrase": {
            "Extract": {
              "query": "%s",
              "slop": 0,
              "boost": 5
            }
          }
        },
        {
          "multi_match": {
            "query": "%s",
            "fields": ["Titlte^2", "Extract"],
            "fuzziness": "AUTO",
            "operator": "and"
          }
        }
      ],
      "minimum_should_match": 1
    }
  }
}`, query, query, query)

	fmt.Println("req body", reqBody)

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	response, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("movies"),
		// r.es.Search.WithQuery(query),
		r.es.Search.WithBody(strings.NewReader(reqBody)),
		r.es.Search.WithSize(20),
		r.es.Search.WithTrackTotalHits(false),
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
