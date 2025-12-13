package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"maushold/ranking-service/service"

	"github.com/gorilla/mux"
)

type RankingHandler struct {
	rankingService service.RankingService
}

func NewRankingHandler(rankingService service.RankingService) *RankingHandler {
	return &RankingHandler{
		rankingService: rankingService,
	}
}

func (h *RankingHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	leaderboard, err := h.rankingService.GetLeaderboard(limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, leaderboard)
}

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

func (h *RankingHandler) SyncRankings(w http.ResponseWriter, r *http.Request) {
	go h.rankingService.SyncRankings()
	respondJSON(w, http.StatusOK, map[string]string{"message": "Sync started"})
}

func (h *RankingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "ranking-service"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
