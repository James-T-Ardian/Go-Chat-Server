package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	Username string
	conn     *websocket.Conn
	room     *Room
	hub      *Hub
}

func newClient(username string, ws *websocket.Conn, hub *Hub) *client {
	return &client{
		Username: username,
		conn:     ws,
		hub:      hub,
	}
}

func (c *client) read() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Println("Unexpected close error: ", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(jsonMessage, &msg); err != nil {
			log.Print(err)
		}
		c.handleMessage(msg)
	}
}

func (c *client) handleMessage(msg Message) {
	switch msg.Action {
	case SendMessage:
		c.room.broadcast <- msg
	case JoinRoom:
		c.joinRoom(msg.Target)
	case LeaveRoom:
		c.room.unregister <- c
	default:
		log.Println("Invalid message 'action' field: ", msg.Action, ", by client: ", c.Username)
	}
}

func (c *client) joinRoom(roomName string) {
	room := c.hub.findRoomByName((roomName))

	if room == nil {
		c.hub.registerRoom(newRoom(roomName))
	} else if c.room != room {
		c.room = room
		c.room.register <- c
	}
}
