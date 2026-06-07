package entity

import "time"

const (
	OrderStatusCreated   = "created"
	OrderStatusPacking   = "packing"
	OrderStatusArriving  = "arriving"
	OrderStatusCompleted = "completed"
	OrderStatusCanceled  = "canceled"
)

const (
	DocumentStatusPendingReview = "pending_review"
	DocumentStatusInReview      = "in_review"
	DocumentStatusApproved      = "approved"
	DocumentStatusRejected      = "rejected"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type Order struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Document struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
