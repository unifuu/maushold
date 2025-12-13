package routes

import (
	"net/http"

	"maushold/battle-service/handler"

	"github.com/gorilla/mux"
)

func SetupBattleRoutes(router *mux.Router, handler *handler.BattleHandler) {
	router.HandleFunc("/battles", handler.CreateBattle).Methods(http.MethodPost)
	router.HandleFunc("/battles/{id}", handler.GetBattle).Methods(http.MethodGet)
	router.HandleFunc("/battles/player/{playerId}", handler.GetPlayerBattles).Methods(http.MethodGet)
	router.HandleFunc("/battles", handler.GetAllBattles).Methods(http.MethodGet)
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
