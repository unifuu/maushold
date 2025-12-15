package routes

import (
	"net/http"

	"maushold/battle-service/handler"

	"github.com/gorilla/mux"
	"github.com/unifuu/lapras"
)

func SetupBattleRoutes(router *mux.Router, handler *handler.BattleHandler) {
	router.Use(lapras.Cors)

	router.HandleFunc("/battles", handler.CreateBattle).Methods(http.MethodPost)
	router.HandleFunc("/battles/{id}", handler.GetBattle).Methods(http.MethodGet)
	router.HandleFunc("/battles/player/{playerId}", handler.GetPlayerBattles).Methods(http.MethodGet)
	router.HandleFunc("/battles", handler.GetAllBattles).Methods(http.MethodGet)
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)
}
