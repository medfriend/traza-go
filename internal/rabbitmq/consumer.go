package rabbitmq

import (
	"log"
	"time"

	"traza-go/internal/observable"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	subject *observable.Subject
	queues  []string
}

func NewConsumer(subject *observable.Subject, queues []string) *Consumer {
	return &Consumer{
		subject: subject,
		queues:  queues,
	}
}

func (c *Consumer) Consume() {
	var conn *amqp091.Connection
	var err error

	// Reintentar la conexión si falla
	for {
		conn, err = amqp091.Dial("amqp://guest:guest@localhost:5672/")
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	for _, queueName := range c.queues {
		_, err := ch.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %v", err)
		}

		msgs, err := ch.Consume(
			queueName,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}

		go func(queue string, msgs <-chan amqp091.Delivery) {
			for msg := range msgs {
				log.Printf("Received a message from %s: %s", queue, msg.Body)
				c.subject.Notify(string(msg.Body))
			}
		}(queueName, msgs)
	}

	select {}
}