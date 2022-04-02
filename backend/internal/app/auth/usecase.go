package auth

import (
	"backend/internal/app/models"
	"context"
)

// UserUsecase represents the user's usecases
type Usecase interface {
	SignUp(ctx context.Context, u *models.User) error
	SignIn(ctx context.Context, u *models.User) (string, error)
	ParseToken(ctx context.Context, accessToken string) (*models.User, error)
}
