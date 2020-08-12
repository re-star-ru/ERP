package apiserver

import (
	authDeliveryHTTP "backend/internal/app/auth/delivery/http"
	authDeliveryMiddleware "backend/internal/app/auth/delivery/http/middleware"
	authRepositoryMongo "backend/internal/app/auth/repository/mongo"
	authUsecase "backend/internal/app/auth/usecase"
	productDeliveryHTTP "backend/internal/app/product/delivery/http"
	productRepositoryMongo "backend/internal/app/product/repository/mongo"
	productUsecase "backend/internal/app/product/usecase"
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
	authRepo := authRepositoryMongo.NewRepository(db, "users")
	us := authUsecase.NewUsecase(authRepo,
		timeoutContext,
		time.Duration(viper.GetInt("jwt.expire"))*time.Hour,
		[]byte(viper.GetString("jwt.signingKey")),
	)

	middl := authDeliveryMiddleware.InitMiddleware(us)
	// TODO: logger middleware
	e.Use(middl.CORS)

	authDeliveryHTTP.NewHandler(e, us)

	productGroup := e.Group("/products")
	productGroup.Use(middl.Authenticator)
	productRepo := productRepositoryMongo.NewRepository(db, "products")
	ps := productUsecase.NewUsecase(productRepo, timeoutContext)
	productDeliveryHTTP.NewHandler(productGroup, ps)

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
