package controllers

import (
	"log"
	"net/http"
)

func Log(message string, r *http.Request) {
	log.Printf(message + "\nHTTP-" + r.Method + "  Host:" + r.Host + "  Protocolo:" + r.Proto + "  IP:" + r.RemoteAddr)
}
