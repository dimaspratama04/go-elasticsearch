package messaging

import (
	"encoding/json"
	models "go-elasticsearch/internal/model"

	"github.com/rabbitmq/amqp091-go"
)

type MoviesPublisher struct {
	Channel   *amqp091.Channel
	QueueName string
}

func NewMoviesPublisher(conn *amqp091.Connection) (*MoviesPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queueName := "movies.created"
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &MoviesPublisher{
		Channel:   ch,
		QueueName: queueName,
	}, nil
}

func (p *MoviesPublisher) Publish(movie *models.Movies) error {
	body, _ := json.Marshal(movie)

	return p.Channel.Publish(
		"",
		p.QueueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
