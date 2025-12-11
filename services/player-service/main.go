package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Models
type Player struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Points    int       `gorm:"default:1000" json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlayerPokemon struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PlayerID   uint      `gorm:"not null;index" json:"player_id"`
	PokemonID  int       `gorm:"not null" json:"pokemon_id"`
	Nickname   string    `json:"nickname"`
	Level      int       `gorm:"default:1" json:"level"`
	Experience int       `gorm:"default:0" json:"experience"`
	HP         int       `gorm:"not null" json:"hp"`
	Attack     int       `gorm:"not null" json:"attack"`
	Defense    int       `gorm:"not null" json:"defense"`
	Speed      int       `gorm:"not null" json:"speed"`
	CreatedAt  time.Time `json:"created_at"`
}

// Global variables
var (
	db          *gorm.DB
	redisClient *redis.Client
	rabbitConn  *amqp.Connection
	rabbitCh    *amqp.Channel
	ctx         = context.Background()
)

func main() {
	// Initialize connections
	initDB()
	initRedis()
	initRabbitMQ()
	defer cleanup()

	// Start consuming messages
	go consumeMessages()

	// Setup router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/players", createPlayer).Methods("POST")
	r.HandleFunc("/players/{id}", getPlayer).Methods("GET")
	r.HandleFunc("/players/{id}", updatePlayer).Methods("PUT")
	r.HandleFunc("/players", getAllPlayers).Methods("GET")
	r.HandleFunc("/players/{id}/pokemon", getPlayerPokemon).Methods("GET")
	r.HandleFunc("/players/{id}/pokemon", addPokemonToPlayer).Methods("POST")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// CORS middleware
	r.Use(corsMiddleware)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("Player Service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func initDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	db.AutoMigrate(&Player{}, &PlayerPokemon{})
	log.Println("Database connected and migrated")
}

func initRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisHost,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Redis connected")
}

func initRabbitMQ() {
	var err error
	rabbitURL := os.Getenv("RABBITMQ_URL")

	rabbitConn, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	rabbitCh, err = rabbitConn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}

	// Declare exchanges
	err = rabbitCh.ExchangeDeclare("player.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	// Declare queue for consuming
	_, err = rabbitCh.QueueDeclare("player.updates", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	// Bind queue
	err = rabbitCh.QueueBind("player.updates", "battle.completed", "battle.events", false, nil)
	if err != nil {
		log.Fatal("Failed to bind queue:", err)
	}

	log.Println("RabbitMQ connected")
}

func cleanup() {
	if rabbitCh != nil {
		rabbitCh.Close()
	}
	if rabbitConn != nil {
		rabbitConn.Close()
	}
	if redisClient != nil {
		redisClient.Close()
	}
}

// HTTP Handlers
func createPlayer(w http.ResponseWriter, r *http.Request) {
	var player Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	player.Points = 1000 // Starting points
	if err := db.Create(&player).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Publish player created event
	publishEvent("player.created", player)

	respondJSON(w, http.StatusCreated, player)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Try cache first
	cacheKey := fmt.Sprintf("player:%s", id)
	cached, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var player Player
		json.Unmarshal([]byte(cached), &player)
		respondJSON(w, http.StatusOK, player)
		return
	}

	var player Player
	if err := db.First(&player, id).Error; err != nil {
		respondError(w, http.StatusNotFound, "Player not found")
		return
	}

	// Cache the result
	data, _ := json.Marshal(player)
	redisClient.Set(ctx, cacheKey, data, 5*time.Minute)

	respondJSON(w, http.StatusOK, player)
}

func updatePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var player Player
	if err := db.First(&player, id).Error; err != nil {
		respondError(w, http.StatusNotFound, "Player not found")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := db.Model(&player).Updates(updates).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("player:%s", id)
	redisClient.Del(ctx, cacheKey)

	publishEvent("player.updated", player)

	respondJSON(w, http.StatusOK, player)
}

func getAllPlayers(w http.ResponseWriter, r *http.Request) {
	var players []Player
	if err := db.Order("points DESC").Find(&players).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, players)
}

func getPlayerPokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var pokemon []PlayerPokemon
	if err := db.Where("player_id = ?", id).Find(&pokemon).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, pokemon)
}

func addPokemonToPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var pokemon PlayerPokemon
	if err := json.NewDecoder(r.Body).Decode(&pokemon); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	pokemon.PlayerID = parseUint(id)
	if err := db.Create(&pokemon).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	publishEvent("player.pokemon.added", pokemon)

	respondJSON(w, http.StatusCreated, pokemon)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// Message Consumer
func consumeMessages() {
	msgs, err := rabbitCh.Consume("player.updates", "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	log.Println("Listening for messages...")

	for msg := range msgs {
		log.Printf("Received message: %s", msg.RoutingKey)

		switch msg.RoutingKey {
		case "battle.completed":
			handleBattleCompleted(msg.Body)
		}
	}
}

func handleBattleCompleted(body []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error parsing battle event: %v", err)
		return
	}

	// Update player stats based on battle result
	log.Printf("Battle completed event: %v", event)
}

// Helpers
func publishEvent(routingKey string, data interface{}) {
	body, _ := json.Marshal(data)
	err := rabbitCh.Publish("player.events", routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		log.Printf("Failed to publish event: %v", err)
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func parseUint(s string) uint {
	var result uint
	fmt.Sscanf(s, "%d", &result)
	return result
}
