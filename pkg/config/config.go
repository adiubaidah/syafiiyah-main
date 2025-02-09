package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
	Environment            string        `mapstructure:"ENVIRONMENT"`
	ServerAddress          string        `mapstructure:"SERVER_ADDRESS"`
	ServerPublicUrl        string        `mapstructure:"SERVER_PUBLIC_URL"`
	DBDriver               string        `mapstructure:"DB_DRIVER"`
	DBSource               string        `mapstructure:"DB_SOURCE"`
	MigrationURL           string        `mapstructure:"MIGRATION_URL"`
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration   time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	GoogleOauthClient      string        `mapstructure:"GOOGLE_OAUTH_CLIENT"`
	MQTTBroker             string        `mapstructure:"MQTT_BROKER"`
	RedisAddress           string        `mapstructure:"REDIS_ADDRESS"`
	DBRedis                int           `mapstructure:"DB_REDIS"`
	AWSRegion              string        `mapstructure:"AWS_REGION"`
	AWSAccessKey           string        `mapstructure:"AWS_ACCESS_KEY"`
	AWSSecretKey           string        `mapstructure:"AWS_SECRET_KEY"`
	AWSBucketName          string        `mapstructure:"AWS_BUCKET_NAME"`
	ScheduleServiceAddress string        `mapstructure:"SCHEDULE_SERVICE_ADDRESS"`
}

const PathPhoto = "internal/storage/photo"

func Load(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
