package server

import (
	"net/http"

	"github.com/gabrielluizsf/api.skyhawk/controllers"
)

func Start() {
	http.HandleFunc("/addplayer", controllers.AddPlayer)
	http.HandleFunc("/allplayers", controllers.GetAllPlayers)
	http.HandleFunc("/login", controllers.Login)
	http.ListenAndServe(":7900", nil)
}
