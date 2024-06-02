package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	customerror "github.com/firzatullahd/blog-api/internal/model/error"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (u *Usecase) Register(ctx context.Context, in *model.RegisterRequest) (*model.AuthResponse, error) {
	logCtx := fmt.Sprintf("%T.Login", u)

	if err := validateRegister(in); err != nil {
		return nil, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(in.Password), u.conf.BcryptSalt)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}
	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	checkUser, err := u.repo.FindUsers(ctx, &model.FilterFindUser{Email: &in.Email})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	if checkUser != nil {
		return nil, customerror.ErrEmailExists
	}

	user := entity.User{
		Email:    in.Email,
		Password: string(password),
		Name:     in.Name,
		Role:     entity.RoleUser.String(),
	}
	res, err := u.repo.InsertUser(ctx, tx, &user)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	accessToken, err := u.generateAccessToken(res.ID, res.Email, res.Role)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	_ = tx.Commit()

	return &model.AuthResponse{
		Email:       in.Email,
		Name:        in.Name,
		AccessToken: accessToken,
	}, nil
}

func validateRegister(in *model.RegisterRequest) error {
	_, err := mail.ParseAddress(in.Email)
	if err != nil {
		return customerror.ErrValidation
	}

	if len(in.Name) < 5 || len(in.Name) > 50 {
		return customerror.ErrValidation
	}

	if len(in.Password) < 5 || len(in.Password) > 15 {
		return customerror.ErrValidation
	}

	return nil
}

func (u *Usecase) Login(ctx context.Context, in *model.LoginRequest) (*model.AuthResponse, error) {
	logCtx := fmt.Sprintf("%T.Login", u)
	users, err := u.repo.FindUsers(ctx, &model.FilterFindUser{
		Email: &in.Email,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	if len(users) == 0 {
		return nil, customerror.ErrNotFound
	}
	user := users[0]
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, customerror.ErrWrongPass
	}

	accessToken, err := u.generateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &model.AuthResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil

}

func (u *Usecase) generateAccessToken(userId uint64, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.MyClaim{
		UserData: model.UserData{
			ID:    userId,
			Email: email,
			Role:  role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	return token.SignedString([]byte(u.conf.JWTSecretKey))
}

func (u *Usecase) GrantAdmin(ctx context.Context, in *model.LoginRequest) (*model.AuthResponse, error) {
	logCtx := fmt.Sprintf("%T.GrantAdmin", u)

	if in.SecretKey != u.conf.AdminSecretKey {
		return nil, customerror.ErrUnauthorized
	}
	users, err := u.repo.FindUsers(ctx, &model.FilterFindUser{
		Email: &in.Email,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerror.ErrNotFound
		}
		return nil, err
	}

	if len(users) == 0 {
		return nil, customerror.ErrNotFound
	}

	user := users[0]
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, customerror.ErrWrongPass
	}

	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	defer func() {
		if err != nil {
			logger.Error(ctx, logCtx, err)
			_ = tx.Rollback()
		}
	}()

	err = u.repo.UpdateUser(ctx, tx, &entity.User{
		Role:  entity.RoleAdmin.String(),
		Email: in.Email,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	tx.Commit()

	accessToken, err := u.generateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	return &model.AuthResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil

}
