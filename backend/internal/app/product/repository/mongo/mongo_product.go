package productRepositoryMongo

import (
	"backend/internal/app/domain"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type Product struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`

	GUID         string `json:"guid" bson:"guid"`
	SKU          string `json:"sku" bson:"sku"`
	Description  string `json:"description" bson:"description"`
	Manufacturer string `json:"manufacturer"  bson:"manufacturer"`
	//TypeGUID     string `json:"typeGUID" bson:""`
	//TypeName     string `json:"typeName"`

	//Characteristics []Characteristic `json:"characteristics"`

	//Properties []Property `json:"properties"`

	Creator domain.User `json:"creator" bson:"creator"`

	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	LastModified time.Time `json:"lastModified" `
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

	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Println(err)
		}
	}()

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
