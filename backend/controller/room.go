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
				delete(r.Map, freq)
			}
		}
		r.mutex.Unlock()
		time.Sleep(1 * time.Minute)
	}
}

func (r *RoomMap) Init() {
	r.Map = make(map[float64][]Client)
	go r.cleanupRooms()
	go Broadcast()
}

func (r *RoomMap) Get(freq float64) []Client {
	r.mutex.RLock()
	defer r.mutex.Unlock()
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
