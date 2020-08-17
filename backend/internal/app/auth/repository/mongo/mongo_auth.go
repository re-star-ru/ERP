package authRepositoryMongo

import (
	"backend/internal/app/domain"
	"context"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	DB *mongo.Collection
}

type user struct {
	ID                string `json:"id" bson:"_id"`
	Email             string `json:"email" bson:"email" validate:"required,email"`
	Name              string `json:"name"`
	EncryptedPassword string `json:"encrypted_password,omitempty" bson:"encrypted_password"`
}

func NewRepository(DB *mongo.Database, collection string) domain.UserRepository {
	return &repository{DB.Collection(collection)}
}

func (m repository) Create(ctx context.Context, user *domain.User) error {
	if err := user.BeforeCreate(); err != nil {
		return err
	}

	// TODO: автоматическую конвертацию в тип базы данных и обратно
	_, err := m.DB.InsertOne(ctx, bson.D{
		{"email", user.Email},
		{"name", user.Name},
		{"encrypted_password", user.EncryptedPassword},
	})
	if err != nil {
		return err
	}
	return nil
}

func (m repository) GetByID(ctx context.Context, id int) (domain.User, error) {
	panic("implement me")
}

func (m repository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user := domain.User{}

	res := m.DB.FindOne(ctx, bson.M{"email": email})

	// Если не найдено возвращаем пустого пользователя
	if res.Err() == mongo.ErrNoDocuments {
		return user, domain.ErrNotFound
	}

	// если
	if res.Err() != nil {
		return user, res.Err()
	}

	if err := res.Decode(&user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (m repository) ValidateUser(ctx context.Context, user *domain.User) (bool, error) {
	res := m.DB.FindOne(ctx, bson.M{"email": user.Email})
	if res.Err() == mongo.ErrNoDocuments {
		return false, domain.ErrInvalidCredentials
	}

	if res.Err() != nil {
		return false, res.Err()
	}

	u := user
	if err := res.Decode(&u); err != nil {
		return false, err
	}

	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(u.Password)) == nil, nil
}

func toUser(u user) *domain.User {
	return &domain.User{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}
}
