package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"maushold/ranking-service/model"
)

type BattleClient struct {
	baseURL string
}

func NewBattleClient(baseURL string) *BattleClient {
	return &BattleClient{baseURL: baseURL}
}

func (c *BattleClient) GetAllBattles() ([]model.Battle, error) {
	url := fmt.Sprintf("%s/battles", c.baseURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var battles []model.Battle
	if err := json.Unmarshal(body, &battles); err != nil {
		return nil, err
	}

	return battles, nil
}
