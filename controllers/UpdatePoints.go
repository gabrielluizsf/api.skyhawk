package controllers

import (
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
		if err := json.Unmarshal(message, &playerData); err != nil {
			log.Println("Error decoding player data:", err)
			return
		}
		// Update points in database
		playerDoc := bson.M{"username": playerData.Username}
		update := bson.M{"$set": bson.M{"points": playerData.Points + playerData.Points}}

		playerCollection, client := database.Connect(r.Context())

		_, err = playerCollection.UpdateOne(r.Context(), playerDoc, update)
		if err != nil {
			log.Println("Error updating points in database:", err)
			return
		}

		client.Disconnect(r.Context())

		log.Printf("Points updated for player %s: %d", playerData.Username, playerData.Points)
	}
}
func UpdateWSSendRequest(w http.ResponseWriter, r *http.Request) {
	var player player.Configure
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		log.Fatal(err)
	}

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

	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Fatal(err)
	}
	pointsToSTRING := strconv.Itoa(player.Points)
	Log("Adicionando mais "+pointsToSTRING+" pontos para o player "+player.Username, r)
}
