package publisher

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const EmailSendEvent = "email.send"

type EmailPublisher struct {
	queue *amqp.Channel
}

type EmailMessage struct {
	Type string         `json:"type"`
	To   string         `json:"to"`
	Data map[string]any `json:"data,omitempty"`
}

func (p *EmailPublisher) Send(ctx context.Context, message EmailMessage) error {
	body, err := json.Marshal(map[string]any{
		"id":      uuid.NewString(),
		"type":    EmailSendEvent,
		"payload": message,
	})
	if err != nil {
		return err
	}

	return p.queue.PublishWithContext(ctx, "", EmailSendEvent, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}
