package main

import (
	"log"
	"net/http"
	"os"

	"maushold/monster-service/config"
	"maushold/monster-service/handler"
	"maushold/monster-service/messaging"
	"maushold/monster-service/repository"
	"maushold/monster-service/routes"
	"maushold/monster-service/service"

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
	err := config.RegisterService(consulClient, "monster-service", cfg.ServicePort)
	if err != nil {
		log.Printf("Failed to register with Consul: %v", err)
	}
	defer config.DeregisterService(consulClient, "monster-service")

	monsterRepo := repository.NewPokemonRepository(db)
	monsterService := service.NewPokemonService(monsterRepo, redisClient)

	// Seed initial data
	service.SeedPokemon(monsterRepo)

	messageProducer := messaging.NewProducer(rabbitCh)
	monsterHandler := handler.NewPokemonHandler(monsterService, messageProducer)

	router := mux.NewRouter()
	routes.SetupPokemonRoutes(router, monsterHandler)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("Pokemon Service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
