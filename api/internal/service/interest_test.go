package service

import "testing"

func TestNewInterestService(t *testing.T) {
	s := NewInterestService(nil, nil)
	if s == nil {
		t.Fatal("NewInterestService returned nil")
	}
	if s.postStore != nil {
		t.Error("expected postStore to be nil")
	}
	if s.interestStore != nil {
		t.Error("expected interestStore to be nil")
	}
}
