package config

import (
	"github.com/spf13/viper"
	"sync"
	"teneta-tg/internal/bot"
	"teneta-tg/pkg/pgsql"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	TelegramConfig *bot.Config
	PgSQLConfig    *pgsql.Config
}

func New() (*Config, error) {
	var err error

	once.Do(func() {
		config = &Config{}

		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		if err = viper.ReadInConfig(); err != nil {
			return
		}

		telegramConfig := viper.Sub("telegram")
		pgsqlConfig := viper.Sub("pgsql")

		if err = telegramConfig.Unmarshal(&config.TelegramConfig); err != nil {
			return
		}

		if err = pgsqlConfig.Unmarshal(&config.PgSQLConfig); err != nil {
			return
		}
	})

	if err != nil {
		return nil, err
	}

	return config, nil
}
