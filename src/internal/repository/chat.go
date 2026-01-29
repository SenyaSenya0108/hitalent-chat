package repository

import (
	"context"

	"chat/internal/model"
	"chat/internal/storage"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{db: storage.GetDB()}
}

func (r *ChatRepository) Create(chat *model.Chat) (*model.Chat, error) {
	ctx := context.Background()

	if err := gorm.G[model.Chat](r.db).Create(ctx, chat); err != nil {
		return chat, err
	}

	return chat, nil
}
