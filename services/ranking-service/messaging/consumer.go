package messaging

import (
	"encoding/json"
	"log"

	"maushold/ranking-service/service"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel        *amqp.Channel
	rankingService service.RankingService
}

func NewConsumer(channel *amqp.Channel, rankingService service.RankingService) *Consumer {
	return &Consumer{
		channel:        channel,
		rankingService: rankingService,
	}
}

func (c *Consumer) Start() {
	msgs, err := c.channel.Consume(
		"ranking.updates",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	log.Println("Listening for battle events...")

	for msg := range msgs {
		log.Printf("Received message: %s", msg.RoutingKey)

		switch msg.RoutingKey {
		case "battle.completed":
			c.handleBattleCompleted(msg.Body)
		case "player.deleted":
			c.handlePlayerDeleted(msg.Body)
		}
	}
}

func (c *Consumer) handlePlayerDeleted(body []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error parsing player deleted event: %v", err)
		return
	}

	playerIDRaw, ok := event["player_id"]
	if !ok {
		log.Printf("Error: player_id missing in player.deleted event")
		return
	}

	playerID := uint(playerIDRaw.(float64))
	log.Printf("Processing player deletion: PlayerID=%d", playerID)

	err := c.rankingService.DeletePlayerRanking(playerID)
	if err != nil {
		log.Printf("Error processing player deletion: %v", err)
	} else {
		log.Printf("Player %d ranking deleted successully", playerID)
	}
}

func (c *Consumer) handleBattleCompleted(body []byte) {
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

	c.rankingService.UpdatePlayerRanking(winnerID, pointsWon, true)
	c.rankingService.UpdatePlayerRanking(loserID, -pointsLost, false)

	log.Printf("Battle completed event processed")
}
