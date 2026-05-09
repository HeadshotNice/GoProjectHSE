package usecase

import (
	"context"
	"errors"
	"time"

	"HSE/internal/entity"
	"HSE/internal/usecase/authjwt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrBadRequest   = errors.New("bad request")
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
	ListByUser(ctx context.Context, userID int64) ([]entity.Order, error)
}

type DocumentsRepo interface {
	Create(ctx context.Context, userID int64, title, content string) (int64, error)
	ListByUser(ctx context.Context, userID int64) ([]entity.Document, error)
}

type Usecase struct {
	test   TestRepo
	dbtest DBTestRepo
	users  UsersRepo
	orders OrdersRepo
	docs   DocumentsRepo

	jwt *authjwt.Manager
}

func New(
	test TestRepo,
	dbtest DBTestRepo,
	users UsersRepo,
	orders OrdersRepo,
	docs DocumentsRepo,
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
	if email == "" || password == "" {
		return 0, ErrBadRequest
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	return u.users.Create(ctx, email, string(hash))
}

func (u *Usecase) Login(ctx context.Context, email, password string) (string, error) {
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

func (u *Usecase) CreateOrder(ctx context.Context, userID int64) (int64, error) {
	if userID <= 0 {
		return 0, ErrUnauthorized
	}
	return u.orders.Create(ctx, userID)
}

func (u *Usecase) ListOrders(ctx context.Context, userID int64) ([]entity.Order, error) {
	if userID <= 0 {
		return nil, ErrUnauthorized
	}
	return u.orders.ListByUser(ctx, userID)
}

func (u *Usecase) CreateDocument(ctx context.Context, userID int64, title, content string) (int64, error) {
	if userID <= 0 {
		return 0, ErrUnauthorized
	}
	if title == "" || content == "" {
		return 0, ErrBadRequest
	}
	return u.docs.Create(ctx, userID, title, content)
}

func (u *Usecase) ListDocuments(ctx context.Context, userID int64) ([]entity.Document, error) {
	if userID <= 0 {
		return nil, ErrUnauthorized
	}
	return u.docs.ListByUser(ctx, userID)
}
