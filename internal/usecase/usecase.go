package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"HSE/internal/entity"
	"HSE/internal/usecase/authjwt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrBadRequest    = errors.New("bad request")
	ErrAlreadyExists = errors.New("already exists")
)

type TestRepo interface {
	Hello(ctx context.Context) (string, error)
}

type DBTestRepo interface {
	InsertLine(ctx context.Context, line string) error
}

type UsersRepo interface {
	Create(ctx context.Context, email, passwordHash string) (int64, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type OrdersRepo interface {
	Create(ctx context.Context, userID int64) (int64, error)
	ListByUser(ctx context.Context, userID int64, activeOnly bool) ([]entity.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, status string) error
}

type DocumentsRepo interface {
	Create(ctx context.Context, userID int64, title, content string) (int64, error)
	ListByUser(ctx context.Context, userID int64) ([]entity.Document, error)
	UpdateStatus(ctx context.Context, documentID int64, status string) error
}

type DocumentEvents interface {
	PublishDocumentSubmitted(ctx context.Context, documentID, userID int64) error
}

type Usecase struct {
	test   TestRepo
	dbtest DBTestRepo
	users  UsersRepo
	orders OrdersRepo
	docs   DocumentsRepo
	events DocumentEvents

	jwt *authjwt.Manager
}

func New(
	test TestRepo,
	dbtest DBTestRepo,
	users UsersRepo,
	orders OrdersRepo,
	docs DocumentsRepo,
	events DocumentEvents,
	jwtSecret string,
	jwtIssuer string,
	jwtTTL time.Duration,
) *Usecase {
	return &Usecase{
		test:   test,
		dbtest: dbtest,
		users:  users,
		orders: orders,
		docs:   docs,
		events: events,
		jwt:    authjwt.New(jwtSecret, jwtIssuer, jwtTTL),
	}
}

func (u *Usecase) TestHello(ctx context.Context) (string, error) {
	return u.test.Hello(ctx)
}

func (u *Usecase) DBTestInsert(ctx context.Context, line string) error {
	if len(line) == 0 {
		return ErrBadRequest
	}
	return u.dbtest.InsertLine(ctx, line)
}

func (u *Usecase) Register(ctx context.Context, email, password string) (int64, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return 0, ErrBadRequest
	}
	if len(password) < 6 {
		return 0, fmt.Errorf("%w: password must contain at least 6 characters", ErrBadRequest)
	}
	existing, err := u.users.FindByEmail(ctx, email)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, ErrAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	return u.users.Create(ctx, email, string(hash))
}

func (u *Usecase) Login(ctx context.Context, email, password string) (string, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return "", ErrBadRequest
	}
	user, err := u.users.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrUnauthorized
	}
	return u.jwt.Issue(user.ID)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (u *Usecase) CreateOrder(ctx context.Context, userID int64) (int64, error) {
	if userID <= 0 {
		return 0, ErrUnauthorized
	}
	orderID, err := u.orders.Create(ctx, userID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (u *Usecase) ListOrders(ctx context.Context, userID int64, activeOnly bool) ([]entity.Order, error) {
	if userID <= 0 {
		return nil, ErrUnauthorized
	}
	return u.orders.ListByUser(ctx, userID, activeOnly)
}

func (u *Usecase) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	if orderID <= 0 || status == "" {
		return ErrBadRequest
	}
	switch status {
	case entity.OrderStatusCreated, entity.OrderStatusPacking, entity.OrderStatusArriving, entity.OrderStatusCompleted, entity.OrderStatusCanceled:
		return u.orders.UpdateStatus(ctx, orderID, status)
	default:
		return ErrBadRequest
	}
}

func (u *Usecase) CreateDocument(ctx context.Context, userID int64, title, content string) (int64, error) {
	if userID <= 0 {
		return 0, ErrUnauthorized
	}
	if title == "" || content == "" {
		return 0, ErrBadRequest
	}
	documentID, err := u.docs.Create(ctx, userID, title, content)
	if err != nil {
		return 0, err
	}
	if u.events != nil {
		if err := u.events.PublishDocumentSubmitted(ctx, documentID, userID); err != nil {
			return 0, err
		}
	}
	return documentID, nil
}

func (u *Usecase) ListDocuments(ctx context.Context, userID int64) ([]entity.Document, error) {
	if userID <= 0 {
		return nil, ErrUnauthorized
	}
	return u.docs.ListByUser(ctx, userID)
}

func (u *Usecase) UpdateDocumentStatus(ctx context.Context, documentID int64, status string) error {
	if documentID <= 0 || status == "" {
		return ErrBadRequest
	}
	switch status {
	case entity.DocumentStatusPendingReview, entity.DocumentStatusInReview, entity.DocumentStatusApproved, entity.DocumentStatusRejected:
		return u.docs.UpdateStatus(ctx, documentID, status)
	default:
		return ErrBadRequest
	}
}
