package data

import "time"

type Ingredient struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}
