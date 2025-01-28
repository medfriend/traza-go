package zero

import (
	"fmt"
	zmq4 "github.com/pebbe/zmq4/draft"
	"log"
)

func ConnZero() (map[string]int, map[string]*zmq4.Socket) {
	queues := map[string]int{
		"chat-bot": 5555,
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

	return queues, zmqSockets
}
