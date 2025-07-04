package config

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type ImmutableConfigs interface {
	GetPort() int
	GetDBConf() DB
}

type config struct {
	Port int `mapstructure:"APP_PORT"`
	DB   DB  `mapstructure:"DB"`
}

type (
	DB struct {
		Host             string `mapstructure:"HOST"`
		Port             string `mapstructure:"PORT"`
		Name             string `mapstructure:"NAME"`
		User             string `mapstructure:"USER"`
		Password         string `mapstructure:"PASSWORD"`
		MaxIdleConns     int    `mapstructure:"MAX_IDLE_CONNS"`
		MaxOpenConns     int    `mapstructure:"MAX_OPEN_CONNS"`
		MaxLifetimeConns int    `mapstructure:"MAX_LIFETIME_CONNS"`
		SSLMode          string `mapstructure:"SSL_MODE"`
	}
)

func (im *config) GetPort() int {
	return im.Port
}

func (im *config) GetDBConf() DB {
	return im.DB
}

var (
	once sync.Once
	conf *config
)

func NewImmutableConfigs() ImmutableConfigs {
	once.Do(func() {
		v := viper.New()
		appEnv, exists := os.LookupEnv("APP_ENV")
		configName := "app.config.local"
		if exists {
			switch appEnv {
			case "development":
				configName = "app.config.dev"
			default:
				configName = "app.config.local"
			}
		}

		slog.Debug("Config loaded", slog.String("ConfigName", configName), slog.String("Level", "warn"))

		v.SetConfigName("configs/" + configName)
		v.AddConfigPath(".")

		v.SetEnvPrefix("INTERNAL_TRANSFER")
		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			panic(err.Error())
		}

		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		err := v.Unmarshal(&conf)
		if err != nil {
			panic(err.Error())
		}
	})
	return conf
}
