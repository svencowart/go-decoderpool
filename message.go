package go_decoderpool

import "time"

type Message struct {
	ID        int    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}
