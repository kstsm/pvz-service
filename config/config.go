package config

import (
	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

var Config config

type config struct {
	Server   Server
	Postgres Postgres
	JWT      JWT
}

type Server struct {
	Host string
	Port string
}

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

type JWT struct {
	JWTSecret string
}

func init() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		slog.Errorf("Ошибка при чтении конфигурации: %s", err)
	}

	Config = config{
		Server: Server{
			Host: viper.GetString("SRV_HOST"),
			Port: viper.GetString("SRV_PORT"),
		},
		Postgres: Postgres{
			Username: viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			DBName:   viper.GetString("POSTGRES_DB"),
		},
		JWT: JWT{
			JWTSecret: viper.GetString("SECRET_KEY"),
		},
	}
}
