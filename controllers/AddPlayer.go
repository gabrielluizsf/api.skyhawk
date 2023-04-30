package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	player "github.com/gabrielluizsf/api.skyhawk/Player"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	var player player.Configure
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Erro ao carregar variáveis de ambiente: ", err)
		}
	}
	uri := os.Getenv("DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("skyhawk").Collection("players")

	playerDoc := bson.M{"username": player.Username, "points": player.Points}

	_, err = collection.InsertOne(context.Background(), playerDoc)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Jogador adicionado com sucesso!"))
}
