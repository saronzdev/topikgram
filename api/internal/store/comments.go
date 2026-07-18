package store

import (
	"context"
	"topikgram/api/internal/domain"
)

type CommentStore struct {
	db *Pool
}

func NewCommentStore(db *Pool) *CommentStore {
	return &CommentStore{db: db}
}

func (s *CommentStore) Create(ctx context.Context, userID int, input *domain.CreateCommentInput) (*domain.Comment, error) {
	c := &domain.Comment{}
	err := s.db.QueryRow(ctx,
		`INSERT INTO comments (content, user_id, post_id) VALUES ($1, $2, $3)
		 RETURNING id, content, created_at, user_id, post_id`,
		input.Content, userID, input.PostID,
	).Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UserID, &c.PostID)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int) ([]domain.Comment, error) {
	rows, err := s.db.Query(ctx,
		`SELECT c.id, c.content, c.created_at, c.user_id, c.post_id, u.id, u.name, u.username
		 FROM comments c
		 JOIN users u ON u.id=c.user_id
		 WHERE c.post_id=$1
		 ORDER BY c.created_at DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []domain.Comment{}
	for rows.Next() {
		c := domain.Comment{}
		u := domain.UserPublic{}
		if err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UserID, &c.PostID, &u.ID, &u.Name, &u.Username); err != nil {
			return nil, err
		}
		c.User = u
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (s *CommentStore) GetByUserID(ctx context.Context, userID int) ([]domain.Comment, error) {
	rows, err := s.db.Query(ctx,
		`SELECT c.id, c.content, c.created_at, c.user_id, c.post_id, u.id, u.name, u.username
		 FROM comments c
		 JOIN users u ON u.id=c.user_id
		 WHERE c.user_id=$1
		 ORDER BY c.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []domain.Comment{}
	for rows.Next() {
		c := domain.Comment{}
		u := domain.UserPublic{}
		if err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UserID, &c.PostID, &u.ID, &u.Name, &u.Username); err != nil {
			return nil, err
		}
		c.User = u
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
