package rabbitObserver

import (
    "context"
    "fmt"
    "log"

    "github.com/rabbitmq/amqp091-go"
)

// Observer representa un observador para una cola de RabbitMQ.
type Observer struct {
    queueName  string
    handleFunc func(msg amqp091.Delivery)
    conn       *amqp091.Connection
    ch         *amqp091.Channel
}

// NewObserver crea una nueva instancia de Observer.
func NewObserver(conn *amqp091.Connection, queueName string, handleFunc func(msg amqp091.Delivery)) (*Observer, error) {
    ch, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("error creando canal de RabbitMQ: %w", err)
    }
    // devuelve la referencia del puntero
    return &Observer{
        queueName:  queueName,
        handleFunc: handleFunc,
        conn:       conn,
        ch:         ch,
    }, nil
}

// Start comienza a escuchar la cola y ejecuta la función de manejo cuando llega un mensaje.
func (o *Observer) Start(ctx context.Context) error {
    msgs, err := o.ch.Consume(
        o.queueName, // nombre de la cola
        "",          // consumidor
        true,        // auto-ack
        false,       // exclusivo
        false,       // no local
        false,       // espera
        nil,         // argumentos
    )
    if err != nil {
        return fmt.Errorf("error al consumir mensajes de RabbitMQ: %w", err)
    }

    go func() {
        for {
            select {
            case msg := <-msgs:
                o.handleFunc(msg)
            case <-ctx.Done():
                log.Printf("Deteniendo el observador para la cola: %s", o.queueName)
                return
            }
        }
    }()

    log.Printf("Observador iniciado para la cola: %s", o.queueName)
    return nil
}

// Stop cierra el canal y la conexión de RabbitMQ.
func (o *Observer) Stop() {
    if err := o.ch.Close(); err != nil {
        log.Printf("Error al cerrar el canal: %v", err)
    }
    if err := o.conn.Close(); err != nil {
        log.Printf("Error al cerrar la conexión: %v", err)
    }
}