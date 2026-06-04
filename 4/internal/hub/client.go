package hub

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Name string
	Hub  *Hub
	Conn *websocket.Conn
	Send chan *Message
}

func NewClient(hub *Hub, conn *websocket.Conn, name string) *Client {
	return &Client{
		Name: name,
		Hub:  hub,
		Conn: conn,
		Send: make(chan *Message, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error:", err)
			}
			break
		}

		parsedMessage := parseMessage(message, c.Name)
		switch parsedMessage.Type {
		case MessageTypeWhisper:
			c.Hub.Whisper <- parsedMessage
		case MessageTypeSystem, MessageTypeChat:
			c.Hub.Broadcast <- parsedMessage
		}
	}
}

func parseMessage(message []byte, from string) *Message {
	text := string(message)
	if strings.HasPrefix(text, "/w ") {
		parts := strings.SplitN(text, " ", 3)
		if len(parts) < 3 || parts[1] == "" || parts[2] == "" {
			return &Message{Type: MessageTypeChat, From: from, Text: text}
		}
		return &Message{Type: MessageTypeWhisper, From: from, To: parts[1], Text: parts[2]}
	}

	return &Message{
		Type: MessageTypeChat,
		From: from,
		Text: text,
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			data, _ := json.Marshal(message)
			w.Write(data)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				return
			}
		}
	}
}
