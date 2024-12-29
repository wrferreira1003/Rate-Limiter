package config

import "github.com/spf13/viper"

type Config struct {
	Port               string `mapstructure:"PORT"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          string `mapstructure:"REDIS_PORT"`
	IpRequestLimit     string `mapstructure:"IP_REQUEST_LIMIT"`
	IpBlockDuration    string `mapstructure:"IP_BLOCK_DURATION"`
	TokenRequestLimit  string `mapstructure:"TOKEN_REQUEST_LIMIT"`
	TokenBlockDuration string `mapstructure:"TOKEN_BLOCK_DURATION"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
