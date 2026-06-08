package entity

import "time"

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

type Document struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
