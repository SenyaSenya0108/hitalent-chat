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
	"chat/internal/validation"
)

type Chat struct {
	provider provider.ChatProvider
	validate *validation.Service
}

func NewChatHandler() *Chat {
	return &Chat{
		provider: provider.NewChatProvider(),
		validate: validation.GetValidator(),
	}
}

func (h *Chat) AddChat(w http.ResponseWriter, r *http.Request) {
	request := model.AddChatRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErr, code := h.validate.ValidationHttpRequest(&request)
	if validationErr != nil {
		http.Error(w, validationErr.Error(), code)
		return
	}

	response, err := h.provider.Create(&request)
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
		if errors.Is(err, repository.ErrChatNotFound) {
			http.Error(w, "chat not found", http.StatusNotFound)
			log.Println(errors.New(fmt.Sprint(err)))
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(errors.New(fmt.Sprint(err)))
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

	request := model.AddMessageRequestDTO{ChatID: uint(id)}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		log.Println(errors.New(fmt.Sprint(err)))
		return
	}

	validationErr, code := h.validate.ValidationHttpRequest(&request)
	if validationErr != nil {
		http.Error(w, validationErr.Error(), code)
		return
	}

	response, err := h.provider.AddMessageToChat(&request)
	if err != nil {
		if errors.Is(err, repository.ErrChatNotFound) {
			http.Error(w, "chat not found", http.StatusNotFound)
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
