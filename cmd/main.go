package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "traza-go/rabbitObserver"

    "github.com/rabbitmq/amqp091-go"
    // gormUtil "github.com/medfriend/shared-commons-go/util/gorm"
    // "github.com/medfriend/shared-commons-go/util/consul"
	// "github.com/medfriend/shared-commons-go/util/env"
    // "gorm.io/gorm"
)

// var db *gorm.DB

func main() {
    // env.LoadEnv()
	// consulClient := consul.ConnectToConsulKey(os.Getenv("SERVICE_NAME"))

    // initDB, err := gormUtil.InitDB(
	// 	db,
	// 	consulClient,
	// 	os.Getenv("SERVICE_STATUS"),
	// )

	// if err != nil {
	// 	return
	// }
    // fmt.Println("Hola parcerito", initDB)

    // Establecer conexión con RabbitMQ
    conn, err := amqp091.Dial("amqp://admin:password@localhost:5672/")
    if err != nil {
        log.Fatalf("Error al conectar con RabbitMQ: %v", err)
    }
    defer conn.Close()

    // Funciones de manejo para las colas
    handleFuncAccion := func(msg amqp091.Delivery) {
        fmt.Printf("Mensaje recibido de la cola trazabilidad-accion: %s\n", msg.Body)
    }

    handleFuncLogin := func(msg amqp091.Delivery) {
        fmt.Printf("Mensaje recibido de la cola trazabilidad-login: %s\n", msg.Body)
    }

    // Crear un contexto que se cancela cuando se recibe una señal de interrupción
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    // Configurar y lanzar los observadores
    observers := []*rabbitObserver.Observer{}
    queues := map[string]func(amqp091.Delivery){
        "trazabilidad-accion": handleFuncAccion,
        "trazabilidad-login":  handleFuncLogin,
    }

    for queue, handler := range queues {
        observer, err := setupObserver(ctx, conn, queue, handler)
        if err != nil {
            log.Fatalf("Error al configurar el observador: %v", err)
        }
        observers = append(observers, observer)
    }

    // Mantener el programa en ejecución
    fmt.Println("Observadores escuchando mensajes en las colas. Presiona Ctrl+C para detener.")
    <-ctx.Done() // Espera hasta que el contexto expire

    // Detener todos los observadores
    for _, observer := range observers {
        observer.Stop()
    }
    fmt.Println("Todos los observadores detenidos.")
}

func setupObserver(ctx context.Context, conn *amqp091.Connection, queueName string, handleFunc func(msg amqp091.Delivery)) (*rabbitObserver.Observer, error) {
    observer, err := rabbitObserver.NewObserver(conn, queueName, handleFunc)
    if err != nil {
        return nil, fmt.Errorf("error al crear el observador para %s: %w", queueName, err)
    }

    // Iniciar el observador
    if err := observer.Start(ctx); err != nil {
        return nil, fmt.Errorf("error al iniciar el observador para %s: %w", queueName, err)
    }

    log.Printf("Observador para la cola '%s' iniciado.", queueName)
    return observer, nil
}
