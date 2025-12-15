package main

import (
	"log"
	"net/http"
	"os"

	"maushold/ranking-service/config"
	"maushold/ranking-service/handler"
	"maushold/ranking-service/messaging"
	"maushold/ranking-service/repository"
	"maushold/ranking-service/routes"
	"maushold/ranking-service/service"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	db := config.InitDB(cfg)
	redisClient := config.InitRedis(cfg)
	rabbitConn, rabbitCh := config.InitRabbitMQ(cfg)
	defer rabbitConn.Close()
	defer rabbitCh.Close()

	consulClient := config.InitConsul(cfg)
	err := config.RegisterService(consulClient, "ranking-service", cfg.ServicePort)
	if err != nil {
		log.Printf("Failed to register with Consul: %v", err)
	}
	defer config.DeregisterService(consulClient, "ranking-service")

	rankingRepo := repository.NewRankingRepository(db)
	playerClient := service.NewPlayerClient(cfg.PlayerServiceURL)
	leaderboardService := service.NewLeaderboardService(redisClient)
	rankingService := service.NewRankingService(rankingRepo, playerClient, leaderboardService)

	// Initialize service discovery
	serviceDiscovery := service.NewServiceDiscovery(consulClient)

	messageConsumer := messaging.NewConsumer(rabbitCh, rankingService)
	go messageConsumer.Start()

	// Start periodic sync
	go rankingService.StartPeriodicSync()

	rankingHandler := handler.NewRankingHandler(rankingService, serviceDiscovery)

	router := mux.NewRouter()
	routes.SetupRankingRoutes(router, rankingHandler)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8004"
	}

	log.Printf("Ranking Service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
