package controller_test

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/likheketo/walkietalkie/controller"
	"github.com/stretchr/testify/require"
)

func TestInsertIntoRoom(t *testing.T) {
	var roomMap controller.RoomMap
	roomMap.Init()
	roomMap.InsertIntoRoom(12.8, &websocket.Conn{})
	roomMap.InsertIntoRoom(12.8, &websocket.Conn{})
	roomMap.InsertIntoRoom(12.81, &websocket.Conn{})
	require.Equal(t, 2, len(roomMap.Get(12.8)))
	require.Equal(t, 1, len(roomMap.Get(12.81)))
}

func TestGetRooms(t *testing.T) {
	var roomMap controller.RoomMap
	roomMap.Init()
	require.Empty(t, roomMap.Get(12.8))
	roomMap.InsertIntoRoom(12.8, &websocket.Conn{})
	roomMap.InsertIntoRoom(12.8, &websocket.Conn{})
	require.Equal(t, 2, len(roomMap.Get(12.8)))
}

func TestRemoveRoom(t *testing.T) {
	var roomMap controller.RoomMap
	roomMap.Init()
	roomMap.InsertIntoRoom(12.8, &websocket.Conn{})
	roomMap.InsertIntoRoom(12.8, &websocket.Conn{})
	roomMap.InsertIntoRoom(122, &websocket.Conn{})
	roomMap.RemoveFromRoom(12.8, &websocket.Conn{})
	roomMap.RemoveFromRoom(122, &websocket.Conn{})
	require.Equal(t, 0, len(roomMap.Get(122)))
	require.Equal(t, 1, len(roomMap.Get(12.8)))
}
