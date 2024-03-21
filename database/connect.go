package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context) (*mongo.Collection, *mongo.Client) {

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Erro ao carregar variáveis de ambiente: %v", err)
		}
	}
	uri := os.Getenv("DB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("skyhawk").Collection("players")

	return collection, client

}
