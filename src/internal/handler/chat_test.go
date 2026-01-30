package handler

import (
	"errors"
	"net/http/httptest"
	"testing"

	"chat/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockChatProvider struct {
	deleteFn func(uint) error
}

func (m *mockChatProvider) Delete(id uint) error {
	return m.deleteFn(id)
}

func (m *mockChatProvider) Create(
	*model.AddChatRequestDTO,
) (*model.AddChatResponseDTO, error) {
	panic("Create should not be called in Delete tests")
}

func (m *mockChatProvider) GetByID(
	*model.GetByIdRequestDTO,
) (*model.GetChatByIDResponseDTO, error) {
	panic("GetByID should not be called in Delete tests")
}

func (m *mockChatProvider) AddMessageToChat(
	*model.AddMessageRequestDTO,
) (*model.AddMessageResponseDTO, error) {
	panic("AddMessageToChat should not be called in Delete tests")
}

func TestChat_Delete_Success(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/chats/10", nil)
	req.SetPathValue("id", "10")

	rec := httptest.NewRecorder()

	provider := &mockChatProvider{
		deleteFn: func(id uint) error {
			require.Equal(t, uint(10), id)
			return nil
		},
	}

	h := &Chat{
		provider: provider,
	}
	h.Delete(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, 204, res.StatusCode)
}

func TestChat_Delete_InvalidID(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/chats/abc", nil)
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()

	provider := &mockChatProvider{
		deleteFn: func(id uint) error {
			t.Fatal("Delete must not be called when id is invalid")
			return nil
		},
	}

	h := &Chat{
		provider: provider,
	}

	h.Delete(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, 400, res.StatusCode)
}

func TestChat_Delete_ProviderError(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/chats/5", nil)
	req.SetPathValue("id", "5")

	rec := httptest.NewRecorder()

	provider := &mockChatProvider{
		deleteFn: func(id uint) error {
			return errors.New("db error")
		},
	}

	h := &Chat{
		provider: provider,
	}

	h.Delete(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, 500, res.StatusCode)
}
