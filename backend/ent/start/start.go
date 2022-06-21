package main

import (
	"backend/configs"
	"backend/ent/ent"
	"backend/ent/ent/user"
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	client, err := ent.Open("postgres", configs.ReadConfig().PG)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err = client.Schema.Create(context.Background()); err != nil {
		log.Fatal(err)
	}

	CreateUser(context.Background(), client)
	u, err := QueryUser(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(u)
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Create().SetAge(10).SetName("test").Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("created user: %v", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Query().Where(user.Name("test")).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	log.Printf("queried user: %v", u)
	return u, nil
}

func CreateCar(ctx context.Context, client *ent.Client) (*ent.User, error) {
	tesla, err := client.Car.Create().SetModel("Tesla").SetRegisteredAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create car: %w", err)
	}

	log.Printf("created car: %v", tesla)

	// Create a new car with model "Ford".
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)

	a8m, err := client.User.Create().SetName("a8m").SetAge(20).AddCars(tesla, ford).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)

	return a8m, nil
}
