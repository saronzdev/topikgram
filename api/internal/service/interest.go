package service

import (
	"context"
	"log"
	"topikgram/api/internal/store"
)

type InterestService struct {
	postStore    *store.PostStore
	interestStore *store.InterestStore
}

func NewInterestService(postStore *store.PostStore, interestStore *store.InterestStore) *InterestService {
	return &InterestService{postStore: postStore, interestStore: interestStore}
}

func (s *InterestService) UpdatePostInterest(ctx context.Context, userID, postID int, delta float64) {
	topics, err := s.postStore.GetTopicsByID(ctx, postID)
	if err != nil {
		log.Printf("[WARN] UpdatePostInterest: GetTopicsByID postID=%d err=%v", postID, err)
		return
	}
	if err := s.interestStore.UpdateWeight(ctx, userID, topics, delta); err != nil {
		log.Printf("[WARN] UpdatePostInterest: UpdateWeight userID=%d postID=%d err=%v", userID, postID, err)
	}
}
