package websocket

const SendMessage = "send-message"
const JoinRoom = "join-room"
const LeaveRoom = "leave-room"
const GetCurrentUsername = "get-current-username"

type Message struct {
	Action    string `json:"action"`
	Timestamp string `json:"timeStamp,omitempty"`
	Body      string `json:"body,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Target    string `json:"target,omitempty"`
}
