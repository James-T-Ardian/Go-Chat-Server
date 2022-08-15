package websocket

import "github.com/gorilla/websocket"

type Hub struct {
	rooms map[*Room]bool
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[*Room]bool),
	}
}

func (hub *Hub) ServeWSHub(ws *websocket.Conn) {

	client := newClient(createPseudoRandomName(), ws, hub)
	client.read()
}

func (hub *Hub) findRoomByName(roomName string) *Room {
	var foundRoom *Room
	for room := range hub.rooms {
		if room.RoomName == roomName {
			foundRoom = room
			break
		}
	}
	return foundRoom
}

func (hub *Hub) registerRoom(room *Room) {
	hub.rooms[room] = true
}
