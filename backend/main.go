package main

import (
	"log"
	"net/http"

	"github.com/likheketo/walkietalkie/controller"
)

func main() {
	controller.AllRooms.Init()
	http.HandleFunc("/connect", controller.ConnectToFreq)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
