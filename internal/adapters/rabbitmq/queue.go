package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"story-board-generator/internal/ports"
)

type Queue struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewQueue(url, queueName string) (*Queue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	if _, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	if err := ch.Qos(1, 0, false); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("set qos: %w", err)
	}

	return &Queue{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (q *Queue) Close() error {
	if q.channel != nil {
		_ = q.channel.Close()
	}
	if q.conn != nil {
		return q.conn.Close()
	}
	return nil
}

func (q *Queue) EnqueueStoryboardGenerate(ctx context.Context, payload ports.StoryboardGeneratePayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal queue payload: %w", err)
	}

	if err := q.channel.PublishWithContext(
		ctx,
		"",
		q.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	return nil
}

func (q *Queue) ConsumeStoryboardGenerate(ctx context.Context, handler ports.MessageHandler) error {
	deliveries, err := q.channel.Consume(
		q.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("start consumer: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("delivery channel closed")
			}

			var payload ports.StoryboardGeneratePayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				_ = msg.Nack(false, false)
				continue
			}

			if err := handler(ctx, payload); err != nil {
				_ = msg.Nack(false, true)
				continue
			}

			if err := msg.Ack(false); err != nil {
				return fmt.Errorf("ack message: %w", err)
			}
		}
	}
}
