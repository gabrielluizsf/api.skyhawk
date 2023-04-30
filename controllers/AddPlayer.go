package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	player "github.com/gabrielluizsf/api.skyhawk/Player"
	"github.com/gabrielluizsf/api.skyhawk/database"
	"go.mongodb.org/mongo-driver/bson"
)

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	var player player.Configure
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playerDoc := bson.M{"username": player.Username, "points": player.Points}

	conect, client := database.Connect()
	_, err = conect.InsertOne(context.Background(), playerDoc)
	logERROR(err)

	err = client.Disconnect(context.Background())
	logERROR(err)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Jogador adicionado com sucesso!"))
}

func logERROR(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
