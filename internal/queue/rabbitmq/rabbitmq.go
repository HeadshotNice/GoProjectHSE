package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"HSE/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DocumentSubmittedEvent struct {
	DocumentID int64 `json:"document_id"`
	UserID     int64 `json:"user_id"`
}

type Client struct {
	conn      *amqp.Connection
	publishCh *amqp.Channel
	consumeCh *amqp.Channel
	cfg       config.RabbitMQConfig
}

func Open(cfg config.RabbitMQConfig) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	publishCh, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	consumeCh, err := conn.Channel()
	if err != nil {
		_ = publishCh.Close()
		_ = conn.Close()
		return nil, err
	}

	c := &Client{conn: conn, publishCh: publishCh, consumeCh: consumeCh, cfg: cfg}
	if err := c.declareTopology(); err != nil {
		_ = c.Close()
		return nil, err
	}
	return c, nil
}

func (c *Client) Close() error {
	if c.publishCh != nil {
		_ = c.publishCh.Close()
	}
	if c.consumeCh != nil {
		_ = c.consumeCh.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) declareTopology() error {
	if err := c.publishCh.ExchangeDeclare(
		c.cfg.Exchange,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := c.declareAndBindQueue(c.cfg.DocumentSubmittedQueue, c.cfg.DocumentSubmittedKey); err != nil {
		return err
	}
	return c.declareAndBindQueue(c.cfg.DocumentStatusQueue, c.cfg.DocumentStatusKey)
}

func (c *Client) declareAndBindQueue(queueName, routingKey string) error {
	if _, err := c.publishCh.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return err
	}
	return c.publishCh.QueueBind(queueName, routingKey, c.cfg.Exchange, false, nil)
}

func (c *Client) PublishDocumentSubmitted(ctx context.Context, documentID, userID int64) error {
	body, err := json.Marshal(DocumentSubmittedEvent{DocumentID: documentID, UserID: userID})
	if err != nil {
		return err
	}

	return c.publishCh.PublishWithContext(
		ctx,
		c.cfg.Exchange,
		c.cfg.DocumentSubmittedKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func (c *Client) ConsumeDocumentSubmitted(ctx context.Context) (<-chan amqp.Delivery, error) {
	deliveries, err := c.consumeCh.Consume(
		c.cfg.DocumentSubmittedQueue,
		"hse-document-review-worker",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	out := make(chan amqp.Delivery)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-deliveries:
				if !ok {
					return
				}
				out <- msg
			}
		}
	}()
	return out, nil
}

func DecodeDocumentSubmitted(body []byte) (DocumentSubmittedEvent, error) {
	var event DocumentSubmittedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return event, err
	}
	if event.DocumentID <= 0 || event.UserID <= 0 {
		return event, fmt.Errorf("bad document-submitted event")
	}
	return event, nil
}
