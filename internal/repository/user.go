package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func (r *Repo) InsertUser(ctx context.Context, tx *sql.Tx, in *entity.User) (*entity.User, error) {
	logCtx := fmt.Sprintf("%T.InsertUser", r)
	var out entity.User
	err := tx.QueryRowContext(ctx, `insert into users(email, password, name) values ($1, $2, $3) returning id, email, name, role`, in.Email, in.Password, in.Name).Scan(&out.ID, &out.Name, &out.Role)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &out, nil
}

func (r *Repo) UpdateUser(ctx context.Context, tx *sql.Tx, in *entity.User) error {
	logCtx := fmt.Sprintf("%T.UpdateUser", r)
	_, err := tx.ExecContext(ctx, `update users set role = $2 where id = $1`, in.ID, in.Role)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return nil
}

func (r *Repo) FindUsers(ctx context.Context, in *model.FilterFindUser) ([]entity.User, error) {
	logCtx := fmt.Sprintf("%T.FindUser", r)
	var users []entity.User

	query, args := buildQueryFindUser(in)

	rows, err := r.dbRead.QueryContext(ctx, query, args)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user)
		if err != nil {
			logger.Error(ctx, logCtx, err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func buildQueryFindUser(filter *model.FilterFindUser) (string, map[string]interface{}) {

	args := make(map[string]interface{}, 0)
	var params []string
	if filter.Email != nil {
		params = append(params, `email = :email`)
		args["email"] = filter.Email
	}

	if len(filter.ID) > 0 {
		params = append(params, `id in (:id)`)
		args["id"] = filter.ID
	}

	query := fmt.Sprintf(`select id, email, password, name, created_at, updated_at, deleted_at from users where deleted_at isnull and %s`, strings.Join(params, "and"))

	return query, args
}
