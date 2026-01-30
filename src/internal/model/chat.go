package model

import "time"

type Chat struct {
	ID        uint
	Title     string
	Messages  []Message `gorm:"foreignKey:ChatID"`
	CreatedAt time.Time
}

type GetByIdRequestDTO struct {
	ChatID uint
	Limit  int
}

type ChatCreateDTO struct {
	Title string `json:"title"`
}
