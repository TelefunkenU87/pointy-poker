package utils

import (
	"time"

	"github.com/TelefunkenU87/pointy-poker/server/models"
)

func CleanupRooms() {
	for {
		time.Sleep(10 * time.Minute)
		models.RoomsMu.Lock()
		for name, room := range models.Rooms {
			room.Mutex.Lock()
			if time.Since(room.LastActive) > 30*time.Minute && len(room.Clients) == 0 {
				delete(models.Rooms, name)
			}
			room.Mutex.Unlock()
		}
		models.RoomsMu.Unlock()
	}
}
