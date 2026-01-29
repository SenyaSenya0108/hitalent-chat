package provider

import (
	"chat/internal/model"
	"chat/internal/repository"
)

type ChatProvider interface {
	Create(*model.ChatCreateDTO) (*model.Chat, error)
}

type Chat struct {
	repo *repository.ChatRepository
}

func NewChatProvider() *Chat {
	return &Chat{
		repo: repository.NewChatRepository(),
	}
}

func (p *Chat) Create(data *model.ChatCreateDTO) (*model.Chat, error) {
	chat := &model.Chat{Title: data.Title}
	response, err := p.repo.Create(chat)
	if err != nil {
		return response, err
	}

	return response, nil
}
