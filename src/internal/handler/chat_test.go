package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"chat/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChatProvider struct {
	mock.Mock
}

func (m *MockChatProvider) Create(dto *model.ChatCreateDTO) (*model.Chat, error) {
	args := m.Called(dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chat), args.Error(1)
}

func newRequest(t *testing.T, method, path string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		assert.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req, err := http.NewRequest(method, path, &buf)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestAddChat(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockChatProvider)
		expectedStatus int
		expectedBody   string
		checkResponse  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:        "Успешное создание чата",
			requestBody: &model.ChatCreateDTO{Title: "Test Chat"},
			mockSetup: func(m *MockChatProvider) {
				m.On("Create", &model.ChatCreateDTO{Title: "Test Chat"}).
					Return(&model.Chat{ID: uint(1), Title: "Test Chat"}, nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
				var resp model.Chat
				assert.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
				assert.Equal(t, uint(1), resp.ID)
				assert.Equal(t, "Test Chat", resp.Title)
			},
		},
		{
			name:           "Некорректный JSON в теле запроса",
			requestBody:    "invalid{json",
			mockSetup:      func(m *MockChatProvider) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name:        "Ошибка валидации (пустое поле)",
			requestBody: &model.ChatCreateDTO{Title: ""},
			mockSetup: func(m *MockChatProvider) {
				m.On("Create", mock.Anything).
					Return(nil, assert.AnError) // Или конкретная ошибка валидации
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "validation failed", // Зависит от реализации ошибки в провайдере
		},
		{
			name:        "Ошибка провайдера",
			requestBody: &model.ChatCreateDTO{Title: "Error Chat"},
			mockSetup: func(m *MockChatProvider) {
				m.On("Create", &model.ChatCreateDTO{Title: "Error Chat"}).
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   assert.AnError.Error(),
		},
		{
			name:           "Пустое тело запроса",
			requestBody:    nil,
			mockSetup:      func(m *MockChatProvider) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := new(MockChatProvider)
			if tt.mockSetup != nil {
				tt.mockSetup(mockProvider)
			}

			handler := &Chat{provider: mockProvider}
			req := newRequest(t, http.MethodPost, "/chats", tt.requestBody)
			rec := httptest.NewRecorder()

			handler.AddChat(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code, "Некорректный статус ответа")

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}

			mockProvider.AssertExpectations(t)
		})
	}
}

func TestAddChat_Headers(t *testing.T) {
	mockProvider := new(MockChatProvider)
	mockProvider.On("Create", mock.Anything).
		Return(&model.Chat{ID: 1}, nil)

	handler := &Chat{provider: mockProvider}
	req := newRequest(t, http.MethodPost, "/chats/", &model.ChatCreateDTO{Title: "Test"})
	rec := httptest.NewRecorder()

	handler.AddChat(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type")) // Двойная проверка для надежности
	mockProvider.AssertExpectations(t)
}
