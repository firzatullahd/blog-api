package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	customerror "github.com/firzatullahd/blog-api/internal/model/error"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func (r *Repo) InsertPost(ctx context.Context, tx *sql.Tx, in *entity.Post) (*entity.Post, error) {
	logCtx := fmt.Sprintf("%T.InsertPost", r)
	logger.Info(ctx, "invoked InsertPost")
	var post entity.Post
	err := tx.QueryRowContext(ctx, `insert into posts(title, content) values ($1, $2) returning id, title, content, status, publish_date, created_at, updated_at, deleted_at`, in.Title, in.Content).Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &post, nil
}

func (r *Repo) UpdatePost(ctx context.Context, tx *sql.Tx, in *entity.Post) error {
	logCtx := fmt.Sprintf("%T.UpdatePost", r)
	logger.Info(ctx, "invoked UpdatePost")
	res, err := tx.ExecContext(ctx, `update posts set updated_at= now(), title = $2, content = $3, status = $4, publish_date = $5 where id = $1`, in.ID, in.Title, in.Content, in.Status, in.PublishDate)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return customerror.ErrNotFound
	}

	return nil
}

func (r *Repo) FindPost(ctx context.Context, in model.FilterFindPost) (*entity.Post, error) {
	logCtx := fmt.Sprintf("%T.FindPost", r)

	var post entity.Post
	err := r.dbRead.QueryRowContext(ctx, `select id, title, content, status, publish_date, created_at, updated_at, deleted_at from posts where deleted_at isnull and id =$1`, in.ID).Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &post, nil
}

func (r *Repo) DeletePost(ctx context.Context, tx *sql.Tx, id uint64) error {
	logCtx := fmt.Sprintf("%T.DeletePost", r)
	logger.Info(ctx, "invoked DeletePost")

	res, err := tx.ExecContext(ctx, `update posts set deleted_at= now() where id = $1`, id)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return customerror.ErrNotFound
	}

	return nil
}
