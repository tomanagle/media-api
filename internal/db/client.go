package db

import (
	"context"
	"log/slog"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateMongoClientParams struct {
	ConnectionString  string
	ConnectionTimeout time.Duration
	Logger            *slog.Logger
}

func createClient(parentCtx context.Context, params CreateMongoClientParams) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(parentCtx, params.ConnectionTimeout)
	defer cancel()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(params.ConnectionString))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	params.Logger.Info("connected to db")

	return client, nil
}

func Must(ctx context.Context, params CreateMongoClientParams) *mongo.Client {
	client, err := createClient(ctx, params)

	if err != nil {
		panic(err)
	}

	return client

}
