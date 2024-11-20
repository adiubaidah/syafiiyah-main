package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedisSetGet(t *testing.T) {
	ctx := context.Background()
	key := "test_key"
	value := "test_value"

	// Set the key-value pair
	err := redisClient.Set(ctx, key, value, 0).Err()
	require.NoError(t, err)

	// Get the value
	val, err := redisClient.Get(ctx, key).Result()
	require.NoError(t, err)

	require.Equal(t, value, val)

	// Clean up
	err = redisClient.Del(ctx, key).Err()
	require.NoError(t, err)
}
