package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"maushold/pokemon-service/messaging"
	"maushold/pokemon-service/model"
	"maushold/pokemon-service/service"

	"github.com/gorilla/mux"
)

type PokemonHandler struct {
	pokemonService  service.PokemonService
	messageProducer *messaging.Producer
}

func NewPokemonHandler(pokemonService service.PokemonService, messageProducer *messaging.Producer) *PokemonHandler {
	return &PokemonHandler{
		pokemonService:  pokemonService,
		messageProducer: messageProducer,
	}
}

func (h *PokemonHandler) CreatePokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon model.Pokemon
	if err := json.NewDecoder(r.Body).Decode(&pokemon); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.pokemonService.CreatePokemon(&pokemon); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.messageProducer.PublishPokemonEvent("pokemon.created", pokemon)
	respondJSON(w, http.StatusCreated, pokemon)
}

func (h *PokemonHandler) GetPokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid pokemon ID")
		return
	}

	pokemon, err := h.pokemonService.GetPokemon(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Pokemon not found")
		return
	}

	respondJSON(w, http.StatusOK, pokemon)
}

func (h *PokemonHandler) GetAllPokemon(w http.ResponseWriter, r *http.Request) {
	pokemon, err := h.pokemonService.GetAllPokemon()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, pokemon)
}

func (h *PokemonHandler) GetRandomPokemon(w http.ResponseWriter, r *http.Request) {
	pokemon, err := h.pokemonService.GetRandomPokemon()
	if err != nil {
		respondError(w, http.StatusNotFound, "No Pokemon found")
		return
	}

	respondJSON(w, http.StatusOK, pokemon)
}

func (h *PokemonHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "pokemon-service"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
