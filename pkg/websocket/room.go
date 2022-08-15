package websocket

import (
	"fmt"
	"log"
	"time"
)

type Room struct {
	RoomName   string
	register   chan *client
	unregister chan *client
	clients    map[*client]bool
	broadcast  chan Message
}

func newRoom(roomName string) *Room {
	return &Room{
		RoomName:   roomName,
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		broadcast:  make(chan Message),
	}
}

func getTimeStamp() string {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Println(err)
	}
	return time.Now().In(loc).Format("01-02-2006 15:04:05")
}

func (room *Room) runRoom() {
	for {
		select {
		case currClient := <-room.register:
			room.clients[currClient] = true
			log.Println("Size of Connection Pool: ", len(room.clients))
			for client := range room.clients {
				client.conn.WriteJSON(Message{Action: JoinRoom, Timestamp: getTimeStamp(), Body: fmt.Sprintf("%s has joined...", currClient.Username), Target: room.RoomName, Sender: currClient.Username})
			}
		case currClient := <-room.unregister:
			delete(room.clients, currClient)
			log.Println("Size of Connection Pool: ", len(room.clients))
			for client := range room.clients {
				client.conn.WriteJSON(Message{Action: LeaveRoom, Timestamp: getTimeStamp(), Body: fmt.Sprintf("%s has left...", currClient.Username), Target: room.RoomName, Sender: currClient.Username})
			}
		case message := <-room.broadcast:
			for client := range room.clients {
				client.conn.WriteJSON(message)
			}
		}
	}
}
