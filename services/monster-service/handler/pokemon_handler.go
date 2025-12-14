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

type PokemonHandler struct {
	monsterService  service.PokemonService
	messageProducer *messaging.Producer
}

func NewPokemonHandler(monsterService service.PokemonService, messageProducer *messaging.Producer) *PokemonHandler {
	return &PokemonHandler{
		monsterService:  monsterService,
		messageProducer: messageProducer,
	}
}

func (h *PokemonHandler) CreatePokemon(w http.ResponseWriter, r *http.Request) {
	var monster model.Pokemon
	if err := json.NewDecoder(r.Body).Decode(&monster); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.monsterService.CreatePokemon(&monster); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.messageProducer.PublishPokemonEvent("monster.created", monster)
	respondJSON(w, http.StatusCreated, monster)
}

func (h *PokemonHandler) GetPokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid monster ID")
		return
	}

	monster, err := h.monsterService.GetPokemon(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Pokemon not found")
		return
	}

	respondJSON(w, http.StatusOK, monster)
}

func (h *PokemonHandler) GetAllPokemon(w http.ResponseWriter, r *http.Request) {
	monster, err := h.monsterService.GetAllPokemon()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, monster)
}

func (h *PokemonHandler) GetRandomPokemon(w http.ResponseWriter, r *http.Request) {
	monster, err := h.monsterService.GetRandomPokemon()
	if err != nil {
		respondError(w, http.StatusNotFound, "No Pokemon found")
		return
	}

	respondJSON(w, http.StatusOK, monster)
}

func (h *PokemonHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
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
