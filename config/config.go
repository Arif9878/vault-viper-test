package config

import (
	"io"

	"github.com/spf13/viper"
)

type Config struct {
	Imam string `mapstructure:"imam"`
}

func LoadConfig(path io.Reader) (config Config, err error) {
	viper.SetConfigType("json")
	err = viper.ReadConfig(path)
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
