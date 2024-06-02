package handler

import (
	"context"
	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
)

type IUsecase interface {
	Register(ctx context.Context, in *model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, in *model.LoginRequest) (*model.AuthResponse, error)
	GrantAdmin(ctx context.Context, in *model.LoginRequest) (*model.AuthResponse, error)

	CreatePost(ctx context.Context, in *model.Post) (*entity.Post, error)
	UpdatePost(ctx context.Context, in *model.Post, id uint64, email string) error
	DeletePost(ctx context.Context, id uint64) error
	GetPost(ctx context.Context, id uint64) (*entity.Post, error)
}

type Handler struct {
	Usecase IUsecase
}

func NewHandler(usecase IUsecase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}
