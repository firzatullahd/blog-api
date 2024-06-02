package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func (r *Repo) InsertUser(ctx context.Context, tx *sql.Tx, in *entity.User) (*entity.User, error) {
	logCtx := fmt.Sprintf("%T.InsertUser", r)
	var out entity.User
	err := tx.QueryRowContext(ctx, `insert into users(email, password, name, role) values ($1, $2, $3, $4) returning id, email, name, role`, in.Email, in.Password, in.Name, in.Role).Scan(&out.ID, &out.Email, &out.Name, &out.Role)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &out, nil
}

func (r *Repo) UpdateUser(ctx context.Context, tx *sql.Tx, in *entity.User) error {
	logCtx := fmt.Sprintf("%T.UpdateUser", r)
	_, err := tx.ExecContext(ctx, `update users set role = $2 where email = $1`, in.Email, in.Role)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}

func (r *Repo) FindUsers(ctx context.Context, in *model.FilterFindUser) ([]entity.User, error) {
	logCtx := fmt.Sprintf("%T.FindUser", r)
	var users []entity.User

	rows, err := r.dbRead.QueryContext(ctx, fmt.Sprintf(`select id, email, password, name, role, created_at, updated_at, deleted_at from users where deleted_at is null and email ='%s'`, *in.Email))
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			logger.Error(ctx, logCtx, err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
