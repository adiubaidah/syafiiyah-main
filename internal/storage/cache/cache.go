package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redisClient *redis.Client
}

func NewClient(redisClient *redis.Client) *Cache {
	return &Cache{redisClient: redisClient}
}

func (c *Cache) CreateSession(session Session) error {
	ctx := context.Background()
	sessionMap := map[string]any{
		"username":      session.Username,
		"refresh_token": session.RefreshToken,
		"user_agent":    session.UserAgent,
		"client_ip":     session.ClientIp,
		"is_blocked":    session.IsBlocked,
		"expires_at":    session.ExpiresAt.Format(time.RFC3339),
		"created_at":    session.CreatedAt.Format(time.RFC3339),
	}

	err := c.redisClient.HMSet(ctx, "session:"+session.ID.String(), sessionMap).Err()
	if err != nil {
		return err
	}

	c.redisClient.Expire(ctx, "session:"+session.ID.String(), time.Until(session.ExpiresAt))
	return nil
}

func (c *Cache) GetSession(id string) (Session, error) {
	ctx := context.Background()
	sessionMap, err := c.redisClient.HGetAll(ctx, "session:"+id).Result()
	if err != nil {

		if err == redis.Nil {
			return Session{}, exception.NewNotFoundError("session not found")
		}

		return Session{}, err
	}

	expiresAt, err := time.Parse(time.RFC3339, sessionMap["expires_at"])
	if err != nil {
		return Session{}, exception.NewValidationError("invalid time format")
	}

	createdAt, err := time.Parse(time.RFC3339, sessionMap["created_at"])
	if err != nil {
		return Session{}, exception.NewValidationError("invalid time format")
	}

	isBlocked, err := strconv.ParseBool(sessionMap["is_blocked"])
	if err != nil {
		return Session{}, err
	}
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return Session{}, err
	}

	return Session{
		ID:           uuidID,
		Username:     sessionMap["username"],
		RefreshToken: sessionMap["refresh_token"],
		UserAgent:    sessionMap["user_agent"],
		ClientIp:     sessionMap["client_ip"],
		IsBlocked:    isBlocked,
		ExpiresAt:    expiresAt,
		CreatedAt:    createdAt,
	}, nil
}

func (c *Cache) DeleteSession(id string) error {
	ctx := context.Background()
	return c.redisClient.Del(ctx, "session:"+id).Err()
}
