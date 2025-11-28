package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	RedisClient *redis.Client
)

// InitRedis initializes the Redis client
func InitRedis(addr, password string, db int) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// Permission cache keys and TTL
const (
	PermissionCacheTTL  = 5 * time.Minute
	PermissionKeyPrefix = "perm:user:"
)

// GetPermissionCacheKey generates a cache key for user permissions
func GetPermissionCacheKey(userID uint, objType string, objID uint) string {
	return fmt.Sprintf("%s%d:obj:%s:%d", PermissionKeyPrefix, userID, objType, objID)
}

// GetPermissionCache retrieves cached permission mask for a user on an object
func GetPermissionCache(userID uint, objType string, objID uint) (int, bool) {
	if RedisClient == nil {
		return 0, false
	}

	key := GetPermissionCacheKey(userID, objType, objID)
	val, err := RedisClient.Get(ctx, key).Int()
	if err != nil {
		return 0, false
	}

	return val, true
}

// SetPermissionCache caches the permission mask for a user on an object
func SetPermissionCache(userID uint, objType string, objID uint, permMask int) error {
	if RedisClient == nil {
		return nil // Silently skip if Redis not available
	}

	key := GetPermissionCacheKey(userID, objType, objID)
	return RedisClient.Set(ctx, key, permMask, PermissionCacheTTL).Err()
}

// InvalidateUserPermissions clears all permission cache for a specific user
func InvalidateUserPermissions(userID uint) error {
	if RedisClient == nil {
		return nil
	}

	// Use SCAN to find all keys matching the pattern
	pattern := fmt.Sprintf("%s%d:obj:*", PermissionKeyPrefix, userID)
	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	// Delete all found keys
	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

// InvalidateObjectPermissions clears all permission cache for a specific object
func InvalidateObjectPermissions(objType string, objID uint) error {
	if RedisClient == nil {
		return nil
	}

	// Use SCAN to find all keys matching the pattern
	pattern := fmt.Sprintf("%s*:obj:%s:%d", PermissionKeyPrefix, objType, objID)
	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	// Delete all found keys
	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

// InvalidateAllPermissions clears all permission cache
func InvalidateAllPermissions() error {
	if RedisClient == nil {
		return nil
	}

	// Use SCAN to find all permission keys
	pattern := fmt.Sprintf("%s*", PermissionKeyPrefix)
	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	// Delete all found keys
	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}
