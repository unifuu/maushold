package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"maushold/player-service/model"
	"maushold/player-service/testutil"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPlayerService is a mock implementation of PlayerService
type MockPlayerService struct {
	mock.Mock
}

func (m *MockPlayerService) CreatePlayer(player *model.Player) error {
	args := m.Called(player)
	return args.Error(0)
}

func (m *MockPlayerService) GetPlayer(id uint) (*model.Player, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerService) UpdatePlayer(player *model.Player) error {
	args := m.Called(player)
	return args.Error(0)
}

func (m *MockPlayerService) DeletePlayer(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPlayerService) GetAllPlayers() ([]model.Player, error) {
	args := m.Called()
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockPlayerService) UpdatePlayerPoints(id uint, delta int) error {
	args := m.Called(id, delta)
	return args.Error(0)
}

func (m *MockPlayerService) AuthenticatePlayer(username, password string) (*model.Player, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

// MockPlayerMonsterService is a mock implementation of PlayerMonsterService
type MockPlayerMonsterService struct {
	mock.Mock
}

func (m *MockPlayerMonsterService) AddMonsterToPlayer(monster *model.PlayerMonster) error {
	args := m.Called(monster)
	return args.Error(0)
}

func (m *MockPlayerMonsterService) GetPlayerMonster(playerID uint) ([]model.PlayerMonster, error) {
	args := m.Called(playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PlayerMonster), args.Error(1)
}

func TestCreatePlayer(t *testing.T) {
	mockService := new(MockPlayerService)
	mockMonsterService := new(MockPlayerMonsterService)

	// Create handler with nil for unused dependencies in these tests
	handler := &PlayerHandler{
		playerService:        mockService,
		playerMonsterService: mockMonsterService,
		messageProducer:      nil,
		serviceDiscovery:     nil,
	}

	t.Run("successful player creation", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		mockService.On("CreatePlayer", mock.AnythingOfType("*model.Player")).Return(nil).Once()

		req := httptest.NewRequest("POST", "/players", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.CreatePlayer(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response model.Player
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "testuser", response.Username)

		mockService.AssertExpectations(t)
	})

	t.Run("missing username", func(t *testing.T) {
		reqBody := map[string]string{
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/players", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.CreatePlayer(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "required")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/players", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		handler.CreatePlayer(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		mockService.On("CreatePlayer", mock.AnythingOfType("*model.Player")).Return(errors.New("database error")).Once()

		req := httptest.NewRequest("POST", "/players", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.CreatePlayer(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetPlayer(t *testing.T) {
	mockService := new(MockPlayerService)
	mockMonsterService := new(MockPlayerMonsterService)

	handler := &PlayerHandler{
		playerService:        mockService,
		playerMonsterService: mockMonsterService,
		messageProducer:      nil,
		serviceDiscovery:     nil,
	}

	t.Run("successful retrieval", func(t *testing.T) {
		expectedPlayer := testutil.CreateTestPlayer(1, "testuser")
		mockService.On("GetPlayer", uint(1)).Return(expectedPlayer, nil).Once()

		req := httptest.NewRequest("GET", "/players/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetPlayer(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.Player
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedPlayer.Username, response.Username)

		mockService.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		mockService.On("GetPlayer", uint(999)).Return(nil, errors.New("not found")).Once()

		req := httptest.NewRequest("GET", "/players/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		handler.GetPlayer(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid player ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/players/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.GetPlayer(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdatePlayer(t *testing.T) {
	mockService := new(MockPlayerService)
	mockMonsterService := new(MockPlayerMonsterService)

	handler := &PlayerHandler{
		playerService:        mockService,
		playerMonsterService: mockMonsterService,
		messageProducer:      nil,
		serviceDiscovery:     nil,
	}

	t.Run("successful update", func(t *testing.T) {
		existingPlayer := testutil.CreateTestPlayer(1, "oldusername")
		updateBody := map[string]string{
			"username": "newusername",
		}
		body, _ := json.Marshal(updateBody)

		mockService.On("GetPlayer", uint(1)).Return(existingPlayer, nil).Once()
		mockService.On("UpdatePlayer", mock.MatchedBy(func(p *model.Player) bool {
			return p.Username == "newusername"
		})).Return(nil).Once()

		req := httptest.NewRequest("PUT", "/players/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.UpdatePlayer(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		updateBody := map[string]string{
			"username": "newusername",
		}
		body, _ := json.Marshal(updateBody)

		mockService.On("GetPlayer", uint(999)).Return(nil, errors.New("not found")).Once()

		req := httptest.NewRequest("PUT", "/players/999", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		handler.UpdatePlayer(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeletePlayer(t *testing.T) {
	mockService := new(MockPlayerService)
	mockMonsterService := new(MockPlayerMonsterService)

	handler := &PlayerHandler{
		playerService:        mockService,
		playerMonsterService: mockMonsterService,
		messageProducer:      nil,
		serviceDiscovery:     nil,
	}

	t.Run("successful deletion", func(t *testing.T) {
		mockService.On("DeletePlayer", uint(1)).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/players/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DeletePlayer(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Player deleted successfully", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("deletion error", func(t *testing.T) {
		mockService.On("DeletePlayer", uint(1)).Return(errors.New("deletion failed")).Once()

		req := httptest.NewRequest("DELETE", "/players/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DeletePlayer(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockService := new(MockPlayerService)
	mockMonsterService := new(MockPlayerMonsterService)

	handler := &PlayerHandler{
		playerService:        mockService,
		playerMonsterService: mockMonsterService,
		messageProducer:      nil,
		serviceDiscovery:     nil,
	}

	t.Run("successful login", func(t *testing.T) {
		expectedPlayer := testutil.CreateTestPlayer(1, "testuser")
		loginBody := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(loginBody)

		mockService.On("AuthenticatePlayer", "testuser", "password123").Return(expectedPlayer, nil).Once()

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.Player
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "testuser", response.Username)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		loginBody := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(loginBody)

		mockService.On("AuthenticatePlayer", "testuser", "wrongpassword").Return(nil, errors.New("invalid credentials")).Once()

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("missing credentials", func(t *testing.T) {
		loginBody := map[string]string{
			"username": "testuser",
		}
		body, _ := json.Marshal(loginBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAllPlayers(t *testing.T) {
	mockService := new(MockPlayerService)
	mockMonsterService := new(MockPlayerMonsterService)

	handler := &PlayerHandler{
		playerService:        mockService,
		playerMonsterService: mockMonsterService,
		messageProducer:      nil,
		serviceDiscovery:     nil,
	}

	t.Run("successful retrieval", func(t *testing.T) {
		expectedPlayers := []model.Player{
			*testutil.CreateTestPlayer(1, "user1"),
			*testutil.CreateTestPlayer(2, "user2"),
		}

		mockService.On("GetAllPlayers").Return(expectedPlayers, nil).Once()

		req := httptest.NewRequest("GET", "/players", nil)
		w := httptest.NewRecorder()

		handler.GetAllPlayers(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []model.Player
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService.On("GetAllPlayers").Return([]model.Player{}, errors.New("database error")).Once()

		req := httptest.NewRequest("GET", "/players", nil)
		w := httptest.NewRecorder()

		handler.GetAllPlayers(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHealthCheck(t *testing.T) {
	mockService := new(MockPlayerService)
	mockMonsterService := new(MockPlayerMonsterService)

	handler := &PlayerHandler{
		playerService:        mockService,
		playerMonsterService: mockMonsterService,
		messageProducer:      nil,
		serviceDiscovery:     nil,
	}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "player-service", response["service"])
}
