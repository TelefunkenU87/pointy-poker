package models

import (
	"sync"
	"time"
)

type Room struct {
	Name       string
	Points     []string
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	LastActive time.Time
	Mutex      sync.Mutex
}

func NewRoom(name string, points []string) *Room {
	return &Room{
		Name:       name,
		Points:     points,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		LastActive: time.Now(),
	}
}

var (
	Rooms   = make(map[string]*Room)
	RoomsMu sync.Mutex
)

func GetOrCreateRoom(name string, points []string) *Room {
	RoomsMu.Lock()
	defer RoomsMu.Unlock()
	room, exists := Rooms[name]
	if !exists {
		room = NewRoom(name, points)
		Rooms[name] = room
		go room.Run()
	}
	room.LastActive = time.Now()
	return room
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Mutex.Lock()
			r.Clients[client] = true
			r.Mutex.Unlock()
			r.LastActive = time.Now()
		case client := <-r.Unregister:
			r.Mutex.Lock()
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
			r.Mutex.Unlock()
			r.LastActive = time.Now()
		case message := <-r.Broadcast:
			r.Mutex.Lock()
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
			r.Mutex.Unlock()
			r.LastActive = time.Now()
		}
	}
}
