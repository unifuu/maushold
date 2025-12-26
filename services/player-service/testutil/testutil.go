package testutil

import (
	"context"
	"time"

	"maushold/player-service/model"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

// MockRedisClient is a mock implementation of Redis client for testing
type MockRedisClient struct {
	Data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		Data: make(map[string]string),
	}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	if val, ok := m.Data[key]; ok {
		cmd.SetVal(val)
	} else {
		cmd.SetErr(redis.Nil)
	}
	return cmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	m.Data[key] = value.(string)
	cmd.SetVal("OK")
	return cmd
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	count := int64(0)
	for _, key := range keys {
		if _, ok := m.Data[key]; ok {
			delete(m.Data, key)
			count++
		}
	}
	cmd.SetVal(count)
	return cmd
}

// MockRabbitMQChannel is a mock implementation of RabbitMQ channel
type MockRabbitMQChannel struct {
	PublishedMessages []amqp.Publishing
}

func NewMockRabbitMQChannel() *MockRabbitMQChannel {
	return &MockRabbitMQChannel{
		PublishedMessages: make([]amqp.Publishing, 0),
	}
}

func (m *MockRabbitMQChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	m.PublishedMessages = append(m.PublishedMessages, msg)
	return nil
}

// Test fixtures
func CreateTestPlayer(id uint, username string) *model.Player {
	return &model.Player{
		ID:        id,
		Username:  username,
		Password:  "hashedpassword123",
		Points:    100,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func CreateTestPlayerMonster(id, playerID uint, monsterID int) *model.PlayerMonster {
	return &model.PlayerMonster{
		ID:         id,
		PlayerID:   playerID,
		MonsterID:  monsterID,
		Nickname:   "TestMonster",
		Level:      5,
		Experience: 100,
		HP:         50,
		Attack:     30,
		Defense:    20,
		Speed:      25,
		CreatedAt:  time.Now(),
	}
}
