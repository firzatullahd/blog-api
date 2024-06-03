package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
	"github.com/lib/pq"
)

func (r *Repo) InsertTag(ctx context.Context, tx *sql.Tx, in *entity.Tag) (uint64, error) {
	logCtx := fmt.Sprintf("%T.InsertTag", r)
	logger.Info(ctx, fmt.Sprintf("invoked InsertTag with label %s", in.Label))
	var id uint64
	err := tx.QueryRowContext(ctx, `insert into tags(label) values ($1) on conflict (label) do nothing`, in.Label).Scan(&id)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return 0, err
	}

	return id, nil
}

func (r *Repo) FindTag(ctx context.Context, in model.FilterFindTag) ([]entity.Tag, error) {
	logCtx := fmt.Sprintf("%T.FindTag", r)
	logger.Info(ctx, "invoked FindTag")
	query, args := buildQueryFindTag(in)

	rows, err := r.dbRead.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	defer rows.Close()

	tags := make([]entity.Tag, 0)
	for rows.Next() {
		var tag entity.Tag
		err := rows.Scan(&tag.ID, &tag.Label, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			logger.Error(ctx, logCtx, err)
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func buildQueryFindTag(in model.FilterFindTag) (string, []any) {
	query := `select id, label, created_at, updated_at from tags where deleted_at isnull`
	var args []any
	if len(in.Label) > 0 {
		query += " and label = any($1)"
		args = append(args, pq.Array(in.Label))
	} else if len(in.ID) > 0 {
		query += " and id = any($1)"
		args = append(args, pq.Array(in.ID))
	}

	return query, args
}
