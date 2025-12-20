package config

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func InitRabbitMQ(cfg *Config) (*amqp.Connection, *amqp.Channel) {
	var conn *amqp.Connection
	var err error

	// Retry connection
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(cfg.RabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/10): %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ after 10 attempts:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}

	// Declare OUR exchange (player.events)
	err = ch.ExchangeDeclare("player.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare player.events exchange:", err)
	}

	// Declare battle.events exchange (so we can bind to it)
	// This is idempotent - if battle-service already created it, this is fine
	err = ch.ExchangeDeclare("battle.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("Warning: Failed to declare battle.events exchange: %v", err)
	}

	// Declare our queue
	_, err = ch.QueueDeclare("ranking.updates", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	// Bind queue to battle events (safe now that exchange exists)
	err = ch.QueueBind("ranking.updates", "battle.completed", "battle.events", false, nil)
	if err != nil {
		log.Printf("Warning: Failed to bind queue (will retry later): %v", err)
		// Don't fail here - battle service might not be up yet
	}

	log.Println("RabbitMQ connected")
	return conn, ch
}
