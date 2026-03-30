package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/restapi/internal/service/comment"
)

type CommentRow struct {
	ID     uuid.UUID
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

type CommentStore struct {
	db *sqlx.DB
}

func NewCommentStore(db *sqlx.DB) *CommentStore {
	return &CommentStore{db: db}
}

func converCommentRowToComment(c CommentRow) comment.Comment {
	return comment.Comment{
		ID:     c.ID,
		Slug:   c.Slug.String,
		Author: c.Author.String,
		Body:   c.Body.String,
	}
}

func (d *CommentStore) GetComment(ctx context.Context, uuiid string) (comment.Comment, error) {
	stmt := `
		SELECT id, slug, body, author FROM comments
		WHERE id = $1
	`
	var cmtRow CommentRow
	row := d.db.QueryRowContext(ctx, stmt, uuiid)
	err := row.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error fetching the comment by uuid: %w", err)
	}
	return converCommentRowToComment(cmtRow), nil
}

func (d *CommentStore) PostComment(ctx context.Context, cmt comment.Comment) (comment.Comment, error) {
	stmt := `
		INSERT INTO comments(id, slug, author, body) 
		VALUES(:id, :slug, :author, :body)
	`
	cmt.ID = uuid.New()
	postRow := CommentRow{
		ID:     cmt.ID,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
	}

	row, err := d.db.NamedQueryContext(ctx, stmt, postRow)
	if err != nil {
		fmt.Println(err)
		return comment.Comment{}, fmt.Errorf("failed to create comment: %w", err)
	}
	if err := row.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return cmt, nil
}

func (d *CommentStore) UpdateComment(ctx context.Context, cmt comment.Comment) (comment.Comment, error) {
	stmt := `UPDATE comments SET slug = :slug, author = :author, body = :body WHERE id = :id`

	cmtRow := CommentRow{
		ID:     cmt.ID,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
	}

	rows, err := d.db.NamedQueryContext(ctx, stmt, cmtRow)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to update comment: %w", err)
	}
	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return cmt, nil
}

// DeleteComment - deletes a comment from the database
func (d *CommentStore) DeleteComment(ctx context.Context, id string) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM comments where id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment from the database: %w", err)
	}
	return nil
}
