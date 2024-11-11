package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var addr = flag.String("addr", ":8080", "http service address")
var roomMap map[string]*Room = make(map[string]*Room)
var seperator = []byte(":")

func main() {
	flag.Parse()

	log.Println("Starting the chat server.")

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var room *Room

		roomId := r.URL.Query().Get("roomId")

		if roomId == "" {
			newRoomId := uuid.New().String()
			room = newRoom([]byte(newRoomId))
			go room.run()

			roomMap[newRoomId] = room

			log.Println("A room is created because the received room identifier does not exist.")

		} else {
			if _, exists := roomMap[roomId]; exists {
				room = roomMap[roomId]

				log.Println("The room identifier is confirmed and assigned to the corresponding room.")

			} else {
				log.Println("Invalid room id. [ " + roomId + " ]")
				return
			}
		}

		userName := r.URL.Query().Get("userName")

		if strings.TrimSpace(userName) == "" {
			log.Println("Invalid user name. [ " + roomId + " ]")
			return
		}

		serveWs(room, userName, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
