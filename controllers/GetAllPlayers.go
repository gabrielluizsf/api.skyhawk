package controllers

import (
	"encoding/json"
	"net/http"

	player "github.com/gabrielluizsf/api.skyhawk/Player"
	"github.com/gabrielluizsf/api.skyhawk/database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	playerCollection, client := database.Connect(r.Context())

	filter := bson.M{}

	cur, err := playerCollection.Find(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(r.Context())

	var players []player.Public

	for cur.Next(r.Context()) {
		var player player.Public
		err := cur.Decode(&player)
		logERROR(err)

		hashedUsername, err := bcrypt.GenerateFromPassword([]byte(player.Username), 10)
		logERROR(err)

		player.Username = string(hashedUsername)

		players = append(players, player)
	}

	err = client.Disconnect(r.Context())
	logERROR(err)
	Log("Acesso ao banco de dados de players", r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}
