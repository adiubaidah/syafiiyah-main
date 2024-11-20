package cache

import (
	"log"
	"os"
	"testing"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	redis "github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func TestMain(m *testing.M) {
	env, err := config.LoadConfig("../../..")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: env.RedisAddress,
		DB:   env.DBRedis,
	})

	defer redisClient.Close()

	exitCode := m.Run()
	os.Exit(exitCode)
}
