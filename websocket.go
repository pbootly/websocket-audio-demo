package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	hub        *Hub
	transmit   chan []byte
}

func newClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		connection: conn,
		hub:        hub,
		transmit:   make(chan []byte),
	}
}

type Hub struct {
	clients ClientList
	sync.RWMutex
	audio AudioChannel
}

func newHub(a AudioChannel) *Hub {
	return &Hub{
		clients: make(ClientList),
		audio:   a,
	}
}

func (h *Hub) wsServer(w http.ResponseWriter, r *http.Request) {
	log.Println("New WS connection")
	connection, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(connection, h)
	h.addClient(client)

	log.Println("Clients connected: ", len(h.clients))
	go client.routeAudio(&h.audio)
	go client.sendMessage()
}

func (c *Client) routeAudio(audio *AudioChannel) {
	for {
		select {
		case audioData := <-audio.channel:
			c.transmit <- audioData
		default:
		}
	}
}

func (h *Hub) addClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	h.clients[client] = true
}

func (h *Hub) removeClient(client *Client) {
	h.Lock()
	defer h.Unlock()
	// Delete client if exists
	if _, ok := h.clients[client]; ok {
		client.connection.Close()
		delete(h.clients, client)
	}
}

func (c *Client) sendMessage() {
	defer func() {
		c.hub.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.transmit:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed")
				}
				return
			}
			if err := c.connection.WriteMessage(websocket.BinaryMessage, message); err != nil {
				log.Println("Send Error: ", err)
				return
			}
		default:

		}
	}

}
