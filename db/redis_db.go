package db

import (
	"context"
	config "eCommerce-go/utils"
	logger "eCommerce-go/utils"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Client *redis.Client
}

var RedisStoreInstance *RedisStore

func NewRedisStore() (*RedisStore, error) {
	log := logger.GetLogger()
	log.With("component", "redis_store")
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Env.RedisAddr,
		Password: config.Env.RedisPassword,
		DB:       config.Env.RedisDB,
	})

	// Ping Redis to check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Successfully connected to Redis", "address", config.Env.RedisAddr)
	RedisStoreInstance = &RedisStore{Client: rdb}
	return RedisStoreInstance, nil
}

// StoreRefreshToken stores the refresh token details (UUID -> UserID) with expiry
func (rs *RedisStore) StoreRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	log := logger.GetLogger()
	log.With("component", "store_refresh_token")
	err := rs.Client.Set(ctx, tokenID, userID, expiresIn).Err()
	if err != nil {
		log.Error("Failed to store refresh token in Redis", "tokenID", tokenID, "userID", userID, "error", err)
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	log.Debug("Stored refresh token in Redis", "tokenID", tokenID, "userID", userID, "expiresIn", expiresIn)
	return nil
}

// ValidateRefreshToken checks if the token ID exists in Redis and returns the associated UserID
func (rs *RedisStore) ValidateRefreshToken(ctx context.Context, tokenID string) (string, error) {
	log := logger.GetLogger()
	log.With("component", "validate_refresh_token")
	userID, err := rs.Client.Get(ctx, tokenID).Result()
	if err == redis.Nil {
		log.Warn("Refresh token not found in Redis or expired", "tokenID", tokenID)
		return "", fmt.Errorf("refresh token not found or expired")
	} else if err != nil {
		log.Error("Failed to validate refresh token in Redis", "tokenID", tokenID, "error", err)
		return "", fmt.Errorf("failed to validate refresh token: %w", err)
	}
	log.Debug("Validated refresh token in Redis", "tokenID", tokenID, "userID", userID)
	return userID, nil
}

// DeleteRefreshToken deletes a refresh token by its ID
func (rs *RedisStore) DeleteRefreshToken(ctx context.Context, tokenID string) error {
	log := logger.GetLogger()
	log.With("component", "delete_refresh_token")
	deletedCount, err := rs.Client.Del(ctx, tokenID).Result()
	if err != nil {
		log.Error("Failed to delete refresh token from Redis", "tokenID", tokenID, "error", err)
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	if deletedCount == 0 {
		log.Warn("Attempted to delete non-existent refresh token", "tokenID", tokenID)
		// Depending on the use case (logout), this might not be an error
		// return fmt.Errorf("refresh token not found for deletion")
	} else {
		log.Info("Deleted refresh token from Redis", "tokenID", tokenID)
	}
	return nil
}

// Close closes the Redis connection
func (rs *RedisStore) Close() error {
	return rs.Client.Close()
}
