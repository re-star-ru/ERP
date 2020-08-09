package productRepositoryMongo

import (
	"backend/internal/app/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type Product struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"string"`
}

type repository struct {
	db *mongo.Collection
}

func NewRepository(DB *mongo.Database, collection string) domain.ProductRepository {
	return &repository{DB.Collection(collection)}
}

func (r repository) CreateProduct(ctx context.Context, u *domain.User, p *domain.Product) error {
	m := toModel(p)
	res, err := r.db.InsertOne(ctx, m)
	if err != nil {
		return err
	}

	p.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r repository) GetProducts(ctx context.Context) ([]*domain.Product, error) {
	out := make([]*Product, 0)
	cur, err := r.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		product := new(Product)
		err := cur.Decode(product)
		if err != nil {
			return nil, err
		}

		out = append(out, product)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return toProducts(out), nil
}

func toModel(p *domain.Product) *Product {
	// todo: настроить конвертацию id
	//uid, _ := primitive.ObjectIDFromHex(b.)
	return &Product{
		Name: p.Name,
	}
}

func toProduct(b *Product) *domain.Product {
	return &domain.Product{
		ID:   b.ID.Hex(),
		Name: b.Name,
	}
}

func toProducts(bs []*Product) []*domain.Product {
	out := make([]*domain.Product, len(bs))

	for i, b := range bs {
		out[i] = toProduct(b)
	}

	return out
}
