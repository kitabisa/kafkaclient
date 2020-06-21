package kafkaclient

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer promotes the origin consumer of kafka client.
type Consumer struct {
	*kafka.Consumer
}

// EventHandler is a handler for handling event.
type EventHandler func(cons *Consumer, event kafka.Event)

// MessageHandler is a handler for handling message.
type MessageHandler func(cons *Consumer, msg *kafka.Message, err error)

// GetOrigin returns the origin consumer of kafka.
func (c *Consumer) GetOrigin() *kafka.Consumer {
	return c.Consumer
}

func (c *Consumer) consume(args ConsumeArgs) (err error) {
	err = c.SubscribeTopics(args.Topics, args.RebalanceCb)
	if err != nil {
		return
	}

	go func(c *Consumer, args ConsumeArgs) {
		var (
			err error
			msg *kafka.Message
		)
		for {
			msg, err = c.ReadMessage(time.Duration(args.Polling) * time.Millisecond)
			args.Handler(c, msg, err)
		}
	}(c, args)

	return
}

func (c *Consumer) consumeEvent(args ConsumeArgs) (err error) {
	err = c.SubscribeTopics(args.Topics, args.RebalanceCb)
	if err != nil {
		return
	}

	go func(c *Consumer, args ConsumeArgs) {
		var (
			event kafka.Event
		)
		for {
			event = c.Poll(args.Polling)
			args.EventHandler(c, event)
		}
	}(c, args)
	return
}
