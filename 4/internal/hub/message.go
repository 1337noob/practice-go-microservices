package hub

const (
	MessageTypeChat    MessageType = "chat"
	MessageTypeSystem  MessageType = "system"
	MessageTypeWhisper MessageType = "whisper"
)

type MessageType string

type Message struct {
	Type MessageType `json:"type"`
	From string      `json:"from"`
	Text string      `json:"text"`
	To   string      `json:"to,omitempty"`
}
