package services

import (
    "fmt"

    "github.com/rabbitmq/amqp091-go"
)

// HandleLogin maneja mensajes de la cola de trazabilidad-login.
func HandleLogin(msg amqp091.Delivery) {
    fmt.Printf("Mensaje recibido de la cola trazabilidad-login: %s\n", msg.Body)
    // Aquí podrías agregar la lógica para procesar el mensaje.
}
