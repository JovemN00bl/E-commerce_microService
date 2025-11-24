package event

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(ctx context.Context, queueName string, body []byte) error
	Close()
}

type rabbitMQPusblisher struct {
	conn *amqp.Connection
}

func NewRabbitMQPublisher(amqURL string) (Publisher, error) {
	conn, err := amqp.Dial(amqURL)
	if err != nil {
		return nil, fmt.Errorf("Falha ao conectar ao RabbitMQ: %w", err)
	}

	log.Println("Conectado com sucesso!")
	return &rabbitMQPusblisher{conn: conn}, nil
}

func (p *rabbitMQPusblisher) Publish(ctx context.Context, queueName string, body []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("Falha ao abrir canal: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("Erro ao declarar fila: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{ContentType: "applicantion/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    time.Now()})

	if err != nil {
		return fmt.Errorf("Falha ao publicar mensagem: %w", err)
	}
	log.Printf("Evento publicado na fila: '%s'", queueName)
	return nil
}

func (p *rabbitMQPusblisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
