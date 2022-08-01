package websocket

const SendMessage = "send-message"
const JoinRoom = "join-room"
const LeaveRoom = "leave-room"

type Message struct {
	Action    string `json:"action"`
	Timestamp string `json:"timeStamp"`
	Body      string `json:"body"`
	Sender    string `json:"sender,omitempty"`
	Target    string `json:"target"`
}
