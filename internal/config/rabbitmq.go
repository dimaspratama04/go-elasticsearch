package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQConnection(host, port, username, password, vhost string) (*amqp091.Connection, error) {
	var rabbitURL string

	if strings.Contains(host, ":") {
		rabbitURL = fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost)
	} else {
		rabbitURL = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", username, password, host, port, vhost)
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Printf("❌ Failed to connect to RabbitMQ: %v", err)
		return nil, errors.New("Continuing without RabbitMQ connection")
	}

	log.Printf("✅ Connected to RabbitMQ at %s", rabbitURL)
	return conn, nil
}
