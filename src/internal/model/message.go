package model

import "time"

type Message struct {
	ID        uint
	Text      string
	ChatID    uint `gorm:"foreignKey:ChatID"`
	CreatedAt time.Time
}

type AddMessageRequestDTO struct {
	Text   string `validate:"required,min=1,max=5000"`
	ChatID uint
}

type AddMessageResponseDTO struct {
	ID        uint
	Text      string
	ChatID    uint
	CreatedAt time.Time
}
