package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	player "github.com/gabrielluizsf/api.skyhawk/Player"
	"github.com/gabrielluizsf/api.skyhawk/database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	conect, client := database.Connect(r.Context())

	filter := bson.M{}

	cur, err := conect.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())

	var players []player.Public

	for cur.Next(context.Background()) {
		var player player.Public
		err := cur.Decode(&player)
		logERROR(err)

		hashedUsername, err := bcrypt.GenerateFromPassword([]byte(player.Username), 10)
		logERROR(err)

		player.Username = string(hashedUsername)

		players = append(players, player)
	}

	err = client.Disconnect(context.Background())
	logERROR(err)
	Log("Acesso ao banco de dados de players", r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}
