package authUsecase

import (
	"backend/internal/app/models"
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type usecase struct {
	userRepo       models.UserRepository
	contextTimeout time.Duration
	expireDuration time.Duration
	signingKey     []byte
}

func NewUsecase(u models.UserRepository, timeout, expire time.Duration, signingKey []byte) models.UserUsecase {
	return &usecase{
		userRepo:       u,
		contextTimeout: timeout,
		expireDuration: expire,
		signingKey:     signingKey,
	}
}

type AuthClaims struct {
	User *models.User `json:"user"`
	jwt.StandardClaims
}

func (u *usecase) SignUp(c context.Context, m *models.User) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	existedUser, err := u.GetByEmail(ctx, m.Email)

	if existedUser != (models.User{}) {
		return models.ErrUserAlreadyExists
	}
	if err != models.ErrNotFound && err != nil {
		return err
	}

	return u.userRepo.Create(ctx, m)
}

func (u *usecase) SignIn(ctx context.Context, usr *models.User) (string, error) {
	if ok, err := u.userRepo.ValidateUser(ctx, usr); !ok {
		return "", err
	}

	claims := AuthClaims{
		usr,
		jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(u.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(u.signingKey)
}

func (u *usecase) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return u.signingKey, nil
	})

	if err != nil && err.Error() == "token is malformed: token contains an invalid number of segments" {
		return nil, models.ErrInvalidAccessToken
	}

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		claims.User.Sanitize()
		return claims.User, nil
	}

	return nil, models.ErrInvalidAccessToken
}

func (u *usecase) GetByEmail(c context.Context, email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.userRepo.GetByEmail(ctx, email)
}
