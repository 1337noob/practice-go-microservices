package main

import (
	"flag"
	"log"
	"main/internal/hub"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "cmd/web/home.html")
}

func ServeWs(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := hub.NewClient(h, conn, name)
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func main() {
	flag.Parse()
	h := hub.NewHub()
	go h.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(h, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
