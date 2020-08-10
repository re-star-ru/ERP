package apiserver

import (
	productDeliveryHTTP "backend/internal/app/product/delivery/http"
	productRepositoryMongo "backend/internal/app/product/repository/mongo"
	productUsecase "backend/internal/app/product/usecase"
	userDeliveryHTTP "backend/internal/app/user/delivery/http"
	userDeliveryMiddleware "backend/internal/app/user/delivery/http/middleware"
	userRepositoryMongo "backend/internal/app/user/repository/mongo"
	userUsecase "backend/internal/app/user/usecase"
	"context"
	"time"

	"github.com/labstack/echo"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"log"

	"github.com/spf13/viper"
)

func Start() {
	dbHost := viper.GetString("database.host")
	dbName := viper.GetString("database.name")

	db, err := newMongoDB(dbHost, dbName)
	if err != nil {
		log.Fatal(err)
	}
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	e := echo.New()
	middl := userDeliveryMiddleware.InitMiddleware()
	// TODO: logger middleware
	e.Use(middl.CORS)

	userRepo := userRepositoryMongo.NewRepository(db, "user")
	us := userUsecase.NewUsecase(userRepo,
		timeoutContext,
		time.Duration(viper.GetInt("jwt.expire"))*time.Hour,
		[]byte(viper.GetString("jwt.signingKey")),
	)
	userDeliveryHTTP.NewHandler(e, us)

	productRepo := productRepositoryMongo.NewRepository(db, "product")
	ps := productUsecase.NewUsecase(productRepo, timeoutContext)
	productDeliveryHTTP.NewHandler(e, ps)

	log.Fatal(e.Start(viper.GetString("server.address")))
}

func newMongoDB(dbhost, dbname string) (*mongo.Database, error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbhost))

	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("MONGO CONNECTED")

	return client.Database(dbname), nil
}
