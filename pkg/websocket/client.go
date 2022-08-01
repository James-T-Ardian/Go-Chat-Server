package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Username string
	conn     *websocket.Conn
	room     *Room
}

func (c *Client) Read() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

}
