package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "traza-go/pkg/rabbit"
    "traza-go/pkg/services"

    "github.com/joho/godotenv"
    gormUtil "github.com/medfriend/shared-commons-go/util/gorm"
    "github.com/medfriend/shared-commons-go/util/consul"
	"github.com/medfriend/shared-commons-go/util/env"
    "gorm.io/gorm"
)
var db *gorm.DB

func main() {
  env.LoadEnv()
	consulClient := consul.ConnectToConsulKey(os.Getenv("SERVICE_NAME"))

    initDB, err := gormUtil.InitDB(
		db,
		consulClient,
		os.Getenv("SERVICE_STATUS"),
	)
	if err != nil {
		return
	}

    // Cargar variables de entorno

    if err := godotenv.Load(); err != nil {
        log.Println("Advertencia: No se pudo cargar el archivo .env")
    }
    rabbitMQURL := os.Getenv("RABBITMQ_URL")

    if rabbitMQURL == "" {
        log.Fatal("RABBITMQ_URL no está configurada")
    }

    // Configurar conexión a RabbitMQ
    conn, err := rabbit.Connect(rabbitMQURL)
    if err != nil {
        log.Fatalf("Error al conectar con RabbitMQ: %v", err)
    }
    defer conn.Close()

    // Crear observadores para cada cola
    observers := []*rabbit.Observer{
        rabbit.NewObserver(conn, "trazabilidad-accion", services.HandleAccion),
        rabbit.NewObserver(conn, "trazabilidad-login", services.HandleLogin),
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
