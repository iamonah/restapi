package comment

import "errors"

var(
	ErrCreateComment = errors.New("failed to create comment")
	ErrFetchingComment = errors.New("failed to fetch comment by id")
	ErrNotImplemented = errors.New("not implemented")
)