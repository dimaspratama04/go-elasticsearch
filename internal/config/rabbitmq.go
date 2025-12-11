package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQConnection(host, port, username, password, vhost string) *amqp091.Connection {
	var rabbitURL string

	if strings.Contains(host, ":") {
		rabbitURL = fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost)
	} else {
		rabbitURL = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", username, password, host, port, vhost)
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Printf("❌ Failed to connect to RabbitMQ: %v", err)
		log.Println("Continuing without RabbitMQ connection...")
		return nil
	}

	log.Printf("✅ Connected to RabbitMQ at %s", rabbitURL)
	return conn
}
