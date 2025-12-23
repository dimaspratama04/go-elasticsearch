package messaging

import (
	"encoding/json"
	"errors"
	"go-elasticsearch/internal/entity"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type MoviesPublisher struct {
	Channel      *amqp091.Channel
	ExchangeName string
}

func NewMoviesPublisher(conn *amqp091.Connection, exchangeName string) (*MoviesPublisher, error) {
	if conn == nil {
		return nil, errors.New("rabbitmq connection is nil")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare topic exchange
	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &MoviesPublisher{
		Channel:      ch,
		ExchangeName: exchangeName,
	}, nil
}

func (p *MoviesPublisher) Publish(routingKey string, movies *entity.Movies) error {
	if movies == nil {
		return errors.New("movies is nil")
	}

	body, err := json.Marshal(movies)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal movies payload: %v\n", err)
		return err
	}

	err = p.Channel.Publish(
		p.ExchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)

	if err != nil {
		log.Printf(
			"[ERROR] Failed to publish message | exchange=%s | routingKey=%s | err=%v\n",
			p.ExchangeName,
			routingKey,
			err,
		)
		return err
	}

	log.Printf(
		"[INFO] Movie published successfully | exchange=%s | routingKey=%s | size=%d\n",
		p.ExchangeName,
		routingKey,
		len(body),
	)

	return nil

}
