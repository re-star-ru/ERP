package repo

import "go.mongodb.org/mongo-driver/mongo"

type MongoRepo struct {
	c  *mongo.Client
	db *mongo.Database
}

func NewMongoRepo(c *mongo.Client, db string) *MongoRepo {
	return &MongoRepo{
		c:  c,
		db: c.Database(db),
	}
}
