package controllers

import (
	"encoding/json"
	"net/http"

	player "github.com/gabrielluizsf/api.skyhawk/Player"
	"github.com/gabrielluizsf/api.skyhawk/database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var playerDB player.Configure
	var playerResponse player.Public
	var player player.Configure

	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playerDoc := bson.M{"username": player.Username}

	playerCollection, client := database.Connect(r.Context())

	result := playerCollection.FindOne(r.Context(), playerDoc)
	if result.Err() != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}
  if err := result.Decode(&playerDB);err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(playerDB.Password), []byte(player.Password));err != nil {
		http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}
	client.Disconnect(r.Context())

	playerResponse.Username = playerDB.Username
	playerResponse.Points = playerDB.Points
	responseData, err := json.Marshal(playerResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	Log("Login feito por "+player.Username, r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
