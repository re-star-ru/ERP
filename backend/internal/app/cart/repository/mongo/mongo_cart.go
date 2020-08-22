package cartRepositoryMongo

import (
	"backend/internal/app/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	db *mongo.Collection
}

func (r repository) GetUsersCart(ctx context.Context, u *models.User) (*models.Cart, error) {
	cart, err := r.getCartByUser(ctx, u)

	return toCart(cart, u), err
}

func (r repository) getCartByUser(ctx context.Context, u *models.User) (*Cart, error) {
	c := &Cart{}
	c.AddedProducts = AddedProducts{}

	objID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return nil, err
	}

	res := r.db.FindOne(ctx, bson.M{"ownerID": objID})

	if res.Err() == mongo.ErrNoDocuments {
		return c, nil
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	if err := res.Decode(c); err != nil {
		return nil, err
	}
	return c, err
}

func (r repository) AddToCart(ctx context.Context, u *models.User, productID string, count int) error {
	cart, err := r.getCartByUser(ctx, u)
	if err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return err
	}

	cart.OwnerID = oid

	poid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	row := &struct {
		ProductID primitive.ObjectID `json:"productID" bson:"productID"`
		Count     int                `json:"count" bson:"count"`
	}{
		poid,
		count,
	}

	cart.AddedProducts[productID] = row
	sr := r.db.FindOneAndReplace(ctx, bson.M{"ownerID": oid}, cart, options.FindOneAndReplace().SetUpsert(true))
	if sr.Err() != nil {
		return err
	}
	return nil
}

func (r repository) RemoveFromCart(ctx context.Context, u *models.User, cartKey string) error {
	panic("implement me")
}

type Cart struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	OwnerID       primitive.ObjectID `bson:"ownerID,omitempty"`
	AddedProducts `json:"addedProducts" bson:"addedProducts"`
}

type AddedProducts map[string]*struct {
	ProductID primitive.ObjectID `json:"productID" bson:"productID"`
	Count     int                `json:"count" bson:"count"`
}

func NewRepository(DB *mongo.Database, collection string) models.CartRepository {
	return &repository{DB.Collection(collection)}
}

func toCart(cart *Cart, user *models.User) *models.Cart {
	dc := models.NewCart()
	dc.OwnerID = user.ID

	for key, v := range cart.AddedProducts {
		dc.AddedProducts[key] = &struct {
			ProductID string `json:"productID"`
			Count     int    `json:"count"`
		}{
			v.ProductID.Hex(),
			v.Count,
		}
	}

	return dc
}
