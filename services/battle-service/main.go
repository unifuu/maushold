package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
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
type Battle struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Player1ID   uint      `gorm:"not null;index" json:"player1_id"`
	Player2ID   uint      `gorm:"not null;index" json:"player2_id"`
	Pokemon1ID  uint      `gorm:"not null" json:"pokemon1_id"`
	Pokemon2ID  uint      `gorm:"not null" json:"pokemon2_id"`
	WinnerID    uint      `json:"winner_id"`
	Status      string    `gorm:"default:'pending'" json:"status"` // pending, in_progress, completed
	BattleLog   string    `gorm:"type:text" json:"battle_log"`
	PointsWon   int       `json:"points_won"`
	PointsLost  int       `json:"points_lost"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type BattleRequest struct {
	Player1ID  uint `json:"player1_id"`
	Player2ID  uint `json:"player2_id"`
	Pokemon1ID uint `json:"pokemon1_id"`
	Pokemon2ID uint `json:"pokemon2_id"`
}

type PlayerPokemon struct {
	ID       uint   `json:"id"`
	PlayerID uint   `json:"player_id"`
	HP       int    `json:"hp"`
	Attack   int    `json:"attack"`
	Defense  int    `json:"defense"`
	Speed    int    `json:"speed"`
	Level    int    `json:"level"`
	Nickname string `json:"nickname"`
}

type Player struct {
	ID     uint   `json:"id"`
	Points int    `json:"points"`
	Username string `json:"username"`
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
	rand.Seed(time.Now().UnixNano())
	
	initDB()
	initRedis()
	initRabbitMQ()
	defer cleanup()

	r := mux.NewRouter()

	r.HandleFunc("/battles", createBattle).Methods("POST")
	r.HandleFunc("/battles/{id}", getBattle).Methods("GET")
	r.HandleFunc("/battles/player/{playerId}", getPlayerBattles).Methods("GET")
	r.HandleFunc("/battles", getAllBattles).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	r.Use(corsMiddleware)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8003"
	}

	log.Printf("Battle Service starting on port %s", port)
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

	db.AutoMigrate(&Battle{})
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

	err = rabbitCh.ExchangeDeclare("battle.events", "topic", true, false, false, false, nil)
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

// HTTP Handlers
func createBattle(w http.ResponseWriter, r *http.Request) {
	var req BattleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Fetch Pokemon data
	pokemon1, err := fetchPlayerPokemon(req.Player1ID, req.Pokemon1ID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Pokemon 1 not found")
		return
	}

	pokemon2, err := fetchPlayerPokemon(req.Player2ID, req.Pokemon2ID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Pokemon 2 not found")
		return
	}

	// Create battle record
	battle := Battle{
		Player1ID:  req.Player1ID,
		Player2ID:  req.Player2ID,
		Pokemon1ID: req.Pokemon1ID,
		Pokemon2ID: req.Pokemon2ID,
		Status:     "in_progress",
	}

	if err := db.Create(&battle).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Simulate battle
	winner, battleLog := simulateBattle(pokemon1, pokemon2)
	
	// Determine winner
	if winner == 1 {
		battle.WinnerID = req.Player1ID
	} else {
		battle.WinnerID = req.Player2ID
	}

	// Calculate points
	battle.PointsWon = 50 + rand.Intn(50)
	battle.PointsLost = 20 + rand.Intn(30)
	battle.BattleLog = battleLog
	battle.Status = "completed"
	now := time.Now()
	battle.CompletedAt = &now

	db.Save(&battle)

	// Publish battle completed event
	publishBattleEvent(battle)

	respondJSON(w, http.StatusCreated, battle)
}

func getBattle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var battle Battle
	if err := db.First(&battle, id).Error; err != nil {
		respondError(w, http.StatusNotFound, "Battle not found")
		return
	}

	respondJSON(w, http.StatusOK, battle)
}

func getPlayerBattles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerId := vars["playerId"]

	var battles []Battle
	err := db.Where("player1_id = ? OR player2_id = ?", playerId, playerId).
		Order("created_at DESC").
		Limit(20).
		Find(&battles).Error

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, battles)
}

func getAllBattles(w http.ResponseWriter, r *http.Request) {
	var battles []Battle
	if err := db.Order("created_at DESC").Limit(50).Find(&battles).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, battles)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// Battle Logic
func simulateBattle(p1, p2 *PlayerPokemon) (int, string) {
	log := ""
	hp1 := p1.HP
	hp2 := p2.HP

	log += fmt.Sprintf("Battle Start!\n%s (HP: %d) vs %s (HP: %d)\n\n",
		p1.Nickname, hp1, p2.Nickname, hp2)

	round := 1
	for hp1 > 0 && hp2 > 0 && round <= 20 {
		log += fmt.Sprintf("=== Round %d ===\n", round)

		// Determine who attacks first based on speed
		if p1.Speed >= p2.Speed {
			damage := calculateDamage(p1.Attack, p2.Defense)
			hp2 -= damage
			log += fmt.Sprintf("%s attacks for %d damage! %s HP: %d\n",
				p1.Nickname, damage, p2.Nickname, maxInt(hp2, 0))

			if hp2 > 0 {
				damage = calculateDamage(p2.Attack, p1.Defense)
				hp1 -= damage
				log += fmt.Sprintf("%s attacks for %d damage! %s HP: %d\n",
					p2.Nickname, damage, p1.Nickname, maxInt(hp1, 0))
			}
		} else {
			damage := calculateDamage(p2.Attack, p1.Defense)
			hp1 -= damage
			log += fmt.Sprintf("%s attacks for %d damage! %s HP: %d\n",
				p2.Nickname, damage, p1.Nickname, maxInt(hp1, 0))

			if hp1 > 0 {
				damage = calculateDamage(p1.Attack, p2.Defense)
				hp2 -= damage
				log += fmt.Sprintf("%s attacks for %d damage! %s HP: %d\n",
					p1.Nickname, damage, p2.Nickname, maxInt(hp2, 0))
			}
		}

		log += "\n"
		round++
	}

	if hp1 > hp2 {
		log += fmt.Sprintf("🏆 %s wins!\n", p1.Nickname)
		return 1, log
	} else {
		log += fmt.Sprintf("🏆 %s wins!\n", p2.Nickname)
		return 2, log
	}
}

func calculateDamage(attack, defense int) int {
	baseDamage := attack - (defense / 2)
	if baseDamage < 1 {
		baseDamage = 1
	}
	
	// Add some randomness
	variance := rand.Intn(10) - 5
	damage := baseDamage + variance
	
	if damage < 1 {
		damage = 1
	}
	
	return damage
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper Functions
func fetchPlayerPokemon(playerID, pokemonID uint) (*PlayerPokemon, error) {
	url := fmt.Sprintf("%s/players/%d/pokemon", os.Getenv("PLAYER_SERVICE_URL"), playerID)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	var pokemons []PlayerPokemon
	if err := json.Unmarshal(body, &pokemons); err != nil {
		return nil, err
	}

	for _, p := range pokemons {
		if p.ID == pokemonID {
			if p.Nickname == "" {
				p.Nickname = fmt.Sprintf("Pokemon #%d", p.ID)
			}
			return &p, nil
		}
	}

	return nil, fmt.Errorf("pokemon not found")
}

func publishBattleEvent(battle Battle) {
	event := map[string]interface{}{
		"battle_id":   battle.ID,
		"winner_id":   battle.WinnerID,
		"loser_id":    getLoserID(battle),
		"points_won":  battle.PointsWon,
		"points_lost": battle.PointsLost,
		"timestamp":   time.Now().Unix(),
	}

	body, _ := json.Marshal(event)
	err := rabbitCh.Publish("battle.events", "battle.completed", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	
	if err != nil {
		log.Printf("Failed to publish battle event: %v", err)
	} else {
		log.Printf("Published battle.completed event for battle %d", battle.ID)
	}
}

func getLoserID(battle Battle) uint {
	if battle.WinnerID == battle.Player1ID {
		return battle.Player2ID
	}
	return battle.Player1ID
}

// HTTP Helpers
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