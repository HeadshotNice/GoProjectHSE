package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"HSE/internal/entity"
)

type fakeTestRepo struct{}

func (fakeTestRepo) Hello(ctx context.Context) (string, error) {
	return "Hello!", nil
}

type fakeDBTestRepo struct{}

func (fakeDBTestRepo) InsertLine(ctx context.Context, line string) error {
	return nil
}

type fakeUsersRepo struct {
	nextID int64
	users  map[string]*entity.User
}

func newFakeUsersRepo() *fakeUsersRepo {
	return &fakeUsersRepo{
		nextID: 1,
		users:  make(map[string]*entity.User),
	}
}

func (r *fakeUsersRepo) Create(ctx context.Context, email, passwordHash string) (int64, error) {
	if _, ok := r.users[email]; ok {
		return 0, errors.New("duplicate email")
	}
	id := r.nextID
	r.nextID++
	r.users[email] = &entity.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	return id, nil
}

func (r *fakeUsersRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := r.users[email]
	if user == nil {
		return nil, nil
	}
	return user, nil
}

type fakeDocumentsRepo struct{}

func (fakeDocumentsRepo) Create(ctx context.Context, userID int64, title, content string) (int64, error) {
	return 20, nil
}

func (fakeDocumentsRepo) ListByUser(ctx context.Context, userID int64) ([]entity.Document, error) {
	return []entity.Document{{ID: 20, UserID: userID, Title: "Doc", Content: "Content"}}, nil
}

func (fakeDocumentsRepo) UpdateStatus(ctx context.Context, documentID int64, status string) error {
	return nil
}

func newTestUsecase(users *fakeUsersRepo) *Usecase {
	return New(
		fakeTestRepo{},
		fakeDBTestRepo{},
		users,
		fakeDocumentsRepo{},
		nil,
		"test-secret",
		"test-issuer",
		time.Hour,
	)
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	uc := newTestUsecase(newFakeUsersRepo())

	if _, err := uc.Register(context.Background(), "User@Example.com", "secret1"); err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	_, err := uc.Register(context.Background(), " user@example.com ", "another1")
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestLoginChecksPassword(t *testing.T) {
	uc := newTestUsecase(newFakeUsersRepo())

	if _, err := uc.Register(context.Background(), "user@example.com", "right-password"); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if _, err := uc.Login(context.Background(), "USER@example.com", "wrong-password"); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for wrong password, got %v", err)
	}

	token, err := uc.Login(context.Background(), "USER@example.com", "right-password")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty JWT token")
	}
}

func TestRegisterRejectsWeakPassword(t *testing.T) {
	uc := newTestUsecase(newFakeUsersRepo())

	_, err := uc.Register(context.Background(), "user@example.com", "12345")
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestProtectedActionsRequireUserID(t *testing.T) {
	uc := newTestUsecase(newFakeUsersRepo())

	if _, err := uc.CreateDocument(context.Background(), 0, "title", "content"); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for document, got %v", err)
	}
}
