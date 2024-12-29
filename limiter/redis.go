package strategy

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStrategy struct {
	client *redis.Client
}

func NewRedisStrategy(host, port string) *RedisStrategy {
	return &RedisStrategy{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", host, port),
		}),
	}
}

// Incrementa o contador para a chave especificada
func (r *RedisStrategy) IncrementCounter(key string) (int, error) {
	ctx := context.Background()
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Obtém o valor atual do contador para a chave especificada
func (r *RedisStrategy) GetCounter(key string) (int, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil // chave não existe
	} else if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Define a expiração para a chave especificada
func (r *RedisStrategy) SetExpiration(key string, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Expire(ctx, key, expiration).Err()
}

// Verifica se a chave está bloqueada
func (r *RedisStrategy) IsBlocked(key string) (bool, error) {
	ctx := context.Background()
	blocked, err := r.client.Get(ctx, "blocked:"+key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Se existe a chave "blocked:<key>", retorna true
	return blocked == "true", nil
}

// Define a chave como bloqueada por um período de tempo
func (r *RedisStrategy) SetBlocked(key string, duration time.Duration) error {
	ctx := context.Background()
	// Define "blocked:<key>" = "true" com TTL de `duration`
	err := r.client.Set(ctx, "blocked:"+key, "true", duration).Err()
	if err != nil {
		return err
	}
	return nil
}
