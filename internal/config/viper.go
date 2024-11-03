package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Env struct {
	Environment    string   `mapstructure:"ENVIRONMENT"`
	AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`
	DBSource       string   `mapstructure:"DB_SOURCE"`
	DBSourceTest   string   `mapstructure:"DB_SOURCE_TEST"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadEnv(path string) (env Env, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&env)
	return
}
