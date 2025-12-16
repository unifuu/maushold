package routes

import (
	"net/http"

	"maushold/player-service/handler"

	"github.com/gorilla/mux"
	"github.com/unifuu/lapras"
)

func SetupPlayerRoutes(router *mux.Router, handler *handler.PlayerHandler) {
	router.Use(lapras.Cors)

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
