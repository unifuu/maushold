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
type Pokemon struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Type1       string    `gorm:"not null" json:"type1"`
	Type2       string    `json:"type2"`
	BaseHP      int       `gorm:"not null" json:"base_hp"`
	BaseAttack  int       `gorm:"not null" json:"base_attack"`
	BaseDefense int       `gorm:"not null" json:"base_defense"`
	BaseSpeed   int       `gorm:"not null" json:"base_speed"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
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
	initDB()
	initRedis()
	initRabbitMQ()
	defer cleanup()

	// Seed initial Pokemon data
	seedPokemon()

	r := mux.NewRouter()

	r.HandleFunc("/pokemon", getAllPokemon).Methods("GET")
	r.HandleFunc("/pokemon/{id}", getPokemon).Methods("GET")
	r.HandleFunc("/pokemon", createPokemon).Methods("POST")
	r.HandleFunc("/pokemon/random", getRandomPokemon).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	r.Use(corsMiddleware)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("Pokemon Service starting on port %s", port)
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

	db.AutoMigrate(&Pokemon{})
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

	err = rabbitCh.ExchangeDeclare("pokemon.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
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

func seedPokemon() {
	var count int64
	db.Model(&Pokemon{}).Count(&count)
	if count > 0 {
		return
	}

	starterPokemon := []Pokemon{
		{ID: 1, Name: "Bulbasaur", Type1: "Grass", Type2: "Poison", BaseHP: 45, BaseAttack: 49, BaseDefense: 49, BaseSpeed: 45, Description: "A strange seed was planted on its back at birth."},
		{ID: 4, Name: "Charmander", Type1: "Fire", Type2: "", BaseHP: 39, BaseAttack: 52, BaseDefense: 43, BaseSpeed: 65, Description: "Obviously prefers hot places. When it rains, steam is said to spout from the tip of its tail."},
		{ID: 7, Name: "Squirtle", Type1: "Water", Type2: "", BaseHP: 44, BaseAttack: 48, BaseDefense: 65, BaseSpeed: 43, Description: "After birth, its back swells and hardens into a shell."},
		{ID: 25, Name: "Pikachu", Type1: "Electric", Type2: "", BaseHP: 35, BaseAttack: 55, BaseDefense: 40, BaseSpeed: 90, Description: "When several of these Pokémon gather, their electricity could build and cause lightning storms."},
		{ID: 39, Name: "Jigglypuff", Type1: "Normal", Type2: "Fairy", BaseHP: 115, BaseAttack: 45, BaseDefense: 20, BaseSpeed: 20, Description: "When its huge eyes light up, it sings a mysteriously soothing melody that lulls its enemies to sleep."},
		{ID: 133, Name: "Eevee", Type1: "Normal", Type2: "", BaseHP: 55, BaseAttack: 55, BaseDefense: 50, BaseSpeed: 55, Description: "Its genetic code is irregular. It may mutate if exposed to radiation from element stones."},
		{ID: 143, Name: "Snorlax", Type1: "Normal", Type2: "", BaseHP: 160, BaseAttack: 110, BaseDefense: 65, BaseSpeed: 30, Description: "Very lazy. Just eats and sleeps. As its rotund bulk builds, it becomes steadily more slothful."},
		{ID: 150, Name: "Mewtwo", Type1: "Psychic", Type2: "", BaseHP: 106, BaseAttack: 110, BaseDefense: 90, BaseSpeed: 130, Description: "It was created by a scientist after years of horrific gene splicing and DNA engineering experiments."},
		{ID: 94, Name: "Gengar", Type1: "Ghost", Type2: "Poison", BaseHP: 60, BaseAttack: 65, BaseDefense: 60, BaseSpeed: 110, Description: "Under a full moon, this Pokémon likes to mimic the shadows of people and laugh at their fright."},
		{ID: 6, Name: "Charizard", Type1: "Fire", Type2: "Flying", BaseHP: 78, BaseAttack: 84, BaseDefense: 78, BaseSpeed: 100, Description: "Spits fire that is hot enough to melt boulders. Known to cause forest fires unintentionally."},
	}

	for _, p := range starterPokemon {
		db.Create(&p)
	}

	log.Println("Seeded initial Pokemon data")
}

// HTTP Handlers
func getAllPokemon(w http.ResponseWriter, r *http.Request) {
	// Try cache first
	cached, err := redisClient.Get(ctx, "pokemon:all").Result()
	if err == nil {
		var pokemon []Pokemon
		json.Unmarshal([]byte(cached), &pokemon)
		respondJSON(w, http.StatusOK, pokemon)
		return
	}

	var pokemon []Pokemon
	if err := db.Find(&pokemon).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Cache for 10 minutes
	data, _ := json.Marshal(pokemon)
	redisClient.Set(ctx, "pokemon:all", data, 10*time.Minute)

	respondJSON(w, http.StatusOK, pokemon)
}

func getPokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cacheKey := fmt.Sprintf("pokemon:%s", id)
	cached, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var pokemon Pokemon
		json.Unmarshal([]byte(cached), &pokemon)
		respondJSON(w, http.StatusOK, pokemon)
		return
	}

	var pokemon Pokemon
	if err := db.First(&pokemon, id).Error; err != nil {
		respondError(w, http.StatusNotFound, "Pokemon not found")
		return
	}

	data, _ := json.Marshal(pokemon)
	redisClient.Set(ctx, cacheKey, data, 10*time.Minute)

	respondJSON(w, http.StatusOK, pokemon)
}

func createPokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon Pokemon
	if err := json.NewDecoder(r.Body).Decode(&pokemon); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := db.Create(&pokemon).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Invalidate cache
	redisClient.Del(ctx, "pokemon:all")

	publishEvent("pokemon.created", pokemon)

	respondJSON(w, http.StatusCreated, pokemon)
}

func getRandomPokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon Pokemon
	if err := db.Order("RANDOM()").First(&pokemon).Error; err != nil {
		respondError(w, http.StatusNotFound, "No Pokemon found")
		return
	}

	respondJSON(w, http.StatusOK, pokemon)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// Helpers
func publishEvent(routingKey string, data interface{}) {
	body, _ := json.Marshal(data)
	err := rabbitCh.Publish("pokemon.events", routingKey, false, false, amqp.Publishing{
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