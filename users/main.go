package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"users/delivery"
	"users/repo"
)

const MongoDB = "mongodb://localhost:27017"
const MongoMainDB = "coursera"
const MongoTestingDB = "testDB"

// Users store users in mongo db, check permissions for requested
// methods via casbin, api via grpc or http
func main() {
	// read configs
	host := getEnv("HOST", ":11000")
	mongoAddr := getEnv("MONGO_ADDR", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "users")

	// for dev
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})

	// init app
	repo := repo.NewMongoRepo(connectMongo(context.Background(), mongoAddr), mongoDB)
	d := delivery.NewHttpDelivery(repo)
	r := chi.NewRouter()

	// setup handlers
	r.Get("/", d.GetUser)

	// run http service
	log.Info().Msgf("user service ready at %s", host)
	log.Err(http.ListenAndServe(host, r)).Send()
}

func connectMongo(ctx context.Context, path string) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(path))
	if err != nil {
		log.Fatal().Err(err).Msg("cant create client mongo")
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal().Err(err).Msg("cant ping mongo")
	}
	log.Info().Msg("MONGO OK")
	return client
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
