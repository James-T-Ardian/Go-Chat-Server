package websocket

type MessageType int

const (
	ClientMessage MessageType = iota
	ServerMessage
)

func (mt MessageType) String() string {
	switch mt {
	case ClientMessage:
		return "ClientMessage"
	case ServerMessage:
		return "ServerMessage"
	default:
		return "Undefined"
	}
}

type Message struct {
	Type      MessageType `json:"messageType"`
	Timestamp string      `json:"timeStamp"`
	Body      string      `json:"body"`
	Sender    string      `json:"sender,omitempty"`
	Target    string      `json:"target"`
}
