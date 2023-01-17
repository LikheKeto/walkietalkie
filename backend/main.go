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

	if err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil); err != nil {
		log.Fatal(err)
	}
}
