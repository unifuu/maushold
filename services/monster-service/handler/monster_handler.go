package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"maushold/monster-service/messaging"
	"maushold/monster-service/model"
	"maushold/monster-service/service"

	"github.com/gorilla/mux"
)

type MonsterHandler struct {
	monsterService  service.MonsterService
	messageProducer *messaging.Producer
}

func NewMonsterHandler(monsterService service.MonsterService, messageProducer *messaging.Producer) *MonsterHandler {
	return &MonsterHandler{
		monsterService:  monsterService,
		messageProducer: messageProducer,
	}
}

func (h *MonsterHandler) CreateMonster(w http.ResponseWriter, r *http.Request) {
	var monster model.Monster
	if err := json.NewDecoder(r.Body).Decode(&monster); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.monsterService.CreateMonster(&monster); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.messageProducer.PublishMonsterEvent("monster.created", monster)
	respondJSON(w, http.StatusCreated, monster)
}

func (h *MonsterHandler) GetMonster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid monster ID")
		return
	}

	monster, err := h.monsterService.GetMonster(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Monster not found")
		return
	}

	respondJSON(w, http.StatusOK, monster)
}

func (h *MonsterHandler) GetAllMonster(w http.ResponseWriter, r *http.Request) {
	monster, err := h.monsterService.GetAllMonster()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, monster)
}

func (h *MonsterHandler) GetRandomMonster(w http.ResponseWriter, r *http.Request) {
	monster, err := h.monsterService.GetRandomMonster()
	if err != nil {
		respondError(w, http.StatusNotFound, "No Monster found")
		return
	}

	respondJSON(w, http.StatusOK, monster)
}

func (h *MonsterHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "monster-service"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
