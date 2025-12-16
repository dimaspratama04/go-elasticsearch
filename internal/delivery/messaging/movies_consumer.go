package messaging

import (
	"encoding/json"
	"go-elasticsearch/internal/model"
	"go-elasticsearch/internal/repository"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type MoviesConsumer struct {
	channel      *amqp091.Channel
	exchangeName string
	esRepository *repository.MoviesESRepository
}

func NewMoviesConsumer(
	conn *amqp091.Connection,
	exchangeName string,
	esRepo *repository.MoviesESRepository,
) (*MoviesConsumer, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// QoS → penting untuk batching
	if err := ch.Qos(100, 0, false); err != nil {
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	return &MoviesConsumer{
		channel:      ch,
		exchangeName: exchangeName,
		esRepository: esRepo,
	}, nil
}

func (w *MoviesConsumer) Consumer(routingKey string) error {
	queueName := "movies.es.indexer"

	q, err := w.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err := w.channel.QueueBind(
		q.Name,
		routingKey,
		w.exchangeName,
		false,
		nil,
	); err != nil {
		return err
	}

	msgs, err := w.channel.Consume(
		q.Name,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("[INFO] Movies consumer started")

	// ===== BATCH CONFIG =====
	const batchSize = 100
	const flushInterval = 3 * time.Second

	batch := make([]model.Movies, 0, batchSize)
	pendingAcks := make([]amqp091.Delivery, 0, batchSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case msg := <-msgs:
			var movie model.Movies

			if err := json.Unmarshal(msg.Body, &movie); err != nil {
				log.Printf("[ERROR] Invalid payload: %v\n", err)
				msg.Nack(false, false)
				continue
			}

			batch = append(batch, movie)
			pendingAcks = append(pendingAcks, msg)

			if len(batch) >= batchSize {
				w.flush(batch, pendingAcks)
				batch = batch[:0]
				pendingAcks = pendingAcks[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				w.flush(batch, pendingAcks)
				batch = batch[:0]
				pendingAcks = pendingAcks[:0]
			}
		}
	}
}

func (w *MoviesConsumer) flush(
	batch []model.Movies,
	acks []amqp091.Delivery,
) {
	if err := w.esRepository.BulkIndex(batch); err != nil {
		log.Printf("[ERROR] Bulk index failed: %v\n", err)
		for _, msg := range acks {
			msg.Nack(false, true) // retry
		}
		return
	}

	for _, msg := range acks {
		msg.Ack(false)
	}

	log.Printf("[SUCCESS] Indexed %d movies to Elasticsearch\n", len(batch))
}
