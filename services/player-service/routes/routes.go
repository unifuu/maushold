package routes

import (
	"net/http"

	"maushold/player-service/handler"

	"github.com/gorilla/mux"
)

func SetupPlayerRoutes(router *mux.Router, handler *handler.PlayerHandler) {
	router.HandleFunc("/players", handler.CreatePlayer).Methods(http.MethodPost)
	router.HandleFunc("/players/{id}", handler.GetPlayer).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}", handler.UpdatePlayer).Methods(http.MethodPut)
	router.HandleFunc("/players", handler.GetAllPlayers).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}/pokemon", handler.GetPlayerPokemon).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}/pokemon", handler.AddPokemonToPlayer).Methods(http.MethodPost)
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	// CORS middleware
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
