package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectToFreq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a request to connect!")
	freq := r.URL.Query().Get("freq")
	if freq == "" {
		fmt.Fprintln(w, "Invalid query param")
		return
	}
	floatFreq, err := strconv.ParseFloat(freq, 64)
	if err != nil {
		fmt.Fprintln(w, "Invalid query param")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintln(w, "Error upgrading connection")
		return
	}
	AllRooms.InsertIntoRoom(floatFreq, conn)
	for {
		var msg broadcastMsg
		err := conn.ReadJSON(&msg.message)
		if websocket.IsCloseError(err, 1001) {
			log.Printf("A client left freq - %f", floatFreq)
			AllRooms.RemoveFromRoom(floatFreq, conn)
			conn.Close()
			break
		}
		if err != nil {
			AllRooms.RemoveFromRoom(floatFreq, conn)
			conn.Close()
			log.Print(err)
			break
		}
		msg.client.Conn = conn
		msg.freq = floatFreq

		broadcastLine <- msg
	}
}
