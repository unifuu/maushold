package routes

import (
	"net/http"

	"maushold/ranking-service/handler"

	"github.com/gorilla/mux"
)

func SetupRankingRoutes(router *mux.Router, handler *handler.RankingHandler) {
	router.HandleFunc("/rankings", handler.GetLeaderboard).Methods(http.MethodGet)
	router.HandleFunc("/rankings/player/{playerId}", handler.GetPlayerRanking).Methods(http.MethodGet)
	router.HandleFunc("/rankings/sync", handler.SyncRankings).Methods(http.MethodPost)
	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

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
