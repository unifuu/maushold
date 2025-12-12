package messaging

import (
	"encoding/json"
	"log"

	"maushold/player-service/service"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel       *amqp.Channel
	playerService service.PlayerService
}

func NewConsumer(channel *amqp.Channel, playerService service.PlayerService) *Consumer {
	return &Consumer{
		channel:       channel,
		playerService: playerService,
	}
}

func (c *Consumer) Start() {
	msgs, err := c.channel.Consume(
		"player.updates",
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

	log.Println("Listening for messages...")

	for msg := range msgs {
		log.Printf("Received message: %s", msg.RoutingKey)

		switch msg.RoutingKey {
		case "battle.completed":
			c.handleBattleCompleted(msg.Body)
		}
	}
}

func (c *Consumer) handleBattleCompleted(body []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error parsing battle event: %v", err)
		return
	}

	// Update player points based on battle result
	if winnerID, ok := event["winner_id"].(float64); ok {
		if pointsWon, ok := event["points_won"].(float64); ok {
			c.playerService.UpdatePlayerPoints(uint(winnerID), int(pointsWon))
		}
	}

	if loserID, ok := event["loser_id"].(float64); ok {
		if pointsLost, ok := event["points_lost"].(float64); ok {
			c.playerService.UpdatePlayerPoints(uint(loserID), -int(pointsLost))
		}
	}

	log.Printf("Battle completed event processed: %v", event)
}
