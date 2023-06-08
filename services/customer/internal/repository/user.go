package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/forstes/besafe-go/customer/services/customer/internal/domain"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByCredentials(ctx context.Context, email string, password string) (*domain.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{db: db}
}

func (s *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (email, password, first_name, last_name, phone)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at`

	args := []any{user.Email, user.PasswordHash, user.Details.FirstName, user.Details.LastName, user.Details.Phone}

	err := s.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_unique"):
			return ErrDuplicate
		default:
			return err
		}
	}
	return nil
}

func (s *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
	SELECT id, first_name, last_name, phone, email, password, created_at
	FROM users
	WHERE id = $1`

	var user domain.User
	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Details.FirstName,
		&user.Details.LastName,
		&user.Details.Phone,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (s *userRepo) GetByCredentials(ctx context.Context, email string, password string) (*domain.User, error) {
	query := `
	SELECT id, first_name, last_name, phone, email, password, created_at
	FROM users
	WHERE email = $1 AND password = $2`

	var user domain.User
	err := s.db.QueryRow(ctx, query, email, password).Scan(
		&user.ID,
		&user.Details.FirstName,
		&user.Details.LastName,
		&user.Details.Phone,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
