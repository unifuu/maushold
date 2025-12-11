package config

import (
	"log"

	"github.com/streadway/amqp"
)

func InitRabbitMQ(cfg *Config) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}

	// Declare exchanges
	err = ch.ExchangeDeclare("player.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	// Declare queue
	_, err = ch.QueueDeclare("player.updates", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	// Bind queue
	err = ch.QueueBind("player.updates", "battle.completed", "battle.events", false, nil)
	if err != nil {
		log.Fatal("Failed to bind queue:", err)
	}

	log.Println("RabbitMQ connected")
	return conn, ch
}
