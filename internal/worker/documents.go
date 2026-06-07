package worker

import (
	"context"
	"log"
	"time"

	"HSE/internal/entity"
	"HSE/internal/queue/rabbitmq"
	"HSE/internal/usecase"
)

type DocumentReviewWorker struct {
	events *rabbitmq.Client
	uc     *usecase.Usecase
	log    *log.Logger
}

func NewDocumentReviewWorker(events *rabbitmq.Client, uc *usecase.Usecase, logger *log.Logger) *DocumentReviewWorker {
	return &DocumentReviewWorker{
		events: events,
		uc:     uc,
		log:    logger,
	}
}

func (w *DocumentReviewWorker) Run(ctx context.Context) error {
	deliveries, err := w.events.ConsumeDocumentSubmitted(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-deliveries:
			if !ok {
				return nil
			}
			if err := w.process(ctx, msg.Body); err != nil {
				w.log.Printf("document review worker failed: %v", err)
				_ = msg.Nack(false, true)
				continue
			}
			_ = msg.Ack(false)
		}
	}
}

func (w *DocumentReviewWorker) process(ctx context.Context, body []byte) error {
	event, err := rabbitmq.DecodeDocumentSubmitted(body)
	if err != nil {
		return err
	}

	statuses := []string{
		entity.DocumentStatusInReview,
		entity.DocumentStatusApproved,
	}

	for _, status := range statuses {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
		if err := w.uc.UpdateDocumentStatus(ctx, event.DocumentID, status); err != nil {
			return err
		}
		w.log.Printf("document %d status changed to %s", event.DocumentID, status)
	}
	return nil
}
