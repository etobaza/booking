package models

import "time"

type Reservation struct {
	ID        int       `json:"id"`
	RoomID    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
