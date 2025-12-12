package routes

import (
	"net/http"

	"maushold/pokemon-service/handler"

	"github.com/gorilla/mux"
)

func SetupPokemonRoutes(router *mux.Router, handler *handler.PokemonHandler) {
	router.HandleFunc("/pokemon", handler.CreatePokemon).Methods(http.MethodPost)
	router.HandleFunc("/pokemon/{id}", handler.GetPokemon).Methods(http.MethodGet)
	router.HandleFunc("/pokemon", handler.GetAllPokemon).Methods(http.MethodGet)
	router.HandleFunc("/pokemon/random", handler.GetRandomPokemon).Methods(http.MethodGet)
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	router.Use(corsMiddleware)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
