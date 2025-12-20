package routes

import (
	"net/http"

	"maushold/player-service/handler"

	"github.com/gorilla/mux"
)

func SetupPlayerRoutes(router *mux.Router, handler *handler.PlayerHandler) {
	// API routes
	router.HandleFunc("/players/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/players", handler.CreatePlayer).Methods(http.MethodPost)
	router.HandleFunc("/players/{id}", handler.GetPlayer).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}", handler.UpdatePlayer).Methods(http.MethodPut)
	router.HandleFunc("/players/{id}", handler.DeletePlayer).Methods(http.MethodDelete)
	router.HandleFunc("/players", handler.GetAllPlayers).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}/monster", handler.GetPlayerMonster).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}/monster", handler.AddMonsterToPlayer).Methods(http.MethodPost)
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)
}
