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
		client, ok := c.Hub.Clients[c.Name]
		if !ok {
			log.Println("client not found", c.Name)
			continue
		}

		parsedMessage := parseMessage(message, client.Name)
		switch parsedMessage.Type {
		case MessageTypeWhisper:
			c.Hub.Whisper <- parsedMessage
		case MessageTypeSystem, MessageTypeChat:
			c.Hub.Broadcast <- parsedMessage
		}
	}
}

func parseMessage(message []byte, from string) *Message {
	if strings.HasPrefix(string(message), "/w") {
		parts := strings.SplitN(string(message), " ", 3)
		return &Message{
			Type: MessageTypeWhisper,
			From: from,
			Text: parts[2],
			To:   parts[1],
		}
	}

	return &Message{
		Type: MessageTypeChat,
		From: from,
		Text: string(message),
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
