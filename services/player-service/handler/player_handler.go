package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"maushold/player-service/messaging"
	"maushold/player-service/model"
	"maushold/player-service/service"

	"github.com/gorilla/mux"
)

type PlayerHandler struct {
	playerService        service.PlayerService
	playerMonsterService service.PlayerMonsterService
	messageProducer      *messaging.Producer
	serviceDiscovery     *service.ServiceDiscovery
}

func NewPlayerHandler(
	playerService service.PlayerService,
	playerMonsterService service.PlayerMonsterService,
	messageProducer *messaging.Producer,
	serviceDiscovery *service.ServiceDiscovery,
) *PlayerHandler {
	return &PlayerHandler{
		playerService:        playerService,
		playerMonsterService: playerMonsterService,
		messageProducer:      messageProducer,
		serviceDiscovery:     serviceDiscovery,
	}
}

func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate inputs
	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	player := model.Player{
		Username: req.Username,
		Password: req.Password, // Will be hashed in service layer
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

func (h *PlayerHandler) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	if err := h.playerService.DeletePlayer(uint(id)); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.messageProducer.PublishPlayerEvent("player.deleted", map[string]interface{}{"id": id})

	respondJSON(w, http.StatusOK, map[string]string{"message": "Player deleted successfully"})
}

func (h *PlayerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate inputs
	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	player, err := h.playerService.AuthenticatePlayer(req.Username, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	respondJSON(w, http.StatusOK, player)
}

func (h *PlayerHandler) GetMonsterInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	monsterID := vars["monsterId"]

	monsterServiceURL, err := h.serviceDiscovery.DiscoverService("monster-service")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Make request to discovered service
	resp, err := http.Get(monsterServiceURL + "/monster/" + monsterID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to call monster service")
		return
	}
	defer resp.Body.Close()

	var monster interface{}
	if err := json.NewDecoder(resp.Body).Decode(&monster); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to parse monster data")
		return
	}

	respondJSON(w, http.StatusOK, monster)
}

func (h *PlayerHandler) GetPlayerMonster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	monster, err := h.playerMonsterService.GetPlayerMonster(uint(id))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, monster)
}

func (h *PlayerHandler) AddMonsterToPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	var monster model.PlayerMonster
	if err := json.NewDecoder(r.Body).Decode(&monster); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	monster.PlayerID = uint(id)
	if err := h.playerMonsterService.AddMonsterToPlayer(&monster); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.messageProducer.PublishPlayerEvent("player.monster.added", monster)

	respondJSON(w, http.StatusCreated, monster)
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
