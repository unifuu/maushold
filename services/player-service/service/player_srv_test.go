package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"maushold/player-service/model"
	"maushold/player-service/testutil"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockPlayerRepository is a mock implementation of PlayerRepository
type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) Create(player *model.Player) error {
	args := m.Called(player)
	return args.Error(0)
}

func (m *MockPlayerRepository) FindByID(id uint) (*model.Player, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepository) Update(player *model.Player) error {
	args := m.Called(player)
	return args.Error(0)
}

func (m *MockPlayerRepository) FindAll() ([]model.Player, error) {
	args := m.Called()
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockPlayerRepository) UpdatePoints(id uint, points int) error {
	args := m.Called(id, points)
	return args.Error(0)
}

func (m *MockPlayerRepository) Delete(player *model.Player) error {
	args := m.Called(player)
	return args.Error(0)
}

// MockRedisClient is a simple mock for Redis operations
type MockRedisClient struct {
	data map[string]string
	ctx  context.Context
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
		ctx:  context.Background(),
	}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	if val, ok := m.data[key]; ok {
		cmd.SetVal(val)
	} else {
		cmd.SetErr(redis.Nil)
	}
	return cmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	// Handle both string and []byte values
	switch v := value.(type) {
	case string:
		m.data[key] = v
	case []byte:
		m.data[key] = string(v)
	default:
		m.data[key] = fmt.Sprintf("%v", v)
	}
	cmd.SetVal("OK")
	return cmd
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	count := int64(0)
	for _, key := range keys {
		if _, ok := m.data[key]; ok {
			delete(m.data, key)
			count++
		}
	}
	cmd.SetVal(count)
	return cmd
}

