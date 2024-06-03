package usecase

import (
	"context"
	"database/sql"

	"github.com/firzatullahd/blog-api/internal/config"
	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
)

type Irepository interface {
	WithTransaction() (*sql.Tx, error)

	InsertUser(ctx context.Context, tx *sql.Tx, in *entity.User) (*entity.User, error)
	FindUsers(ctx context.Context, in *model.FilterFindUser) ([]entity.User, error)
	UpdateUser(ctx context.Context, tx *sql.Tx, in *entity.User) error

	InsertTag(ctx context.Context, tx *sql.Tx, in *entity.Tag) (uint64, error)
	FindTag(ctx context.Context, in model.FilterFindTag) ([]entity.Tag, error)

	InsertPost(ctx context.Context, tx *sql.Tx, in *entity.Post) (*entity.Post, error)
	UpdatePost(ctx context.Context, tx *sql.Tx, in *entity.Post) error
	DeletePost(ctx context.Context, tx *sql.Tx, id uint64) error
	FindPosts(ctx context.Context, in model.FilterFindPost) ([]entity.Post, error)
	CountPost(ctx context.Context) (int, error)

	InsertRPostTag(ctx context.Context, tx *sql.Tx, postID, tagID uint64) error
	DeleteRPostTag(ctx context.Context, tx *sql.Tx, postID, tagID uint64) error
	FindRPostTag(ctx context.Context, in model.FilterFindRPost) ([]entity.RPost, error)
}

type Usecase struct {
	conf *config.Config
	repo Irepository
}

func NewUsecase(conf *config.Config, repo Irepository) *Usecase {
	return &Usecase{
		conf: conf,
		repo: repo,
	}
}
