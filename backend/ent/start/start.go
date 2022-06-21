package main

import (
	"backend/configs"
	"backend/ent/ent"
	"backend/ent/ent/user"
	"context"
	"fmt"
	"log"

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
