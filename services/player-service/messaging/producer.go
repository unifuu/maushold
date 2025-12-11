package messaging

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	channel *amqp.Channel
}

func NewProducer(channel *amqp.Channel) *Producer {
	return &Producer{channel: channel}
}

func (p *Producer) PublishPlayerEvent(routingKey string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"player.events",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish event %s: %v", routingKey, err)
		return err
	}

	log.Printf("Published event: %s", routingKey)
	return nil
}
