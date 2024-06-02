package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func (r *Repo) FindRPostTag(ctx context.Context, in model.FilterFindRPost) ([]entity.RPost, error) {
	logCtx := fmt.Sprintf("%T.FindRPostTag", r)

	var rposts []entity.RPost
	rows, err := r.dbRead.QueryContext(ctx, `select id, post_id, tag_id, created_at, deleted_at from r_post_tag where post_id = $1 and deleted_at isnull`, in.PostID)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rpost entity.RPost
		err = rows.Scan(&rpost.ID, &rpost.PostID, &rpost.TagID, &rpost.CreatedAt, &rpost.DeletedAt)
		if err != nil {
			logger.Error(ctx, logCtx, err)
			return nil, err
		}

		rposts = append(rposts, rpost)
	}

	return rposts, nil
}

func (r *Repo) InsertRPostTag(ctx context.Context, tx *sql.Tx, postID, tagID uint64) error {
	logCtx := fmt.Sprintf("%T.InsertRPost", r)
	logger.Info(ctx, fmt.Sprintf("invoked InsertRPost with postID %d, tagId %d", postID, tagID))
	_, err := tx.ExecContext(ctx, `insert into public.r_post_tag(post_id, tag_id) values ($1, $2)`, postID, tagID)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}

func (r *Repo) DeleteRPostTag(ctx context.Context, tx *sql.Tx, postID, tagID uint64) error {
	logCtx := fmt.Sprintf("%T.DeleteRPostTag", r)
	logger.Info(ctx, fmt.Sprintf("invoked DeleteRPostTag with postId %d, tagId %d", postID, tagID))
	_, err := tx.ExecContext(ctx, `update r_post_tag set deleted_at = now() where post_id = $1 and tag_id = $2`, postID, tagID)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}
