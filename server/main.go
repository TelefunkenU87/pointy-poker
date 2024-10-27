package main

import (
	"log"
	"net/http"

	"github.com/TelefunkenU87/pointy-poker/server/handlers"
	"github.com/TelefunkenU87/pointy-poker/server/utils"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.HandleFunc("/ws/{roomName}", handlers.ServeWs)

	// serve static files from app
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./app")))

	go utils.CleanupRooms()

	addr := ":8080"
	log.Printf("Server started on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
