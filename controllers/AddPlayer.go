package controllers

import (
	"encoding/json"
	"net/http"

	player "github.com/gabrielluizsf/api.skyhawk/Player"
	"github.com/gabrielluizsf/api.skyhawk/database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	const (
		successMessage = "Jogador adicionado com sucesso!"
		errorMessage   = "Erro ao adicionar jogador"
	)

	var existingPlayer player.Configure
	var player player.Configure
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(player.Password), 10)
	if err != nil {
		Log("Não conseguiu criptografar a senha do player "+player.Username, r)
		http.Error(w, "PASSWORD CRYPTO ERROR", http.StatusInternalServerError)
		return
	}

	playerDoc := bson.M{
		"username": player.Username,
		"points":   player.Points,
		"password": string(passwordEncrypted),
	}

	playerCollection, client := database.Connect(r.Context())
	defer client.Disconnect(r.Context())

	// Verifica se já existe algum jogador com o mesmo username
	existingPlayerDoc := bson.M{"username": player.Username}
  if err = playerCollection.FindOne(r.Context(), existingPlayerDoc).Decode(&existingPlayer); err == nil {
		Log("Um usuário tentou se cadastrar com um username já existente", r)
		http.Error(w, "REQUEST ERROR", http.StatusBadRequest)
		return
	}

	result, err := playerCollection.InsertOne(r.Context(), playerDoc)
	if err != nil || result == nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	Log(player.Username+" criado com sucesso", r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
