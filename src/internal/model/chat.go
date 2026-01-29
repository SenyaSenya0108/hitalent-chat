package model

import "time"

type Chat struct {
	ID        uint
	Title     string
	CreatedAt time.Time
}

type ChatCreateDTO struct {
	Title string `json:"title"`
}
