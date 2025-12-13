package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"maushold/ranking-service/model"
)

type PlayerClient struct {
	baseURL string
}

func NewPlayerClient(baseURL string) *PlayerClient {
	return &PlayerClient{baseURL: baseURL}
}

func (c *PlayerClient) GetPlayer(playerID uint) (*model.Player, error) {
	url := fmt.Sprintf("%s/players/%d", c.baseURL, playerID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var player model.Player
	if err := json.Unmarshal(body, &player); err != nil {
		return nil, err
	}

	return &player, nil
}
