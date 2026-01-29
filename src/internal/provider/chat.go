package provider

import (
	"chat/internal/model"
	"chat/internal/repository"
)

type ChatProvider struct {
	repo *repository.ChatRepository
}

func NewChatProvider() *ChatProvider {
	return &ChatProvider{
		repo: repository.NewChatRepository(),
	}
}

func (p *ChatProvider) Create(data *model.ChatCreateDTO) (*model.Chat, error) {
	chat := &model.Chat{Title: data.Title}
	response, err := p.repo.Create(chat)
	if err != nil {
		return response, err
	}

	return response, nil
}
