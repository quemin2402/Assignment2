package usecase

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"user-service/internal/domain"
)

type UserUC interface {
	Register(ctx context.Context, u *domain.User, rawPwd string) (*domain.User, error)
	Auth(ctx context.Context, username, pwd string) (string, error) // возвращаем «токен»
	GetProfile(ctx context.Context, id string) (*domain.User, error)
}

type uc struct{ repo domain.UserRepository }

func New(r domain.UserRepository) UserUC { return &uc{r} }

func (u *uc) Register(ctx context.Context, usr *domain.User, raw string) (*domain.User, error) {
	if usr.Username == "" || raw == "" {
		return nil, errors.New("invalid")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(raw), 12)
	usr.Password = string(hash)
	if err := u.repo.Create(ctx, usr); err != nil {
		return nil, err
	}
	return usr, nil
}

func (u *uc) Auth(ctx context.Context, username, pwd string) (string, error) {
	user, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd)) != nil {
		return "", errors.New("invalid credentials")
	}
	return user.ID + ":" + user.Username, nil
}

func (u *uc) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	return u.repo.GetByID(ctx, id)
}
