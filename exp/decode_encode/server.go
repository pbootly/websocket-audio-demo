package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentTime := time.Now().Format(time.RFC3339)
			err := conn.WriteMessage(websocket.TextMessage, []byte(currentTime))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func RunServer() {
	http.HandleFunc("/ws", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
