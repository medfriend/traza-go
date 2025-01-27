package main

import (
	"fmt"
	zmq4 "github.com/pebbe/zmq4"
)

func main() {
	// Crear un socket PULL
	socket, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	// Conectar al servidor PUSH
	err = socket.Connect("tcp://localhost:5557")
	if err != nil {
		panic(err)
	}

	fmt.Println("Cliente ZeroMQ (PULL) conectado al servidor en tcp://localhost:5555")

	// Escuchar mensajes
	for {
		msg, err := socket.Recv(0)
		if err != nil {
			fmt.Printf("Error recibiendo mensaje: %v\n", err)
			continue
		}
		fmt.Printf("Mensaje recibido: %s\n", msg)
	}
}
