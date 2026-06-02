package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"main/pkg/eventbus"
	"main/pkg/eventbus/contracts"
	"sync"

	"github.com/IBM/sarama"
)

type KafkaEventBus struct {
	producer sarama.SyncProducer
	handlers map[eventbus.EventType][]eventbus.EventHandler
	mu       sync.RWMutex
}

func NewKafkaEventBus(brokers []string) (*KafkaEventBus, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaEventBus{
		producer: producer,
		handlers: make(map[eventbus.EventType][]eventbus.EventHandler),
	}, nil
}

func (b *KafkaEventBus) Subscribe(eventType eventbus.EventType, handler eventbus.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *KafkaEventBus) Publish(ctx context.Context, event eventbus.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: string(event.Type()),
		Value: sarama.ByteEncoder(payload),
	}

	_, _, err = b.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *KafkaEventBus) StartConsumerGroup(ctx context.Context, brokers []string, groupID string) error {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}
	defer func(client sarama.ConsumerGroup) {
		err := client.Close()
		if err != nil {
			log.Printf("Failed to close consumer group client: %v", err)
		} else {
			log.Println("Consumer group client closed")
		}
	}(client)

	b.mu.RLock()
	topics := make([]string, 0, len(b.handlers))
	for t := range b.handlers {
		topics = append(topics, string(t))
	}
	b.mu.RUnlock()

	if len(topics) == 0 {
		return errors.New("no topics defined")
	}

	handler := &consumerGroupHandler{bus: b}

	for {
		if err := client.Consume(ctx, topics, handler); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type consumerGroupHandler struct {
	bus *KafkaEventBus
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := h.handleMessage(msg)
		if err != nil {
			return err
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func (h *consumerGroupHandler) handleMessage(msg *sarama.ConsumerMessage) error {
	eventType := eventbus.EventType(msg.Topic)

	h.bus.mu.RLock()
	handlers, exists := h.bus.handlers[eventType]
	h.bus.mu.RUnlock()

	if !exists {
		return nil
	}

	for _, handler := range handlers {
		ev, err := h.unmarshalEvent(eventType, msg.Value)
		if err != nil {
			return err
		}
		err = handler.Handle(ev)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *consumerGroupHandler) unmarshalEvent(eventType eventbus.EventType, e []byte) (eventbus.Event, error) {
	switch eventType {
	case contracts.UserCreatedEventType:
		var event contracts.UserCreated
		err := json.Unmarshal(e, &event)
		if err != nil {
			log.Println("Error unmarshaling event: ", err)
		}
		return event, nil
	default:
		return nil, errors.New("unknown event type")
	}
}
