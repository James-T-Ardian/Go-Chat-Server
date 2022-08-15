package websocket

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"strings"

	petname "github.com/dustinkirkland/golang-petname"
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

// Name will be in the form of "adjective-name-pseudoRandom6DigitNumber"
func createPseudoRandomName() string {
	name := petname.Generate(2, "-")
	numberTag := pseudoRandomNumberString(6)
	return name + "-" + numberTag
}

func pseudoRandomNumberString(max int) string {
	buff := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, buff, max)
	if n != max {
		log.Println("Random number for client name failed to be generated. Will return all zeroes. Err: ", err)
		zeroes := make([]string, max)
		for i := 0; i < max; i++ {
			zeroes[i] = "0"
		}
		return strings.Join(zeroes, "")
	}

	table := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	for i := 0; i < len(buff); i++ {
		buff[i] = table[int(buff[i])%len(table)]
	}
	return string(buff)
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
		c.leaveRoom(msg.Target)
	default:
		log.Println("Invalid message 'action' field: ", msg.Action, ", by client: ", c.Username)
	}
}

func (c *client) joinRoom(roomName string) {
	room := c.hub.findRoomByName((roomName))

	if room == nil {
		createdRoom := *(newRoom(roomName))
		go createdRoom.runRoom()
		c.hub.registerRoom(&createdRoom)
	} else if c.room != room {
		c.room = room
		c.room.register <- c
	}
}

func (c *client) leaveRoom(roomName string) {
	c.room.unregister <- c
	c.room = nil
}
