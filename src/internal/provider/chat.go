package provider

import (
	"chat/internal/model"
	"chat/internal/repository"
)

type ChatProvider interface {
	Create(*model.AddChatRequestDTO) (*model.AddChatResponseDTO, error)
	GetByID(data *model.GetByIdRequestDTO) (*model.GetChatByIDResponseDTO, error)
	AddMessageToChat(data *model.AddMessageRequestDTO) (*model.AddMessageResponseDTO, error)
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

func (p *Chat) Create(data *model.AddChatRequestDTO) (*model.AddChatResponseDTO, error) {
	chat := &model.Chat{Title: data.Title}
	response, err := p.repo.Create(chat)
	if err != nil {
		return nil, err
	}

	chatDTO := &model.AddChatResponseDTO{
		ID:        response.ID,
		Title:     response.Title,
		CreatedAt: response.CreatedAt,
	}

	return chatDTO, nil
}

func (p *Chat) GetByID(data *model.GetByIdRequestDTO) (*model.GetChatByIDResponseDTO, error) {
	chat, err := p.repo.GetByID(data.ChatID, data.Limit)
	if err != nil {
		return nil, err
	}

	chatDTO := &model.GetChatByIDResponseDTO{
		ID:        chat.ID,
		Title:     chat.Title,
		Messages:  chat.Messages,
		CreatedAt: chat.CreatedAt,
	}

	return chatDTO, nil
}

func (p *Chat) AddMessageToChat(data *model.AddMessageRequestDTO) (*model.AddMessageResponseDTO, error) {
	message := &model.Message{Text: data.Text, ChatID: data.ChatID}
	dbData, err := p.repo.AddMessageToChat(message)
	if err != nil {
		return nil, err
	}

	messageDTO := &model.AddMessageResponseDTO{
		ID:        dbData.ID,
		Text:      dbData.Text,
		ChatID:    dbData.ChatID,
		CreatedAt: dbData.CreatedAt,
	}

	return messageDTO, nil
}

func (p *Chat) Delete(id uint) error {
	return p.repo.Delete(id)
}
