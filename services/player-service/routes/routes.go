package routes

import (
	"net/http"

	"maushold/player-service/handler"

	"github.com/gorilla/mux"
)

func SetupPlayerRoutes(router *mux.Router, handler *handler.PlayerHandler) {
	// CORS middleware - apply first
	router.Use(corsMiddleware)

	// API routes
	router.HandleFunc("/players", handler.CreatePlayer).Methods(http.MethodPost)
	router.HandleFunc("/players/{id}", handler.GetPlayer).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}", handler.UpdatePlayer).Methods(http.MethodPut)
	router.HandleFunc("/players", handler.GetAllPlayers).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}/monster", handler.GetPlayerMonster).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}/monster", handler.AddMonsterToPlayer).Methods(http.MethodPost)
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
