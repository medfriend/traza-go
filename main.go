package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/medfriend/shared-commons-go/util/consul"
	"github.com/medfriend/shared-commons-go/util/env"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"traza-go/rabbit"
	"traza-go/zero"
)

func main() {
	env.LoadEnv()
	consulClient := consul.ConnectToConsulKey(os.Getenv("SERVICE_NAME"))

	conn, err := rabbit.ConnRabbitMQ(consulClient)

	if err != nil {
		return
	}

	defer func(conn *amqp091.Connection) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	queues, zmqSockets := zero.ConnZero()

	observers := make([]*rabbit.Observer, 0)
	for queue := range queues {
		zmqSocket := zmqSockets[queue]

		handleFunc := func(msg amqp091.Delivery) {
			log.Printf("Mensaje recibido de la cola %s: %s", queue, msg.Body)

			// Procesar el mensaje (si es necesario)
			message := map[string]interface{}{
				"message": string(msg.Body),
				"time":    time.Now().Format(time.RFC3339),
			}

			// Serializar el mensaje a JSON
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error serializando el mensaje de %s: %v", queue, err)
				return
			}

			// Enviar el mensaje a través del socket ZeroMQ
			_, err = zmqSocket.SendBytes(jsonMessage, 0)
			if err != nil {
				log.Printf("Error enviando mensaje a través de ZeroMQ desde %s: %v", queue, err)
				return
			}

			log.Printf("Mensaje enviado a ZeroMQ desde %s: %s", queue, string(jsonMessage))
		}

		// Crear observador para la cola actual
		observer := rabbit.NewObserver(conn, queue, handleFunc)
		observers = append(observers, observer)
	}

	// Iniciar cada observador
	for _, observer := range observers {
		if err := observer.Start(context.Background()); err != nil {
			log.Fatalf("Error al iniciar observador: %v", err)
		}
	}

	// Configurar manejo de señales
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Println("Observadores escuchando mensajes en las colas. Presiona Ctrl+C para detener.")
	<-ctx.Done() // Mantener el programa activo hasta que se interrumpa.

	// Detener todos los observadores
	for _, observer := range observers {
		observer.Stop()
	}
	fmt.Println("Todos los observadores detenidos.")
}
