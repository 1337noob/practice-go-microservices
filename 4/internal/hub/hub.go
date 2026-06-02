package hub

import "log"

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan *Message
	Whisper    chan *Message
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan *Message, 256),
		Whisper:    make(chan *Message, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.Name] = client
			msg := &Message{
				From: client.Name,
				Type: MessageTypeSystem,
				Text: "joined to chat",
			}
			h.Broadcast <- msg
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.Name]; ok {
				delete(h.Clients, client.Name)
				close(client.Send)
				msg := &Message{
					From: client.Name,
					Type: MessageTypeSystem,
					Text: "leave chat",
				}
				h.Broadcast <- msg
			}
		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.Name)
				}
			}
		case message := <-h.Whisper:
			to, ok := h.Clients[message.To]
			if !ok {
				log.Println("No to", message.To)
				continue
			}
			to.Send <- message
		}
	}
}
