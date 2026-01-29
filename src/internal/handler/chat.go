package handler

import (
	"encoding/json"
	"net/http"

	"chat/internal/model"
	"chat/internal/provider"
)

type ChatHandler struct {
	provider provider.ChatProvider
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		provider: provider.NewChatProvider(),
	}
}

func (h *ChatHandler) AddChat(w http.ResponseWriter, r *http.Request) {
	request := &model.ChatCreateDTO{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.provider.Create(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ChatHandler) AddMessageToChat(w http.ResponseWriter, r *http.Request) {

}

func (h *ChatHandler) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (h *ChatHandler) Delete(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
