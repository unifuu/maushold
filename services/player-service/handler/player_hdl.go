package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"player-service/messaging"
	"player-service/model"
	"player-service/service"

	"github.com/gorilla/mux"
)

type PlayerHandler struct {
	playerService        service.PlayerService
	playerPokemonService service.PlayerPokemonService
	messageProducer      *messaging.Producer
}

func NewPlayerHandler(
	playerService service.PlayerService,
	playerPokemonService service.PlayerPokemonService,
	messageProducer *messaging.Producer,
) *PlayerHandler {
	return &PlayerHandler{
		playerService:        playerService,
		playerPokemonService: playerPokemonService,
		messageProducer:      messageProducer,
	}
}

func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var player model.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.playerService.CreatePlayer(&player); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Publish event
	h.messageProducer.PublishPlayerEvent("player.created", player)

	respondJSON(w, http.StatusCreated, player)
}

func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	player, err := h.playerService.GetPlayer(uint(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "Player not found")
		return
	}

	respondJSON(w, http.StatusOK, player)
}

func (h *PlayerHandler) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	player, err := h.playerService.GetPlayer(uint(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "Player not found")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Apply updates
	if username, ok := updates["username"].(string); ok {
		player.Username = username
	}

	if err := h.playerService.UpdatePlayer(player); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.messageProducer.PublishPlayerEvent("player.updated", player)

	respondJSON(w, http.StatusOK, player)
}

func (h *PlayerHandler) GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := h.playerService.GetAllPlayers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, players)
}

func (h *PlayerHandler) GetPlayerPokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	pokemon, err := h.playerPokemonService.GetPlayerPokemon(uint(id))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, pokemon)
}

func (h *PlayerHandler) AddPokemonToPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	var pokemon model.PlayerPokemon
	if err := json.NewDecoder(r.Body).Decode(&pokemon); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	pokemon.PlayerID = uint(id)
	if err := h.playerPokemonService.AddPokemonToPlayer(&pokemon); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.messageProducer.PublishPlayerEvent("player.pokemon.added", pokemon)

	respondJSON(w, http.StatusCreated, pokemon)
}

func (h *PlayerHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "player-service"})
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
