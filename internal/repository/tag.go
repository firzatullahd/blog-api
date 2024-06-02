package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func (r *Repo) InsertTag(ctx context.Context, tx *sql.Tx, in *entity.Tag) (uint64, error) {
	logCtx := fmt.Sprintf("%T.InsertTag", r)
	var id uint64
	err := tx.QueryRowContext(ctx, `insert into tags(label) values ($1) returning id`, in.Label).Scan(&id)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return 0, err
	}

	return id, nil
}

// func (r *Repo) UpdateTag(ctx context.Context, tx *sql.Tx, in *entity.Tag) error {
// logCtx := fmt.Sprintf("%T.UpdateTag", r)
// _, err := tx.ExecContext(ctx, `update tags set posts = $2, updated_at = now() where id = $1`, in.Label, in.Posts)
// if err != nil {
// logger.Error(ctx, logCtx, err)
// return err
// }
//
// return nil
// }
func (r *Repo) FindTag(ctx context.Context, in model.FilterFindTag) ([]entity.Tag, error) {
	logCtx := fmt.Sprintf("%T.FindTag", r)

	rows, err := r.dbRead.QueryContext(ctx, `select id, label, created_at, updated_at from tags where label deleted_at isnull and in ($1)`, in.Label)
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
