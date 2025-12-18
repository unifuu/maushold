package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"maushold/ranking-service/service"

	"github.com/gorilla/mux"
)

type RankingHandler struct {
	rankingService   service.RankingService
	serviceDiscovery *service.ServiceDiscovery
}

func NewRankingHandler(rankingService service.RankingService, serviceDiscovery *service.ServiceDiscovery) *RankingHandler {
	return &RankingHandler{
		rankingService:   rankingService,
		serviceDiscovery: serviceDiscovery,
	}
}

// GetLeaderboard returns the top N players with metadata
func (h *RankingHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 10000 {
			limit = parsed
		}
	}

	response, err := h.rankingService.GetLeaderboard(limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// GetPlayerRanking returns a player's ranking information
func (h *RankingHandler) GetPlayerRanking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["playerId"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	ranking, err := h.rankingService.GetPlayerRanking(uint(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "Player ranking not found")
		return
	}

	respondJSON(w, http.StatusOK, ranking)
}

// GetPlayerRankWithContext returns a player's rank with surrounding players
func (h *RankingHandler) GetPlayerRankWithContext(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["playerId"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	contextSize := 5 // Default context size
	if c := r.URL.Query().Get("context"); c != "" {
		if parsed, err := strconv.Atoi(c); err == nil && parsed > 0 && parsed <= 50 {
			contextSize = parsed
		}
	}

	context, err := h.rankingService.GetPlayerRankWithContext(uint(id), contextSize)
	if err != nil {
		respondError(w, http.StatusNotFound, "Player not found")
		return
	}

	respondJSON(w, http.StatusOK, context)
}

// UpdateCombatPower updates a player's combat power (internal endpoint)
func (h *RankingHandler) UpdateCombatPower(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerID    uint  `json:"player_id"`
		CombatPower int64 `json:"combat_power"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PlayerID == 0 || req.CombatPower < 0 {
		respondError(w, http.StatusBadRequest, "Invalid player ID or combat power")
		return
	}

	err := h.rankingService.UpdatePlayerCombatPower(req.PlayerID, req.CombatPower)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Combat power updated successfully"})
}

// SyncRankings triggers a full sync from DB to Redis
func (h *RankingHandler) SyncRankings(w http.ResponseWriter, r *http.Request) {
	go h.rankingService.SyncRankings()
	respondJSON(w, http.StatusAccepted, map[string]string{"message": "Sync started"})
}

// RefreshMaterializedView refreshes the materialized view
func (h *RankingHandler) RefreshMaterializedView(w http.ResponseWriter, r *http.Request) {
	go h.rankingService.RefreshMaterializedView()
	respondJSON(w, http.StatusAccepted, map[string]string{"message": "Materialized view refresh started"})
}

// HealthCheck returns service health status
func (h *RankingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "ranking-service",
		"version": "2.0.0",
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