func TestCreatePlayer(t *testing.T) {
	mockRepo := new(MockPlayerRepository)
	mockRedis := NewMockRedisClient()
	service := NewPlayerService(mockRepo, mockRedis)

	t.Run("successful player creation", func(t *testing.T) {
		player := &model.Player{
			Username: "testuser",
			Password: "password123",
		}

		mockRepo.On("Create", mock.AnythingOfType("*model.Player")).Return(nil).Once()

		err := service.CreatePlayer(player)

		assert.NoError(t, err)
		assert.Equal(t, 0, player.Points)
		// Verify password was hashed
		assert.NotEqual(t, "password123", player.Password)
		err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte("password123"))
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		player := &model.Player{
			Username: "testuser",
			Password: "password123",
		}

		mockRepo.On("Create", mock.AnythingOfType("*model.Player")).Return(errors.New("database error")).Once()

		err := service.CreatePlayer(player)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPlayer(t *testing.T) {
	mockRepo := new(MockPlayerRepository)
	mockRedis := NewMockRedisClient()
	service := NewPlayerService(mockRepo, mockRedis)

	t.Run("cache hit", func(t *testing.T) {
		expectedPlayer := testutil.CreateTestPlayer(1, "testuser")
		cachedData, _ := json.Marshal(expectedPlayer)
		mockRedis.data["player:1"] = string(cachedData)

		player, err := service.GetPlayer(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedPlayer.ID, player.ID)
		assert.Equal(t, expectedPlayer.Username, player.Username)
		// Repository should not be called on cache hit
		mockRepo.AssertNotCalled(t, "FindByID")
	})

	t.Run("cache miss - fetch from database", func(t *testing.T) {
		expectedPlayer := testutil.CreateTestPlayer(2, "testuser2")
		mockRepo.On("FindByID", uint(2)).Return(expectedPlayer, nil).Once()

		player, err := service.GetPlayer(2)

		assert.NoError(t, err)
		assert.Equal(t, expectedPlayer.ID, player.ID)
		assert.Equal(t, expectedPlayer.Username, player.Username)

		// Verify data was cached
		cachedData, exists := mockRedis.data["player:2"]
		assert.True(t, exists)

		var cachedPlayer model.Player
		json.Unmarshal([]byte(cachedData), &cachedPlayer)
		assert.Equal(t, expectedPlayer.ID, cachedPlayer.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("record not found")).Once()

		player, err := service.GetPlayer(999)

		assert.Error(t, err)
		assert.Nil(t, player)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdatePlayer(t *testing.T) {
	mockRepo := new(MockPlayerRepository)
	mockRedis := NewMockRedisClient()
	service := NewPlayerService(mockRepo, mockRedis)

	t.Run("successful update with cache invalidation", func(t *testing.T) {
		player := testutil.CreateTestPlayer(1, "updateduser")

		// Pre-populate cache
		cachedData, _ := json.Marshal(player)
		mockRedis.data["player:1"] = string(cachedData)

		mockRepo.On("Update", player).Return(nil).Once()

		err := service.UpdatePlayer(player)

		assert.NoError(t, err)
		// Verify cache was invalidated
		_, exists := mockRedis.data["player:1"]
		assert.False(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		player := testutil.CreateTestPlayer(1, "testuser")
		mockRepo.On("Update", player).Return(errors.New("update failed")).Once()

		err := service.UpdatePlayer(player)

		assert.Error(t, err)
		assert.Equal(t, "update failed", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestDeletePlayer(t *testing.T) {

	t.Run("successful deletion", func(t *testing.T) {
		player := testutil.CreateTestPlayer(1, "testuser")

		// Pre-populate cache
		mockRepo := new(MockPlayerRepository)
		mockRedis := NewMockRedisClient()
		service := NewPlayerService(mockRepo, mockRedis)

		cachedData, _ := json.Marshal(player)
		mockRedis.data["player:1"] = string(cachedData)

		mockRepo.On("FindByID", uint(1)).Return(player, nil).Once()
		mockRepo.On("Delete", mock.MatchedBy(func(p *model.Player) bool {
			return p.ID == 1 && p.Username == "testuser"
		})).Return(nil).Once()

		err := service.DeletePlayer(1)

		assert.NoError(t, err)
		// Verify cache was invalidated
		_, exists := mockRedis.data["player:1"]
		assert.False(t, exists)
		mockRepo.AssertExpectations(t)
	})

		mockRepo := new(MockPlayerRepository)
		mockRedis := NewMockRedisClient()
		service := NewPlayerService(mockRepo, mockRedis)

	t.Run("player not found", func(t *testing.T) {
		mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found")).Once()

		err := service.DeletePlayer(999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllPlayers(t *testing.T) {
	mockRepo := new(MockPlayerRepository)
	mockRedis := NewMockRedisClient()
	service := NewPlayerService(mockRepo, mockRedis)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedPlayers := []model.Player{
			*testutil.CreateTestPlayer(1, "user1"),
			*testutil.CreateTestPlayer(2, "user2"),
		}

		mockRepo.On("FindAll").Return(expectedPlayers, nil).Once()

		players, err := service.GetAllPlayers()

		assert.NoError(t, err)
		assert.Len(t, players, 2)
		assert.Equal(t, expectedPlayers[0].Username, players[0].Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo.On("FindAll").Return([]model.Player{}, nil).Once()

		players, err := service.GetAllPlayers()

		assert.NoError(t, err)
		assert.Len(t, players, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdatePlayerPoints(t *testing.T) {
	mockRepo := new(MockPlayerRepository)
	mockRedis := NewMockRedisClient()
	service := NewPlayerService(mockRepo, mockRedis)

	t.Run("successful points update", func(t *testing.T) {
		player := testutil.CreateTestPlayer(1, "testuser")
		player.Points = 100

		mockRepo.On("FindByID", uint(1)).Return(player, nil).Once()
		mockRepo.On("Update", mock.MatchedBy(func(p *model.Player) bool {
			return p.ID == 1 && p.Points == 150
		})).Return(nil).Once()

		err := service.UpdatePlayerPoints(1, 50)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("negative points delta", func(t *testing.T) {
		player := testutil.CreateTestPlayer(1, "testuser")
		player.Points = 100

		mockRepo.On("FindByID", uint(1)).Return(player, nil).Once()
		mockRepo.On("Update", mock.MatchedBy(func(p *model.Player) bool {
			return p.ID == 1 && p.Points == 70
		})).Return(nil).Once()

		err := service.UpdatePlayerPoints(1, -30)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthenticatePlayer(t *testing.T) {
	mockRepo := new(MockPlayerRepository)
	mockRedis := NewMockRedisClient()
	service := NewPlayerService(mockRepo, mockRedis)

	t.Run("successful authentication", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		player := &model.Player{
			ID:       1,
			Username: "testuser",
			Password: string(hashedPassword),
			Points:   100,
		}

		mockRepo.On("FindAll").Return([]model.Player{*player}, nil).Once()

		result, err := service.AuthenticatePlayer("testuser", "correctpassword")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "testuser", result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid username", func(t *testing.T) {
		mockRepo.On("FindAll").Return([]model.Player{}, nil).Once()

		result, err := service.AuthenticatePlayer("nonexistent", "password")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		player := &model.Player{
			ID:       1,
			Username: "testuser",
			Password: string(hashedPassword),
			Points:   100,
		}

		mockRepo.On("FindAll").Return([]model.Player{*player}, nil).Once()

		result, err := service.AuthenticatePlayer("testuser", "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("FindAll").Return([]model.Player{}, errors.New("database error")).Once()

		result, err := service.AuthenticatePlayer("testuser", "password")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
