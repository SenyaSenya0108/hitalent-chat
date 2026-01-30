package model

import "time"

type Message struct {
	ID        uint
	Text      string
	ChatID    uint `gorm:"foreignKey:ChatID"`
	CreatedAt time.Time
}

type AddMessageDTO struct {
	Text   string
	ChatID uint
}
