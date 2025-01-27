package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/medfriend/shared-commons-go/util/consul"
	"github.com/medfriend/shared-commons-go/util/env"
	zmq4 "github.com/pebbe/zmq4/draft"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"traza-go/rabbit"
)

func main() {
	env.LoadEnv()
	consulClient := consul.ConnectToConsulKey(os.Getenv("SERVICE_NAME"))

	rabbitInfo, _ := consul.GetKeyValue(consulClient, "RABBIT")

	var resultRabbitmq map[string]string

	err := json.Unmarshal([]byte(rabbitInfo), &resultRabbitmq)

	if err != nil {
		return
	}

	s := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		resultRabbitmq["RABBIT_USER"],
		resultRabbitmq["RABBIT_PASSWORD"],
		resultRabbitmq["RABBIT_HOST"],
		resultRabbitmq["RABBIT_PORT"])

	conn, err := rabbit.Connect(s)

	defer func(conn *amqp091.Connection) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	queues := map[string]int{
		"queue1": 5555,
		"queue2": 5556,
		"queue3": 5557,
	}

	// Crear sockets ZeroMQ para cada puerto
	zmqSockets := make(map[string]*zmq4.Socket)

	for queue, port := range queues {
		socket, err := zmq4.NewSocket(zmq4.PUSH)
		if err != nil {
			log.Fatalf("Error creando ZeroMQ socket para la cola %s: %v", queue, err)
		}

		err = socket.Bind(fmt.Sprintf("tcp://*:%d", port))
		if err != nil {
			log.Fatalf("Error enlazando el socket ZeroMQ al puerto %d: %v", port, err)
		}

		log.Printf("ZeroMQ socket (PUSH) para %s enlazado en tcp://*:%d", queue, port)
		zmqSockets[queue] = socket
	}

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
