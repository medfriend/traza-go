package services

import (
    "fmt"

    "github.com/rabbitmq/amqp091-go"
)

// HandleAccion maneja mensajes de la cola de trazabilidad-accion.
func HandleAccion(msg amqp091.Delivery) {
    fmt.Printf("Mensaje recibido de la cola trazabilidad-accion: %s\n", msg.Body)
    // Aquí podrías agregar la lógica para procesar el mensaje.
}
