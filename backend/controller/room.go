package controller

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

type RoomMap struct {
	mutex sync.RWMutex
	Map   map[float64][]Client
}
type broadcastMsg struct {
	message map[string]interface{}
	freq    float64
	client  Client
}

var broadcastLine = make(chan broadcastMsg)

func Broadcast() {
	for {
		msg := <-broadcastLine
		for _, client := range AllRooms.Map[msg.freq] {
			if client != msg.client {
				err := client.Conn.WriteJSON(msg.message)
				if err != nil {
					AllRooms.RemoveFromRoom(msg.freq, msg.client.Conn)
					client.Conn.Close()
					log.Println(err)
				}
			}
		}
	}
}

func (r *RoomMap) cleanupRooms() {
	for {
		r.mutex.Lock()
		for freq, clients := range r.Map {
			if len(clients) == 0 {
				log.Printf("Deleting freq %0.2f from map", freq)
				delete(r.Map, freq)
			}
		}
		r.mutex.Unlock()
		time.Sleep(10 * time.Minute)
	}
}

func (r *RoomMap) Init() {
	r.Map = make(map[float64][]Client)
	go r.cleanupRooms()
	go Broadcast()
}

func (r *RoomMap) Get(freq float64) []Client {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	val, ok := r.Map[freq]
	if ok {
		return val
	}
	return []Client{}
}

func (r *RoomMap) InsertIntoRoom(freq float64, conn *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, ok := r.Map[freq]
	log.Printf("A new user joinded freq - %f", freq)
	if ok {
		r.Map[freq] = append(r.Map[freq], Client{Conn: conn})
	} else {
		r.Map[freq] = []Client{
			{Conn: conn},
		}
	}
}

func (r *RoomMap) RemoveFromRoom(freq float64, conn *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	clients, ok := r.Map[freq]
	if ok {
		var indexToDelete int
		for index, client := range clients {
			if client.Conn == conn {
				indexToDelete = index
				break
			}
		}
		clients[indexToDelete] = clients[len(clients)-1]
		r.Map[freq] = clients[:len(clients)-1]
	}
}
