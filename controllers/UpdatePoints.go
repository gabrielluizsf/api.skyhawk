package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	player "github.com/gabrielluizsf/api.skyhawk/Player"
	"github.com/gabrielluizsf/api.skyhawk/database"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

var upgrader = websocket.Upgrader{}

func UpdatePoints(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer conn.Close()

	// Read messages from WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		// Decode player data
		var playerData player.Configure
		err = json.Unmarshal(message, &playerData)
		if err != nil {
			log.Println("Error decoding player data:", err)
			return
		}
		// Update points in database
		playerDoc := bson.M{"username": playerData.Username}
		update := bson.M{"$set": bson.M{"points": playerData.Points + playerData.Points}}

		conect, client := database.Connect()

		_, err = conect.UpdateOne(context.Background(), playerDoc, update)
		if err != nil {
			log.Println("Error updating points in database:", err)
			return
		}

		client.Disconnect(context.Background())

		log.Printf("Points updated for player %s: %d", playerData.Username, playerData.Points)
	}
}
func UpdateWSSendRequest(w http.ResponseWriter, r *http.Request) {
	var player player.Configure
	err := json.NewDecoder(r.Body).Decode(&player)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:7900/ws/updatepoints", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	data := map[string]interface{}{
		"username": player.Username,
		"points":   player.Points,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Fatal(err)
	}
	pointsToSTRING := strconv.Itoa(player.Points)
	Log("Adicionando mais "+pointsToSTRING+" pontos para o player "+player.Username, r)
}
