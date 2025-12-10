package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Models
type PlayerRanking struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PlayerID      uint      `gorm:"unique;not null;index" json:"player_id"`
	Username      string    `json:"username"`
	TotalPoints   int       `gorm:"default:1000;index" json:"total_points"`
	TotalBattles  int       `gorm:"default:0" json:"total_battles"`
	Wins          int       `gorm:"default:0" json:"wins"`
	Losses        int       `gorm:"default:0" json:"losses"`
	WinRate       float64   `json:"win_rate"`
	Rank          int       `json:"rank"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type LeaderboardEntry struct {
	PlayerID    uint    `json:"player_id"`
	Username    string  `json:"username"`
	TotalPoints int     `json:"total_points"`
	Wins        int     `json:"wins"`
	Losses      int     `json:"losses"`
	WinRate     float64 `json:"win_rate"`
	Rank        int     `json:"rank"`
}

type Player struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Points   int    `json:"points"`
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

	// Start consuming battle events
	go consumeBattleEvents()

	// Sync rankings every 5 minutes
	go syncRankingsPeriodically()

	r := mux.NewRouter()

	r.HandleFunc("/rankings", getLeaderboard).Methods("GET")
	r.HandleFunc("/rankings/player/{playerId}", getPlayerRanking).Methods("GET")
	r.HandleFunc("/rankings/sync", syncRankings).Methods("POST")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	r.Use(corsMiddleware)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8004"
	}

	log.Printf("Ranking Service starting on port %s", port)
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

	db.AutoMigrate(&PlayerRanking{})
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

	// Declare queue for battle events
	_, err = rabbitCh.QueueDeclare("ranking.updates", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	// Bind to battle events
	err = rabbitCh.QueueBind("ranking.updates", "battle.completed", "battle.events", false, nil)
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

// Event Consumer
func consumeBattleEvents() {
	msgs, err := rabbitCh.Consume("ranking.updates", "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	log.Println("Listening for battle events...")

	for msg := range msgs {
		log.Printf("Received event: %s", msg.RoutingKey)
		
		if msg.RoutingKey == "battle.completed" {
			handleBattleCompleted(msg.Body)
		}
		
		msg.Ack(false)
	}
}

func handleBattleCompleted(body []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error parsing battle event: %v", err)
		return
	}

	winnerID := uint(event["winner_id"].(float64))
	loserID := uint(event["loser_id"].(float64))
	pointsWon := int(event["points_won"].(float64))
	pointsLost := int(event["points_lost"].(float64))

	log.Printf("Processing battle: Winner=%d (+%d), Loser=%d (-%d)", 
		winnerID, pointsWon, loserID, pointsLost)

	// Update winner
	updatePlayerRanking(winnerID, pointsWon, true)

	// Update loser
	updatePlayerRanking(loserID, -pointsLost, false)

	// Update Redis leaderboard
	updateRedisLeaderboard(winnerID, pointsWon)
	updateRedisLeaderboard(loserID, -pointsLost)
}

func updatePlayerRanking(playerID uint, pointsDelta int, isWin bool) {
	var ranking PlayerRanking
	result := db.Where("player_id = ?", playerID).First(&ranking)

	if result.Error == gorm.ErrRecordNotFound {
		// Fetch player info
		player := fetchPlayerInfo(playerID)
		if player == nil {
			log.Printf("Player %d not found", playerID)
			return
		}

		ranking = PlayerRanking{
			PlayerID:    playerID,
			Username:    player.Username,
			TotalPoints: player.Points + pointsDelta,
			TotalBattles: 1,
		}

		if isWin {
			ranking.Wins = 1
		} else {
			ranking.Losses = 1
		}
	} else {
		ranking.TotalPoints += pointsDelta
		ranking.TotalBattles++
		
		if isWin {
			ranking.Wins++
		} else {
			ranking.Losses++
		}
	}

	// Calculate win rate
	if ranking.TotalBattles > 0 {
		ranking.WinRate = float64(ranking.Wins) / float64(ranking.TotalBattles) * 100
	}

	db.Save(&ranking)
	log.Printf("Updated ranking for player %d: Points=%d, W/L=%d/%d", 
		playerID, ranking.TotalPoints, ranking.Wins, ranking.Losses)
}

func updateRedisLeaderboard(playerID uint, pointsDelta int) {
	// Get current score
	currentScore, err := redisClient.ZScore(ctx, "leaderboard", fmt.Sprintf("%d", playerID)).Result()
	if err != nil {
		// Player not in leaderboard, fetch from DB
		var ranking PlayerRanking
		if db.Where("player_id = ?", playerID).First(&ranking).Error == nil {
			currentScore = float64(ranking.TotalPoints)
		} else {
			currentScore = 1000 // Default starting points
		}
	}

	newScore := currentScore + float64(pointsDelta)
	redisClient.ZAdd(ctx, "leaderboard", &redis.Z{
		Score:  newScore,
		Member: fmt.Sprintf("%d", playerID),
	})

	log.Printf("Updated Redis leaderboard for player %d: %f -> %f", 
		playerID, currentScore, newScore)
}

// HTTP Handlers
func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	// Try Redis first for real-time leaderboard
	result, err := redisClient.ZRevRangeWithScores(ctx, "leaderboard", 0, int64(limit-1)).Result()
	if err == nil && len(result) > 0 {
		leaderboard := make([]LeaderboardEntry, 0, len(result))
		
		for i, z := range result {
			playerID, _ := strconv.ParseUint(z.Member.(string), 10, 64)
			
			var ranking PlayerRanking
			db.Where("player_id = ?", playerID).First(&ranking)
			
			leaderboard = append(leaderboard, LeaderboardEntry{
				PlayerID:    uint(playerID),
				Username:    ranking.Username,
				TotalPoints: int(z.Score),
				Wins:        ranking.Wins,
				Losses:      ranking.Losses,
				WinRate:     ranking.WinRate,
				Rank:        i + 1,
			})
		}
		
		respondJSON(w, http.StatusOK, leaderboard)
		return
	}

	// Fallback to database
	var rankings []PlayerRanking
	db.Order("total_points DESC").Limit(limit).Find(&rankings)

	leaderboard := make([]LeaderboardEntry, len(rankings))
	for i, r := range rankings {
		leaderboard[i] = LeaderboardEntry{
			PlayerID:    r.PlayerID,
			Username:    r.Username,
			TotalPoints: r.TotalPoints,
			Wins:        r.Wins,
			Losses:      r.Losses,
			WinRate:     r.WinRate,
			Rank:        i + 1,
		}
	}

	respondJSON(w, http.StatusOK, leaderboard)
}

func getPlayerRanking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID := vars["playerId"]

	var ranking PlayerRanking
	if err := db.Where("player_id = ?", playerID).First(&ranking).Error; err != nil {
		respondError(w, http.StatusNotFound, "Player ranking not found")
		return
	}

	// Calculate rank
	var rank int64
	db.Model(&PlayerRanking{}).Where("total_points > ?", ranking.TotalPoints).Count(&rank)
	ranking.Rank = int(rank) + 1

	respondJSON(w, http.StatusOK, ranking)
}

func syncRankings(w http.ResponseWriter, r *http.Request) {
	go performRankingSync()
	respondJSON(w, http.StatusOK, map[string]string{"message": "Sync started"})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// Background Tasks
func syncRankingsPeriodically() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		performRankingSync()
	}
}

func performRankingSync() {
	log.Println("Starting ranking sync...")

	var rankings []PlayerRanking
	db.Order("total_points DESC").Find(&rankings)

	// Clear and rebuild Redis leaderboard
	redisClient.Del(ctx, "leaderboard")

	for _, r := range rankings {
		redisClient.ZAdd(ctx, "leaderboard", &redis.Z{
			Score:  float64(r.TotalPoints),
			Member: fmt.Sprintf("%d", r.PlayerID),
		})
	}

	log.Printf("Ranking sync completed. Synced %d players", len(rankings))
}

// Helper Functions
func fetchPlayerInfo(playerID uint) *Player {
	url := fmt.Sprintf("%s/players/%d", os.Getenv("PLAYER_SERVICE_URL"), playerID)
	
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching player %d: %v", playerID, err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	var player Player
	if err := json.Unmarshal(body, &player); err != nil {
		return nil
	}

	return &player
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