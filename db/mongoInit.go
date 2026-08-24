package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init() (*mongo.Client, context.Context, error) {

	fmt.Print(os.Getenv("MONGO_URL"))
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return nil, nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return client, ctx, err
}

func InitNew() (*mongo.Database, context.Context, error) {

	fmt.Print(os.Getenv("MONGO_URL"))
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return nil, nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return client.Database("myutilityx"), ctx, err
}
