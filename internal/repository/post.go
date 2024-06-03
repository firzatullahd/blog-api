package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"

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

func (r *Repo) FindPosts(ctx context.Context, in model.FilterFindPost) ([]entity.Post, error) {
	logCtx := fmt.Sprintf("%T.FindPost", r)

	if in.Limit == 0 {
		in.Limit = 5
	}

	if in.Page == 0 {
		in.Page = 1
	}

	offset := in.Limit * (in.Page - 1)

	var posts []entity.Post
	rows, err := r.dbRead.QueryContext(ctx, `select id, title, content, status, publish_date, created_at, updated_at, deleted_at from posts where deleted_at isnull and id = any($1) order by id asc limit $2 offset $3`, pq.Array(in.IDs), in.Limit, offset)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var post entity.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
		if err != nil {
			logger.Error(ctx, logCtx, err)
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
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

func (r *Repo) CountPost(ctx context.Context) (int, error) {
	logCtx := fmt.Sprintf("%T.CountPost", r)
	logger.Info(ctx, "invoked CountPost")

	var count int
	err := r.dbRead.QueryRowContext(ctx, `select count(id) from posts where deleted_at isnull`).Scan(&count)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return 0, err
	}

	return count, nil
}
