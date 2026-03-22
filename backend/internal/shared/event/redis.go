package event

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisEventBus struct {
	client *redis.Client
	stream string
	group  string
}

func NewRedisEventBus(client *redis.Client, group string) *RedisEventBus {
	return &RedisEventBus{
		client: client,
		stream: "domain_events",
		group:  group,
	}
}

func (b *RedisEventBus) Publish(ctx context.Context, evt Event) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	
	err = b.client.XAdd(ctx, &redis.XAddArgs{
		Stream: b.stream,
		Values: map[string]interface{}{"payload": data},
	}).Err()
	
	if err == nil {
		log.Printf("📤 [Redis] Event published: %s (org_id=%d)", evt.Type, evt.OrganizationID)
	}
	return err
}

func (b *RedisEventBus) Consume(ctx context.Context, handler func(Event) error) {
	// Create consumer group if not exists (ignore error if already exists)
	_ = b.client.XGroupCreateMkStream(ctx, b.stream, b.group, "0").Err()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			streams, err := b.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    b.group,
				Consumer: "backend-main",
				Streams:  []string{b.stream, ">"},
				Count:    10,
				Block:    time.Second * 5,
			}).Result()

			if err != nil {
				if err != redis.Nil {
					log.Printf("❌ [Redis] Error reading events: %v", err)
					time.Sleep(time.Second)
				}
				continue
			}

			for _, s := range streams {
				for _, msg := range s.Messages {
					var evt Event
					payload, ok := msg.Values["payload"].(string)
					if !ok {
						continue
					}
					
					if err := json.Unmarshal([]byte(payload), &evt); err != nil {
						log.Printf("❌ [Redis] Failed to unmarshal event: %v", err)
						continue
					}

					if err := handler(evt); err == nil {
						b.client.XAck(ctx, b.stream, b.group, msg.ID)
					} else {
						log.Printf("❌ [Redis] Handler error for %s: %v", evt.Type, err)
					}
				}
			}
		}
	}
}

func (b *RedisEventBus) Close() error {
	return nil
}
