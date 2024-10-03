package config

import (
	"strings"

	"github.com/bosskrub9992/fuel-management-backend/library/databases"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Use      string
		Postgres databases.PostgresConfig
		SQLite   struct {
			FilePath string `mapstructure:"file_path"`
		}
	}
	Logger struct {
		IsProductionEnv bool `mapstructure:"is_production_env"`
	}
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")     // local
	viper.AddConfigPath("../../config") // unit test
	viper.AddConfigPath("/app/config")  // docker
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.AutomaticEnv() // to overwrite with env var, use upper+snake case
}

func New() *Config {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
