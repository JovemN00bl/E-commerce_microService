package event

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentApprovedEvent struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func ProcessNotifications(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"payment_fanout",
		"fanout", true, false, false, false, nil)
	failOnError(err, "Falha ao declarar exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	failOnError(err, "Falha ao declarar notificação")

	err = ch.QueueBind(
		q.Name,
		"",
		"payment_fanout",
		false,
		nil)
	failOnError(err, "Falha ao fazer o bind")

	msgs, err := ch.Consume(
		q.Name, "", true, false, false, false, nil)

	failOnError(err, "Falha ao registrar consumidor")
	log.Println("Notification service rodando. Esperando pagamentos...")

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var event PaymentApprovedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Erro JSON: %v", err)
				continue
			}
			log.Printf("Enviando E-MAIL para o pedido %s: 'seu pagamento foi APROVADO!'", event.OrderID)
		}
	}()
	<-forever

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
