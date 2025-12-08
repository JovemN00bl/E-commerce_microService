package main

import (
	"log"

	"E-commerce_micro/payments-service/internal/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Iniciando payments service...")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Falha ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Falha ao abrir canal: %v", err)
	}
	defer ch.Close()

	event.ProcessPayments(ch)
}
