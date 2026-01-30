package provider

import (
	"chat/internal/model"
	"chat/internal/repository"
)

type ChatProvider interface {
	Create(*model.ChatCreateDTO) (*model.Chat, error)
	GetByID(data *model.GetByIdRequestDTO) (*model.Chat, error)
	AddMessageToChat(data *model.AddMessageDTO) (*model.Message, error)
	Delete(id uint) error
}

type Chat struct {
	repo *repository.Chat
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

func (p *Chat) GetByID(data *model.GetByIdRequestDTO) (*model.Chat, error) {
	chat, err := p.repo.GetByID(data.ChatID, data.Limit)
	return chat, err
}

func (p *Chat) AddMessageToChat(data *model.AddMessageDTO) (*model.Message, error) {
	message := &model.Message{Text: data.Text, ChatID: data.ChatID}

	return p.repo.AddMessageToChat(message)
}

func (p *Chat) Delete(id uint) error {
	return p.repo.Delete(id)
}
