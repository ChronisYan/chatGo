package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type room struct {
	forward chan []byte      // channel for incoming messages to be forwarded
	join    chan *client     // channel for handling users joining
	leave   chan *client     // channel for handling users leaving
	clients map[*client]bool // keep track of current users in the room
}

// Room methods
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward: // forwards message to all clients
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

// Create new room function
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

// type of a single user
type client struct {
	socket *websocket.Conn // websocket for client
	send   chan []byte     // channel for sending messages
	room   *room           // room user is in
}

// read and write methods for client
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP", err)
		return
	}

	client := &client{room: r, socket: socket, send: make(chan []byte, messageBufferSize)}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
