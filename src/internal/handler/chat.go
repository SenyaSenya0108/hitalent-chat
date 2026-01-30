package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"chat/internal/model"
	"chat/internal/provider"
	"chat/internal/repository"

	"gorm.io/gorm"
)

type Chat struct {
	provider provider.ChatProvider
}

func NewChatHandler() *Chat {
	return &Chat{
		provider: provider.NewChatProvider(),
	}
}

func (h *Chat) AddChat(w http.ResponseWriter, r *http.Request) {
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

func (h *Chat) GetByID(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("id")
	id, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat ID", http.StatusBadRequest)
		return
	}

	limit := r.Context().Value("limit").(int)
	request := &model.GetByIdRequestDTO{ChatID: uint(id), Limit: limit}

	chat, err := h.provider.GetByID(request)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "chat not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chat)
}

func (h *Chat) Delete(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("id")
	id, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	err = h.provider.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Chat) AddMessageToChat(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("id")
	id, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat ID", http.StatusBadRequest)
		return
	}

	request := &model.AddMessageDTO{ChatID: uint(id)}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		log.Println(errors.New(fmt.Sprint(err)))
		return
	}

	response, err := h.provider.AddMessageToChat(request)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			http.Error(w, "chatID error", http.StatusBadRequest)
			log.Println(errors.New(fmt.Sprint(err)))
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println(errors.New(fmt.Sprint(err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
