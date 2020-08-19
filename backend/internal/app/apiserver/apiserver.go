package apiserver

import (
	authDeliveryHTTP "backend/internal/app/auth/delivery/http"
	authDeliveryMiddleware "backend/internal/app/auth/delivery/http/middleware"
	authRepositoryMongo "backend/internal/app/auth/repository/mongo"
	authUsecase "backend/internal/app/auth/usecase"
	cartDeliveryHTTP "backend/internal/app/cart/delivery/http"
	cartRepositoryMongo "backend/internal/app/cart/repository/mongo"
	cartUsecase "backend/internal/app/cart/usecase"
	productDeliveryHTTP "backend/internal/app/product/delivery/http"
	productRepositoryMongo "backend/internal/app/product/repository/mongo"
	productUsecase "backend/internal/app/product/usecase"
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/middleware"

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
	//e.Use(middleware.RequestID())
	// logger connect
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		CustomTimeFormat: "Jan 02 15:04:05.00",
		Format:           "[${time_custom}] [${method}] [${status}]  URI:${uri}; err: ${error}; [${latency_human}]\n",
	}))
	e.Use(middleware.CORS())

	authRepo := authRepositoryMongo.NewRepository(db, "users")
	us := authUsecase.NewUsecase(authRepo,
		timeoutContext,
		time.Duration(viper.GetInt("jwt.expire"))*time.Hour,
		[]byte(viper.GetString("jwt.signingKey")),
	)
	middl := authDeliveryMiddleware.InitMiddleware(us)

	authGroup := e.Group("/auth", middl.Authenticator)
	authDeliveryHTTP.NewHandler(authGroup, us)

	cartGroup := e.Group("/cart", middl.Authenticator)
	cartRepo := cartRepositoryMongo.NewRepository(db, "carts")
	cs := cartUsecase.NewUsecase(cartRepo, timeoutContext)
	cartDeliveryHTTP.NewHandler(cartGroup, cs)

	productGroup := e.Group("/products", middl.Authenticator)
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

	logrus.Println("MONGO CONNECTED")

	return client.Database(dbname), nil
}
