package controllers

import "log"

func logERROR(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
