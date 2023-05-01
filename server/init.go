package server

import (
	"log"
	"net/http"

	"github.com/gabrielluizsf/api.skyhawk/controllers"
	"github.com/gabrielluizsf/api.skyhawk/handlers"
)

func Start() {
	http.Handle("/addplayer", handlers.EnableCors(http.HandlerFunc(controllers.AddPlayer)))
	http.Handle("/allplayers", handlers.EnableCors(http.HandlerFunc(controllers.GetAllPlayers)))
	http.Handle("/login", handlers.EnableCors(http.HandlerFunc(controllers.Login)))
	http.Handle("/updatepoints", handlers.EnableCors(http.HandlerFunc(controllers.UpdateWSSendRequest)))
	//web sockets
	http.Handle("/ws/updatepoints", handlers.EnableCors(http.HandlerFunc(controllers.UpdatePoints)))
	log.Printf("server started")
	http.ListenAndServe(":7900", nil)
}
