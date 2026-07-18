package store

import (
	"context"
	"topikgram/api/internal/domain"

	"github.com/jackc/pgx/v5"
)

type UserStore struct {
	db *Pool
}

func NewUserStore(db *Pool) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) GetByID(ctx context.Context, id int) (*domain.User, error) {
	u := &domain.User{}
	err := s.db.QueryRow(ctx,
		`SELECT id, name, username, email, birthday, created_at
		 FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Birthday, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (s *UserStore) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	u := &domain.User{}
	err := s.db.QueryRow(ctx,
		`SELECT id, name, username, email, birthday, created_at
		 FROM users WHERE username=$1`, username,
	).Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Birthday, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (s *UserStore) List(ctx context.Context) ([]domain.UserPublic, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, name, username
		 FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByPos[domain.UserPublic])
}

func (s *UserStore) GetByIDPublic(ctx context.Context, id int) (*domain.UserPublic, error) {
	u := &domain.UserPublic{}
	err := s.db.QueryRow(ctx,
		`SELECT id, name, username FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Name, &u.Username)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (s *UserStore) GetByUsernamePublic(ctx context.Context, username string) (*domain.UserPublic, error) {
	u := &domain.UserPublic{}
	err := s.db.QueryRow(ctx,
		`SELECT id, name, username FROM users WHERE username=$1`, username,
	).Scan(&u.ID, &u.Name, &u.Username)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (s *UserStore) Follow(ctx context.Context, followerID, followeeID int) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO follows (follower_id, followee_id) VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`, followerID, followeeID)
	return err
}

func (s *UserStore) Unfollow(ctx context.Context, followerID, followeeID int) error {
	_, err := s.db.Exec(ctx,
		`DELETE FROM follows WHERE follower_id=$1 AND followee_id=$2`,
		followerID, followeeID)
	return err
}

func (s *UserStore) Followers(ctx context.Context, userID int) ([]domain.UserPublic, error) {
	rows, err := s.db.Query(ctx,
		`SELECT u.id, u.name, u.username FROM follows f
		 JOIN users u ON u.id=f.follower_id
		 WHERE f.followee_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.UserPublic])
}

func (s *UserStore) Following(ctx context.Context, userID int) ([]domain.UserPublic, error) {
	rows, err := s.db.Query(ctx,
		`SELECT u.id, u.name, u.username FROM follows f
		 JOIN users u ON u.id=f.followee_id
		 WHERE f.follower_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.UserPublic])
}
