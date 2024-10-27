package handlers

import (
	"log"
	"net/http"

	"github.com/TelefunkenU87/pointy-poker/server/models"
	"github.com/gorilla/mux"
)

func ServeWs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomName := vars["roomName"]

	points := []string{"1", "2", "3", "5", "8", "13", "21"}

	room := models.GetOrCreateRoom(roomName, points)

	conn, err := models.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &models.Client{
		Conn: conn,
		Send: make(chan models.Message),
		Room: room,
	}

	room.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
