package routes

import (
	"net/http"

	"maushold/monster-service/handler"

	"github.com/gorilla/mux"
)

func SetupMonsterRoutes(router *mux.Router, handler *handler.MonsterHandler) {
	// API routes
	router.HandleFunc("/monster", handler.CreateMonster).Methods(http.MethodPost)
	router.HandleFunc("/monster/{id}", handler.GetMonster).Methods(http.MethodGet)
	router.HandleFunc("/monster", handler.GetAllMonster).Methods(http.MethodGet)
	router.HandleFunc("/monster/random", handler.GetRandomMonster).Methods(http.MethodGet)
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)
}
