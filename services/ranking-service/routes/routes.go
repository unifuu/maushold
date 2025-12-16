package routes

import (
	"net/http"

	"maushold/ranking-service/handler"

	"github.com/gorilla/mux"
	"github.com/unifuu/lapras"
)

func SetupRankingRoutes(router *mux.Router, handler *handler.RankingHandler) {
	router.Use(lapras.Cors)

	router.HandleFunc("/rankings", handler.GetLeaderboard).Methods(http.MethodGet)
	router.HandleFunc("/rankings/player/{playerId}", handler.GetPlayerRanking).Methods(http.MethodGet)
	router.HandleFunc("/rankings/sync", handler.SyncRankings).Methods(http.MethodPost)
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)
}
