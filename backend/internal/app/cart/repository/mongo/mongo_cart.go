package cartRepositoryMongo

import (
	"backend/internal/app/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	db *mongo.Collection
}

func (r repository) GetUsersCart(ctx context.Context, u *domain.User) (*domain.Cart, error) {
	c := new(Cart)
	res := r.db.FindOne(ctx, bson.M{"ownerID": primitive.ObjectIDFromHex(u.ID)})
	if res.Err() != nil {
		return nil, res.Err()
	}

	if err := res.Decode(c); err != nil {
		return nil, err
	}

	return toCart(c), nil
}

func (r repository) AddToCart(ctx context.Context, u *domain.User, productID string, count int) error {
	panic("implement me")
}

func (r repository) RemoveFromCart(ctx context.Context, u *domain.User, cartKey string) error {
	panic("implement me")
}

type Cart struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	OwnerID       primitive.ObjectID `bson:"ownerID,omitempty"`
	AddedProducts map[string]struct {
		ProductID primitive.ObjectID `json:"productID" bson:"ownerID"`
		Count     int                `json:"count" bson:"count"`
	} `json:"addedProducts" bson:"addedProducts"`
}

func NewRepository(DB *mongo.Database, collection string) domain.CartRepository {
	return &repository{DB.Collection(collection)}
}

func toCart(cart *Cart) *domain.Cart {
	dc := &domain.Cart{
		OwnerID: cart.OwnerID.Hex(),
	}

	for key, v := range cart.AddedProducts {
		dc.AddedProducts[key].Count = v.Count
		pr := struct {
			ProductID string
			Count     int
		}{
			ProductID: v.ProductID.Hex(),
			Count:     v.Count,
		}

		dc.AddedProducts[key] = pr
	}

	return dc
}
