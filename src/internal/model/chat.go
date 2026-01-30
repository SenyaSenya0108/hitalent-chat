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

type AddChatDTO struct {
	Title string `validate:"required,min=1,max=200"`
}
