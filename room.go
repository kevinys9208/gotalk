package main

import (
	"log"
	"slices"
)

type Room struct {
	clients map[*Client]bool

	roomId []byte

	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

func newRoom(roomId []byte) *Room {
	return &Room{
		broadcast:  make(chan []byte, 1),
		register:   make(chan *Client, 1),
		unregister: make(chan *Client, 1),
		clients:    make(map[*Client]bool),
		roomId:     roomId,
	}
}

var fromRoom = []byte("room")

func (h *Room) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.broadcast <- slices.Concat(fromRoom, seperator, h.roomId, seperator, client.userName)
			log.Println("A new client is assigned to the room and basic information is sent.")

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
