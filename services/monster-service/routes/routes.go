package routes

import (
	"net/http"

	"maushold/monster-service/handler"

	"github.com/gorilla/mux"
)

func SetupMonsterRoutes(router *mux.Router, handler *handler.MonsterHandler) {
	// CORS middleware - apply first
	router.Use(corsMiddleware)

	// API routes
	router.HandleFunc("/monster", handler.CreateMonster).Methods(http.MethodPost)
	router.HandleFunc("/monster/{id}", handler.GetMonster).Methods(http.MethodGet)
	router.HandleFunc("/monster", handler.GetAllMonster).Methods(http.MethodGet)
	router.HandleFunc("/monster/random", handler.GetRandomMonster).Methods(http.MethodGet)
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
