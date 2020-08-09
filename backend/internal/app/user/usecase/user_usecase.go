package userUsecase

import (
	"backend/internal/app/domain"
	"context"
	"time"
)

type usecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUsecase(u domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &usecase{
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (u *usecase) SignUp(c context.Context, m *domain.User) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	existedUser, err := u.GetByEmail(ctx, m.Email)

	if existedUser != (domain.User{}) {
		return domain.ErrUserAlreadyExists
	}
	if err != domain.ErrNotFound && err != nil {
		return err
	}

	return u.userRepo.Create(ctx, m)
}

func (u *usecase) SignIn(ctx context.Context, usr *domain.User) (string, error) {
	if ok, err := u.userRepo.ValidateUser(ctx, usr); !ok {
		return "", err
	}
	// TODO: подписывать токен для jwt аутентификации
	return usr.Password, nil
}

func (u *usecase) ParseToken(ctx context.Context, accessToken string) (*domain.User, error) {
	panic("implement me")
}

func (u *usecase) GetByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.userRepo.GetByEmail(ctx, email)
}
