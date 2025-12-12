package main

import (
	"log"
	"net/http"
	"os"

	"maushold/pokemon-service/config"
	"maushold/pokemon-service/handler"
	"maushold/pokemon-service/messaging"
	"maushold/pokemon-service/repository"
	"maushold/pokemon-service/routes"
	"maushold/pokemon-service/service"

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
	err := config.RegisterService(consulClient, "pokemon-service", cfg.ServicePort)
	if err != nil {
		log.Printf("Failed to register with Consul: %v", err)
	}
	defer config.DeregisterService(consulClient, "pokemon-service")

	pokemonRepo := repository.NewPokemonRepository(db)
	pokemonService := service.NewPokemonService(pokemonRepo, redisClient)

	// Seed initial data
	service.SeedPokemon(pokemonRepo)

	messageProducer := messaging.NewProducer(rabbitCh)
	pokemonHandler := handler.NewPokemonHandler(pokemonService, messageProducer)

	router := mux.NewRouter()
	routes.SetupPokemonRoutes(router, pokemonHandler)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("Pokemon Service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
