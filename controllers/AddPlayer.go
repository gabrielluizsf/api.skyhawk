package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	player "github.com/gabrielluizsf/api.skyhawk/Player"
	"github.com/gabrielluizsf/api.skyhawk/database"
	"go.mongodb.org/mongo-driver/bson"
)

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	var existingPlayer player.Configure
	var player player.Configure
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playerDoc := bson.M{"username": player.Username, "points": player.Points}

	conect, client := database.Connect()

	// Verifica se já existe algum jogador com o mesmo username
	existingPlayerDoc := bson.M{"username": player.Username}

	err = conect.FindOne(context.Background(), existingPlayerDoc).Decode(&existingPlayer)
	if err == nil {
		Log("Um usuário tentou se cadastrar com um username já existente", r)
		http.Error(w, "REQUEST ERROR", http.StatusBadRequest)
		return
	}

	_, err = conect.InsertOne(context.Background(), playerDoc)
	logERROR(err)

	err = client.Disconnect(context.Background())
	logERROR(err)
	Log(player.Username+" criado com sucesso", r)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Jogador adicionado com sucesso!"))
}
