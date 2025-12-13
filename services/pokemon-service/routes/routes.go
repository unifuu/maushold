package routes

import (
	"net/http"

	"maushold/pokemon-service/handler"

	"github.com/gorilla/mux"
)

func SetupPokemonRoutes(router *mux.Router, handler *handler.PokemonHandler) {
	// CORS middleware - apply first
	router.Use(corsMiddleware)

	// API routes
	router.HandleFunc("/pokemon", handler.CreatePokemon).Methods(http.MethodPost)
	router.HandleFunc("/pokemon/{id}", handler.GetPokemon).Methods(http.MethodGet)
	router.HandleFunc("/pokemon", handler.GetAllPokemon).Methods(http.MethodGet)
	router.HandleFunc("/pokemon/random", handler.GetRandomPokemon).Methods(http.MethodGet)
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	// Handle preflight requests for all routes
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}).Methods(http.MethodOptions)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
