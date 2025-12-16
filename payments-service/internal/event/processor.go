package event

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderCreatedEvent struct {
	OrderID    string  `json:"order_id"`
	UserID     string  `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
}

func ProcessPayments(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"order_created",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Falha ao declarar fila order_created")

	_, err = ch.QueueDeclare(
		"payment_approved",
		true,
		false, false, false, nil)

	failOnError(err, "Falha ao declarar fila payment_approved")

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	failOnError(err, "Falha ao registrar consumidor.")

	log.Println("Payments worker rodando. Aguardando pedidos...")
	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Recebido pedido: %s", d.Body)

			var event OrderCreatedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				d.Nack(false, false)
				continue
			}

			log.Print("Processando pagamento para Order Id: %s (valor: %.2f)...", event.OrderID, event.TotalPrice)
			time.Sleep(2 * time.Second)

			log.Printf("Pagamento Aprovado para Order Id: %s", event.OrderID)

			publishPaymentApproved(ch, event.OrderID)
			d.Ack(false)
		}
	}()
	<-forever
}

func publishPaymentApproved(ch *amqp.Channel, orderID string) {

	err := ch.ExchangeDeclare(
		"payment_fanout",
		"fanout",
		true,
		false, false, false, nil)

	if err != nil {
		log.Fatalf("Erro ao declarar exchange: %v", err)
		return
	}

	payload := map[string]string{
		"order_id":   orderID,
		"status":     "PAID",
		"updated_at": time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(payload)
	err = ch.Publish("payment_fanout", "", false, false, amqp.Publishing{
		ContentType: "application/json", Body: body,
	})
	if err != nil {
		log.Printf("Erro ao publicar payment_approved: %v", err)
	} else {
		log.Printf("Evento payment enviado!")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}
