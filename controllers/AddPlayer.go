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
func GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	conect, client := database.Connect()

	filter := bson.M{}

	cur, err := conect.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())

	var players []player.Configure

	for cur.Next(context.Background()) {
		var p player.Configure
		err := cur.Decode(&p)
		if err != nil {
			log.Fatal(err)
		}

		players = append(players, p)
	}

	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}
