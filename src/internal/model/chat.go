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

type AddChatRequestDTO struct {
	Title string `validate:"required,min=1,max=200"`
}

type AddChatResponseDTO struct {
	ID        uint
	Title     string
	CreatedAt time.Time
}

type GetChatByIDResponseDTO struct {
	ID        uint
	Title     string
	Messages  []Message
	CreatedAt time.Time
}
