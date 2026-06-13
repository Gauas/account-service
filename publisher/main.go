package publisher

import amqp "github.com/rabbitmq/amqp091-go"

type Registry struct {
	Email *EmailPublisher
}

func New(queue *amqp.Channel) *Registry {
	return &Registry{
		Email: &EmailPublisher{queue: queue},
	}
}
