package main

import (
	"log"
	"net/http"
	"os"

	"maushold/battle-service/config"
	"maushold/battle-service/handler"
	"maushold/battle-service/messaging"
	"maushold/battle-service/repository"
	"maushold/battle-service/routes"
	"maushold/battle-service/service"

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
	err := config.RegisterService(consulClient, "battle-service", cfg.ServicePort)
	if err != nil {
		log.Printf("Failed to register with Consul: %v", err)
	}
	defer config.DeregisterService(consulClient, "battle-service")

	battleRepo := repository.NewBattleRepository(db)
	playerClient := service.NewPlayerClient(cfg.PlayerServiceURL)
	battleEngine := service.NewBattleEngine()
	battleService := service.NewBattleService(battleRepo, playerClient, battleEngine, redisClient)

	messageProducer := messaging.NewProducer(rabbitCh)
	battleHandler := handler.NewBattleHandler(battleService, messageProducer)

	router := mux.NewRouter()
	routes.SetupBattleRoutes(router, battleHandler)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8003"
	}

	log.Printf("Battle Service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
