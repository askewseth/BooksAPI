package main

import (
	"fmt"
	"net/http"

	"github.com/askewseth/kubernetes/api"
	log "github.com/sirupsen/logrus"
)

func main() {

	router := api.GetRouter()

	fmt.Println("Listening on http://localhost:5555/")
	err := http.ListenAndServe(":5555", router)
	if err != nil {
		log.Errorf("Error on ListenAndServe: %v", err)
	}
}
