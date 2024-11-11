package main

import (
	"log"
	"slices"
	"strconv"
)

type Room struct {
	clients map[*Client]bool

	roomId []byte

	roomMap *map[string]*Room

	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

func newRoom(roomId []byte, roomMap *map[string]*Room) *Room {
	return &Room{
		broadcast:  make(chan []byte, 1),
		register:   make(chan *Client, 1),
		unregister: make(chan *Client, 1),
		clients:    make(map[*Client]bool),
		roomId:     roomId,
		roomMap:    roomMap,
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

			if len(h.clients) == 0 {
				delete(roomMap, string(h.roomId))
				log.Println("The client does not exist, so we are closing the room. [ " + string(h.roomId) + " ]")
				log.Println("Number of remaining rooms: " + strconv.Itoa(len(roomMap)))
				return
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

			if len(h.clients) == 0 {
				delete(roomMap, string(h.roomId))
				log.Println("The client does not exist, so we are closing the room. [ " + string(h.roomId) + " ]")
				log.Println("Number of remaining rooms: " + strconv.Itoa(len(roomMap)))
				return
			}
		}
	}
}
