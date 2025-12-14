package consumer

import (
	"encoding/json"
	"go-elasticsearch/internal/delivery/http/usecase"
	"go-elasticsearch/internal/model"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type MoviesWorker struct {
	channel      *amqp091.Channel
	usecase      *usecase.MoviesIndexUseCase
	exchangeName string
}

func NewMoviesConsumerWorker(conn *amqp091.Connection, uc *usecase.MoviesIndexUseCase, exchangeName string) (*MoviesWorker, error) {
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

	return &MoviesWorker{
		channel:      ch,
		usecase:      uc,
		exchangeName: exchangeName,
	}, nil
}

func (w *MoviesWorker) Consumer(exchangeName, routingKey string) {
	q, _ := w.channel.QueueDeclare(
		exchangeName,
		true,
		false,
		false,
		false,
		nil,
	)

	w.channel.QueueBind(
		q.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)

	msgs, _ := w.channel.Consume(
		q.Name,
		"",
		false, // MANUAL ACK
		false,
		false,
		false,
		nil,
	)

	log.Println("[INFO] Movies worker started. Waiting for messages...")

	for msg := range msgs {
		w.handleMessage(msg)
	}

}

func (w *MoviesWorker) handleMessage(msg amqp091.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC] Recovered in consumer: %v\n", r)
			_ = msg.Nack(false, true)
		}
	}()

	var movies model.Movies

	// 1️⃣ Decode payload
	if err := json.Unmarshal(msg.Body, &movies); err != nil {
		log.Printf("[ERROR] Invalid JSON payload: %v\n", err)
		_ = msg.Nack(false, false) // discard
		return
	}

	// 2️⃣ Basic validation
	if movies.Title == "" {
		log.Println("[ERROR] Invalid movie data: title is empty")
		_ = msg.Nack(false, false)
		return
	}

	// 3️⃣ Process business logic
	if err := w.usecase.IndexMovies(&movies); err != nil {
		log.Printf("[ERROR] Failed to index movie ID=%d : %v\n", movies.ID, err)
		_ = msg.Nack(false, true) // retry
		return
	}

	// 4️⃣ ACK success
	log.Printf("[SUCCESS] Movie indexed to Elasticsearch | ID=%d | Title=%s\n", movies.ID, movies.Title)
	_ = msg.Ack(false)
}
