package productRepositoryMongo

import (
	"backend/internal/app/models"
	"backend/internal/app/product"
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

	Creator models.User `json:"creator" bson:"creator"`

	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	LastModified time.Time `json:"lastModified" `
}

type repository struct {
	db *mongo.Collection
}

func (r repository) GetByID(ctx context.Context, id string) (models.Product, error) {
	p := &Product{}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Product{}, err
	}
	result := r.db.FindOne(ctx, bson.M{"_id": oid})
	if result.Err() != nil {
		return models.Product{}, result.Err()
	}
	if err := result.Decode(p); err != nil {
		return models.Product{}, err
	}
	return toProduct(p), nil
}

func NewRepository(DB *mongo.Database, collection string) product.Repository {
	return &repository{DB.Collection(collection)}
}

func (r repository) CreateProduct(ctx context.Context, u *models.User, p *models.Product) error {
	m := toModel(p)
	res, err := r.db.InsertOne(ctx, m)
	if err != nil {
		return err
	}

	p.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r repository) GetProducts(ctx context.Context) ([]models.Product, error) {
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
		pd := &Product{}
		err := cur.Decode(pd)
		if err != nil {
			return nil, err
		}

		out = append(out, pd)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return toProducts(out), nil
}

func toModel(p *models.Product) *Product {
	// todo: настроить конвертацию id
	//uid, _ := primitive.ObjectIDFromHex(b.)
	return &Product{
		Name: p.Name,
	}
}

func toProduct(b *Product) models.Product {
	return models.Product{
		ID:   b.ID.Hex(),
		Name: b.Name,
	}
}

func toProducts(bs []*Product) []models.Product {
	out := make([]models.Product, len(bs))

	for i, b := range bs {
		out[i] = toProduct(b)
	}

	return out
}
