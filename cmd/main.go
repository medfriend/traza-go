package main

import (
	"log"
	"traza-go/internal/observable"
	"traza-go/internal/rabbitmq"
)

func main() {
	subject := &observable.Subject{}
	observer := &observable.MyObserver{}

	subject.Attach(observer)

	queues := []string{"queue1", "queue2", "queue3"}

	consumer := rabbitmq.NewConsumer(subject, queues)
	go consumer.Consume()
	log.Println("Iniciando consumidores de RabbitMQ para varias colas...")
	select {}
}