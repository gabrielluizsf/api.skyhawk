package server

import (
	"net/http"

	"github.com/gabrielluizsf/api.skyhawk/controllers"
)

func Start() {
	http.HandleFunc("/addplayer", controllers.AddPlayer)
	http.ListenAndServe(":7900", nil)
}
