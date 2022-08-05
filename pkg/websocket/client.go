package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Username string
	conn     *websocket.Conn
	room     *Room
}

func NewClient(username string, ws *websocket.Conn, room *Room) *Client {
	return &Client{
		Username: username,
		conn:     ws,
		room:     room,
	}
}

func (c *Client) Read() {
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

func (c *Client) handleMessage(msg Message) {
	switch msg.Action {
	case SendMessage:
		c.room.broadcast <- msg
	case JoinRoom:
		c.room.register <- c
	case LeaveRoom:
		c.room.unregister <- c
	default:
		log.Println("Invalid message 'action' field by client: ", c.Username)
	}
}
