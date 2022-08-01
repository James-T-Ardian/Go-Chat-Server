package main

import (
	"fmt"
	"log"
	"net/http"

	goChatWS "github.com/James-T-Ardian/Go-Chat-Server/pkg/websocket"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var defaultRoom = goChatWS.NewRoom("Default")

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := goChatWS.NewClient("Test User", ws, defaultRoom)
	client.Read()
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})

	http.HandleFunc("/ws", serveWs)
}

func main() {
	go defaultRoom.StartRoom()
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
