package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/likheketo/walkietalkie/controller"
)

func main() {
	controller.AllRooms.Init()
	http.HandleFunc("/connect", controller.ConnectToFreq)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
