package store

import (
	"context"
	"topikgram/api/internal/domain"
)

type InterestStore struct {
	db *Pool
}

func NewInterestStore(db *Pool) *InterestStore {
	return &InterestStore{db: db}
}

func (s *InterestStore) UpdateWeight(ctx context.Context, userID int, topicIDs []domain.TopicID, delta float64) error {
	if len(topicIDs) == 0 {
		return nil
	}
	for _, tid := range topicIDs {
		_, err := s.db.Exec(ctx, `
			INSERT INTO interests (user_id, topic_id, weight)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, topic_id) DO UPDATE
			SET weight = LEAST(1.0, GREATEST(0.0, interests.weight + EXCLUDED.weight))
		`, userID, int(tid), delta)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *InterestStore) HasInterests(ctx context.Context, userID int) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM interests WHERE user_id=$1 AND weight > 0)`, userID).Scan(&exists)
	return exists, err
}
