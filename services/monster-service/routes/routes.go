package routes

import (
	"net/http"

	"maushold/monster-service/handler"

	"github.com/gorilla/mux"
	"github.com/unifuu/lapras"
)

func SetupMonsterRoutes(router *mux.Router, handler *handler.MonsterHandler) {
	router.Use(lapras.Cors)

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
