package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func Redis(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{client: rdb}
}

func (r *RedisRepository) EnqueueChat(ctx context.Context, roomID string) error {
	return r.client.LPush(ctx, "queue:chat", roomID).Err()
}

func (r *RedisRepository) PopChat(ctx context.Context) (string, error) {
	return r.client.RPop(ctx, "queue:chat").Result()
}

func (r *RedisRepository) GetAssignedAgent(ctx context.Context, roomID string) (string, error) {
	return r.client.Get(ctx, "room:"+roomID).Result()
}

func (r *RedisRepository) AssignAgent(ctx context.Context, roomID, agentEmail string) error {
	return r.client.Set(ctx, "room:"+roomID, agentEmail, 0).Err()
}

func (r *RedisRepository) SetAgentChatCount(ctx context.Context, agentEmail string, count int) error {
	return r.client.Set(ctx, "agent:"+agentEmail+":count", count, 0).Err()
}

func (r *RedisRepository) GetAgentChatCount(ctx context.Context, agentEmail string) (int, error) {
	return r.client.Get(ctx, "agent:"+agentEmail+":count").Int()
}

func (r *RedisRepository) IncrementAgentChatCount(ctx context.Context, agentEmail string) error {
	return r.client.Incr(ctx, "agent:"+agentEmail+":count").Err()
}
