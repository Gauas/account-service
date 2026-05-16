package infra

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func connectQueue(dsn string) *amqp.Channel {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatalf("infra: failed to connect to queue: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatalf("infra: failed to open queue channel: %v", err)
	}

	if err := ch.ExchangeDeclare("email_exchange", "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("infra: failed to declare email_exchange: %v", err)
	}

	log.Println("infra: queue connected")

	return ch
}
