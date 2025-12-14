package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"maushold/battle-service/model"
)

type PlayerClient struct {
	baseURL string
}

func NewPlayerClient(baseURL string) *PlayerClient {
	return &PlayerClient{baseURL: baseURL}
}

func (c *PlayerClient) GetPlayerMonster(playerID, monsterID uint) (*model.PlayerMonster, error) {
	url := fmt.Sprintf("%s/players/%d/monster", c.baseURL, playerID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var monsters []model.PlayerMonster
	if err := json.Unmarshal(body, &monsters); err != nil {
		return nil, err
	}

	for _, p := range monsters {
		if p.ID == monsterID {
			if p.Nickname == "" {
				p.Nickname = fmt.Sprintf("Monster #%d", p.ID)
			}
			return &p, nil
		}
	}

	return nil, fmt.Errorf("monster not found")
}
