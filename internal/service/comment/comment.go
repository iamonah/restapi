package comment

import (
	"context"
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

type Comment struct {
	ID     uuid.UUID
	Slug   string
	Body   string
	Author string
}

type CommentStore interface {
	GetComment(ctx context.Context, id string) (Comment, error)
	PostComment(ctx context.Context, cmt Comment) (Comment, error)
	UpdateComment(ctx context.Context, cmt Comment) (Comment, error)
	DeleteComment(ctx context.Context, id string) error
}

type Service struct {
	store CommentStore
}

func NewService(Store CommentStore) *Service {
	return &Service{
		store: Store,
	}
}

func (s *Service) GetComment(ctx context.Context, id string) (Comment, error) {

	cmt, err := s.store.GetComment(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return Comment{}, fmt.Errorf("user %v not found", id) // Formatted for client
		}
		return Comment{}, fmt.Errorf("failed to get user: %w", err) // Generic, wrap for debugging
	}

	return cmt, nil
}

func (s *Service) PostComment(ctx context.Context, cmt Comment) (Comment, error) {
	cmt, err := s.store.PostComment(ctx, cmt)
	if err != nil {
		return Comment{}, ErrCreateComment
	}
	return cmt, nil
}

func (s *Service) UpdateComment(ctx context.Context, newComment Comment) (Comment, error) {
	cmt, err := s.store.UpdateComment(ctx, newComment)
	if err != nil {
		log.Errorf("an error occurred updating the comment: %s", err.Error())
		return Comment{}, err
	}
	return cmt, err
}

// DeleteComment - deletes a comment from the database by ID
func (s *Service) DeleteComment(ctx context.Context, ID string) error {
	return s.store.DeleteComment(ctx, ID)
}

func (s *Service) CreateComment(ctx context.Context, id string) (Comment, error) {
	return Comment{}, ErrNotImplemented
}
