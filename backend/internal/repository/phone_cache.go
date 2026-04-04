package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/redis/go-redis/v9"
)

type PhoneCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewPhoneCache(client *redis.Client) *PhoneCache {
	if client == nil {
		return &PhoneCache{
			client: nil,
			ttl:    10 * time.Minute,
		}
	}
	return &PhoneCache{
		client: client,
		ttl:    10 * time.Minute,
	}
}

func (c *PhoneCache) isEnabled() bool {
	return c.client != nil
}

// SetList stores the entire phone list in Redis as a JSON string
func (c *PhoneCache) SetList(ctx context.Context, phones []model.Phone) error {
	if !c.isEnabled() {
		return nil
	}
	data, err := json.Marshal(phones)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "phones:all", data, c.ttl).Err()
}

// GetList retrieves the phone list from Redis
func (c *PhoneCache) GetList(ctx context.Context) ([]model.Phone, error) {
	if !c.isEnabled() {
		return nil, nil
	}
	data, err := c.client.Get(ctx, "phones:all").Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	var phones []model.Phone
	err = json.Unmarshal([]byte(data), &phones)
	return phones, err
}

// SetByID stores a single phone using its ID as part of the key
func (c *PhoneCache) SetByID(ctx context.Context, id int, phone model.Phone) error {
	if !c.isEnabled() {
		return nil
	}
	data, err := json.Marshal(phone)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("phone:%d", id)
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// GetByID tries to fetch a single phone from Redis
func (c *PhoneCache) GetByID(ctx context.Context, id int) (*model.Phone, error) {
	if !c.isEnabled() {
		return nil, nil
	}
	key := fmt.Sprintf("phone:%d", id)
	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	var phone model.Phone
	err = json.Unmarshal([]byte(data), &phone)
	return &phone, err
}

// ClearByID removes a specific phone from the cache
func (c *PhoneCache) ClearByID(ctx context.Context, id int) {
	if !c.isEnabled() {
		return
	}
	key := fmt.Sprintf("phone:%d", id)
	c.client.Del(ctx, key)
}

// Clear removes the cached list (useful after an admin updates/deletes a phone)
func (c *PhoneCache) Clear(ctx context.Context) {
	if !c.isEnabled() {
		return
	}
	c.client.Del(ctx, "phones:all")
}
