package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/minhducta/agent-core-service/pkg/config"
	"github.com/redis/go-redis/v9"
)

// Cache wraps Redis client with application-specific caching logic
type Cache struct {
	client *redis.Client
	cfg    config.CacheConfig
}

// New creates a new cache instance
func New(redisCfg config.RedisConfig, cacheCfg config.CacheConfig) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisCfg.Address(),
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConns,
		MaxRetries:   redisCfg.MaxRetries,
		DialTimeout:  redisCfg.DialTimeout,
		ReadTimeout:  redisCfg.ReadTimeout,
		WriteTimeout: redisCfg.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{client: client, cfg: cacheCfg}, nil
}

// NewWithClient creates a cache instance with an existing Redis client — useful for tests
func NewWithClient(client *redis.Client, cfg config.CacheConfig) *Cache {
	return &Cache{client: client, cfg: cfg}
}

// Client returns the underlying Redis client
func (c *Cache) Client() *redis.Client {
	return c.client
}

// HealthCheck checks if Redis is reachable
func (c *Cache) HealthCheck(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Set stores a value in the cache with the given TTL
func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from the cache
func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key from the cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// --- Bot cache helpers ---

// SetBot caches a bot profile
func (c *Cache) SetBot(ctx context.Context, id string, value any) error {
	return c.Set(ctx, fmt.Sprintf("agent:bot:%s", id), value, c.cfg.BotTTL)
}

// GetBot retrieves a cached bot profile
func (c *Cache) GetBot(ctx context.Context, id string, dest any) error {
	return c.Get(ctx, fmt.Sprintf("agent:bot:%s", id), dest)
}

// DeleteBot removes a cached bot profile
func (c *Cache) DeleteBot(ctx context.Context, id string) error {
	return c.Delete(ctx, fmt.Sprintf("agent:bot:%s", id))
}

// SetBotByAPIKeyHash caches bot_id by api_key_hash for fast auth lookup
func (c *Cache) SetBotByAPIKeyHash(ctx context.Context, hash string, botID string) error {
	return c.client.Set(ctx, fmt.Sprintf("agent:apikey:%s", hash), botID, c.cfg.SessionTTL).Err()
}

// GetBotByAPIKeyHash retrieves bot_id from api_key_hash cache
func (c *Cache) GetBotByAPIKeyHash(ctx context.Context, hash string) (string, error) {
	return c.client.Get(ctx, fmt.Sprintf("agent:apikey:%s", hash)).Result()
}

// DeleteBotByAPIKeyHash removes a cached api_key → bot_id mapping
func (c *Cache) DeleteBotByAPIKeyHash(ctx context.Context, hash string) error {
	return c.Delete(ctx, fmt.Sprintf("agent:apikey:%s", hash))
}

// --- Memory cache helpers ---

// SetMemories caches a bot's memories list
func (c *Cache) SetMemories(ctx context.Context, botID string, value any) error {
	return c.Set(ctx, fmt.Sprintf("agent:memories:%s", botID), value, c.cfg.MemoryTTL)
}

// GetMemories retrieves cached memories for a bot
func (c *Cache) GetMemories(ctx context.Context, botID string, dest any) error {
	return c.Get(ctx, fmt.Sprintf("agent:memories:%s", botID), dest)
}

// InvalidateMemories removes cached memories for a bot
func (c *Cache) InvalidateMemories(ctx context.Context, botID string) error {
	return c.Delete(ctx, fmt.Sprintf("agent:memories:%s", botID))
}
