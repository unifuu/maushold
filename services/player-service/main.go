package main

import (
	"log"
	"net/http"
	"os"

	"maushold/player-service/config"
	"maushold/player-service/handler"
	"maushold/player-service/messaging"
	"maushold/player-service/repository"
	"maushold/player-service/routes"
	"maushold/player-service/service"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := config.InitDB(cfg)

	// Initialize Redis
	redisClient := config.InitRedis(cfg)

	// Initialize RabbitMQ
	rabbitConn, rabbitCh := config.InitRabbitMQ(cfg)
	defer rabbitConn.Close()
	defer rabbitCh.Close()

	// Initialize Consul
	consulClient := config.InitConsul(cfg)

	// Register service with Consul
	err := config.RegisterService(consulClient, "player-service", cfg.ServicePort)
	if err != nil {
		log.Printf("Failed to register with Consul: %v", err)
	}
	defer config.DeregisterService(consulClient, "player-service")

	// Initialize repositories
	playerRepo := repository.NewPlayerRepository(db)
	playerMonsterRepo := repository.NewPlayerMonsterRepository(db)

	// Initialize services
	playerService := service.NewPlayerService(playerRepo, redisClient)
	playerMonsterService := service.NewPlayerMonsterService(playerMonsterRepo, redisClient)

	// Initialize messaging
	messageProducer := messaging.NewProducer(rabbitCh)
	messageConsumer := messaging.NewConsumer(rabbitCh, playerService)

	// Start consuming messages
	go messageConsumer.Start()

	// Initialize handlers
	playerHandler := handler.NewPlayerHandler(playerService, playerMonsterService, messageProducer)

	// Setup routes
	router := mux.NewRouter()
	routes.SetupPlayerRoutes(router, playerHandler)

	// Start server
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("Player Service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
