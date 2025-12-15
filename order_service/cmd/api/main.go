package main

import (
	"log"

	"E-commerce_micro/order_service/internal/database"
	"E-commerce_micro/order_service/internal/event"
	"E-commerce_micro/order_service/internal/order"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Iniciando o Serviço de Pedidos...")

	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: .env não encontrado")
	}

	// 1. Conectar ao Banco de Dados (Postgres)
	dbPool, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("FATAL: Falha no Banco de Dados: %v", err)
	}
	defer dbPool.Close()

	rabbitMQURL := "amqp://guest:guest@localhost:5672/"
	eventPublisher, err := event.NewRabbitMQPublisher(rabbitMQURL)
	if err != nil {
		log.Fatalf("FATAL: Falha no RabbitMQ: %v", err)
	}
	defer eventPublisher.Close()

	rabbitMQConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Falha ao conectar RabbitMQ para o Listener: %v", err)
	}
	defer rabbitMQConn.Close()

	productGrpcUrl := "localhost:50051"
	productClient, err := order.NewProductClient(productGrpcUrl)
	if err != nil {
		log.Fatalf("FATAL: Falha no cliente gRPC: %v", err)
	}
	defer productClient.Close()

	orderRepo := order.NewPostgresRepository(dbPool)

	order.StartListening(rabbitMQConn, orderRepo)

	orderService := order.NewService(orderRepo, productClient, eventPublisher)

	orderHandler := order.NewHandler(orderService)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	orderHandler.RegisterRoutes(apiV1)

	port := ":8082"
	log.Printf("Servidor de Pedidos ouvindo na porta %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("FATAL: Falha no Gin: %v", err)
	}
}
