package order

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentApprovedEvent struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func StartListening(conn *amqp.Connection, repo repository) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Erro ao abrir canal no listener: %v", err)
		return
	}
	err = ch.ExchangeDeclare("payment_fanuot", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Printf("Erro ao declarar exchange: %v", err)
		return
	}

	q, err := ch.QueueDeclare("orders_update_queue",
		true, false, false, false, nil)
	if err != nil {
		log.Printf("Erro ao declarar file: %v", err)
		return
	}

	err = ch.QueueBind(q.Name, "", "payment_fanout", false, nil)
	if err != nil {
		log.Printf("Erro ao fazer bind: %v", err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Erro ao consumir: %v", err)
		return
	}
	go func() {
		log.Println("👂 Orders Service ouvindo pagamentos...")
		for d := range msgs {
			var event PaymentApprovedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Erro JSON: %v", err)
				continue
			}

			log.Printf("Atualizando pedido %s para PAID", event.OrderID)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := repo.UpdateStatus(ctx, event.OrderID, StatusPaid)
			cancel()

			if err != nil {
				log.Printf("Erro ao atualizar banco: %v", err)
			}
		}
	}()
}
