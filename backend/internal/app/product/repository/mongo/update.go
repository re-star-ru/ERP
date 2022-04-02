package productRepositoryMongo

import (
	"backend/internal/app/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r repository) UpdateProduct(ctx context.Context, user *models.User, product models.Product) error {
	oid, _ := primitive.ObjectIDFromHex(product.ID)
	upsert := true
	opt := &options.FindOneAndReplaceOptions{
		Upsert: &upsert,
	}

	result := r.db.FindOneAndReplace(ctx, bson.M{"_id": oid}, toModel(product), opt)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
