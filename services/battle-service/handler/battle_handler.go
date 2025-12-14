package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"maushold/battle-service/messaging"
	"maushold/battle-service/service"

	"github.com/gorilla/mux"
)

type BattleHandler struct {
	battleService   service.BattleService
	messageProducer *messaging.Producer
}

func NewBattleHandler(battleService service.BattleService, messageProducer *messaging.Producer) *BattleHandler {
	return &BattleHandler{
		battleService:   battleService,
		messageProducer: messageProducer,
	}
}

func (h *BattleHandler) CreateBattle(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Player1ID  uint `json:"player1_id"`
		Player2ID  uint `json:"player2_id"`
		Pokemon1ID uint `json:"monster1_id"`
		Pokemon2ID uint `json:"monster2_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	battle, err := h.battleService.CreateBattle(req.Player1ID, req.Player2ID, req.Pokemon1ID, req.Pokemon2ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.messageProducer.PublishBattleEvent("battle.completed", map[string]interface{}{
		"battle_id":   battle.ID,
		"winner_id":   battle.WinnerID,
		"loser_id":    getLoserID(battle),
		"points_won":  battle.PointsWon,
		"points_lost": battle.PointsLost,
	})

	respondJSON(w, http.StatusCreated, battle)
}

func (h *BattleHandler) GetBattle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid battle ID")
		return
	}

	battle, err := h.battleService.GetBattle(uint(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "Battle not found")
		return
	}

	respondJSON(w, http.StatusOK, battle)
}

func (h *BattleHandler) GetPlayerBattles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["playerId"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	battles, err := h.battleService.GetPlayerBattles(uint(id))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, battles)
}

func (h *BattleHandler) GetAllBattles(w http.ResponseWriter, r *http.Request) {
	battles, err := h.battleService.GetRecentBattles()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, battles)
}

func (h *BattleHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "battle-service"})
}

func getLoserID(battle interface{}) uint {
	// Type assertion to get battle fields
	return 0 // Implement based on your battle struct
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
