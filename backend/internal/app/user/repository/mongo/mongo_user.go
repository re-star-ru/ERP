package userRepositoryMongo

import (
	"backend/internal/app/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	DB *mongo.Collection
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
