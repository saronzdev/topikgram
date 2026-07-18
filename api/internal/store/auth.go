package store

import (
	"context"
	"topikgram/api/internal/domain"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthStore struct {
	db *Pool
}

func NewAuthStore(db *Pool) *AuthStore {
	return &AuthStore{db: db}
}

func (s *AuthStore) Create(ctx context.Context, input *domain.RegisterInput) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &domain.User{}
	err = s.db.QueryRow(ctx,
		`INSERT INTO users (name, username, email, password, birthday)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, name, username, email, birthday, created_at`,
		input.Name, input.Username, input.Email, string(hash), input.Birthday,
	).Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Birthday, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthStore) Login(ctx context.Context, input *domain.LoginInput) (*domain.User, error) {
	var hash string
	u := &domain.User{}

	err := s.db.QueryRow(ctx,
		`SELECT id, name, username, email, password, birthday, created_at
		 FROM users WHERE username=$1 OR email=$1`, input.Identifier,
	).Scan(&u.ID, &u.Name, &u.Username, &u.Email, &hash, &u.Birthday, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)); err != nil {
		return nil, domain.ErrIncorrectPassword
	}
	return u, nil
}
