package main

import (
	"E-commerce_micro/notification-service/internal/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	log.Println("Iniciando Notification Service...")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Erro ao conectar RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal: %v", err)
	}
	defer ch.Close()

	event.ProcessNotifications(ch)
}
