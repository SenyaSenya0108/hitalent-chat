package repository

import (
	"context"
	"fmt"

	"chat/internal/model"
	"chat/internal/storage"

	"gorm.io/gorm"
)

type Chat struct {
	db *gorm.DB
}

func NewChatRepository() *Chat {
	return &Chat{db: storage.GetDB()}
}

func (r *Chat) Create(chat *model.Chat) (*model.Chat, error) {
	ctx := context.Background()

	if err := gorm.G[model.Chat](r.db).Create(ctx, chat); err != nil {
		return chat, err
	}

	return chat, nil
}

func (r *Chat) GetByID(id uint, limit int) (*model.Chat, error) {
	chat := &model.Chat{}
	result := r.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(limit) // лимит для сообщений
	}).Take(chat, id)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("%w", ErrNotFound)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return chat, nil
}

func (r *Chat) Delete(id uint) error {
	result := r.db.Delete(&model.Chat{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Chat) AddMessageToChat(message *model.Message) (*model.Message, error) {
	result := r.db.Create(message)

	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}
