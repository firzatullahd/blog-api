package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func (r *Repo) InsertPost(ctx context.Context, tx *sql.Tx, in *entity.Post) (uint64, error) {
	logCtx := fmt.Sprintf("%T.InsertPost", r)
	var id uint64
	err := tx.QueryRowContext(ctx, `insert into posts(title, content) values ($1, $2, $3) returning id`, in.Title, in.Content).Scan(&id)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return 0, err
	}

	return id, nil
}

func (r *Repo) UpdatePost(ctx context.Context, tx *sql.Tx, in *entity.Post) error {
	logCtx := fmt.Sprintf("%T.UpdatePost", r)
	_, err := tx.ExecContext(ctx, `update posts set updated_at= now(), title = $2, content = $3, publish_date = $4 where id = $1`, in.Title, in.Content, in.Status, in.PublishDate)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}

func (r *Repo) InsertRPostTag(ctx context.Context, tx *sql.Tx, postID, tagID uint64) error {
	logCtx := fmt.Sprintf("%T.InsertRPost", r)
	_, err := tx.ExecContext(ctx, `insert into public.r_post_tag(post_id, tag_id) values ($1, $2)`, postID, tagID)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}
