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

func (c *PlayerClient) GetPlayerPokemon(playerID, pokemonID uint) (*model.PlayerPokemon, error) {
	url := fmt.Sprintf("%s/players/%d/pokemon", c.baseURL, playerID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var pokemons []model.PlayerPokemon
	if err := json.Unmarshal(body, &pokemons); err != nil {
		return nil, err
	}

	for _, p := range pokemons {
		if p.ID == pokemonID {
			if p.Nickname == "" {
				p.Nickname = fmt.Sprintf("Pokemon #%d", p.ID)
			}
			return &p, nil
		}
	}

	return nil, fmt.Errorf("pokemon not found")
}
