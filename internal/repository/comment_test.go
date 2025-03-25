//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/restapi/internal/db"
	"github.com/restapi/internal/service/comment"
	"github.com/stretchr/testify/assert"
)

func TestCommentDatabase(t *testing.T) {
	t.Run("test create", func(t *testing.T) {
		db, err := db.NewDatabase()
		assert.NoError(t, err)

		cmtStore := NewCommentStore(db.Client)
		cmt, err := cmtStore.PostComment(context.Background(), comment.Comment{
			Slug:   "slug",
			Author: "author",
			Body:   "body",
		})
		assert.NoError(t, err)

		newcmt, err := cmtStore.GetComment(context.Background(), cmt.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, "slug", newcmt.Slug)

	})
}
