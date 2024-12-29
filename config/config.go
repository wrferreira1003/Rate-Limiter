package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          string `mapstructure:"REDIS_PORT"`
	IpRequestLimit     int    `mapstructure:"IP_REQUEST_LIMIT"`
	IpBlockDuration    int    `mapstructure:"IP_BLOCK_DURATION"`
	TokenRequestLimit  int    `mapstructure:"TOKEN_REQUEST_LIMIT"`
	TokenBlockDuration int    `mapstructure:"TOKEN_BLOCK_DURATION"`

	// Esse é lido diretamente pelo viper.Unmarshal
	CustomTokenLimitsString string `mapstructure:"CUSTOM_TOKEN_LIMITS"`

	// você popula manualmente depois de parsear a string.
	CustomTokenLimits map[string]int `mapstructure:"-"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg *Config
	viper.SetConfigType("env")  // Tipo do arquivo de configuração
	viper.AddConfigPath(path)   // Diretório onde o arquivo está localizado
	viper.SetConfigName(".env") // Nome do arquivo de configuração
	viper.AutomaticEnv()        // Carrega as variáveis de ambiente

	// Carrega as configuracoes do arquivo .env
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("erro ao carregar as configuracoes do arquivo .env: %w", err)
	}

	// Converte as configurações do arquivo .env para a struct Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("erro ao decodificar as configurações: %w", err)
	}

	// Agora, parseamos manualmente o campo `CustomTokenLimitsString`
	cfg.CustomTokenLimits = parseCustomTokenLimits(cfg.CustomTokenLimitsString)

	return cfg, nil
}

// parseCustomTokenLimits("token1=50,token2=20") -> {"token1": 50, "token2": 20}
func parseCustomTokenLimits(envVal string) map[string]int {
	result := make(map[string]int)
	if envVal == "" {
		return result
	}

	pairs := strings.Split(envVal, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			token := kv[0]
			limit, err := strconv.Atoi(kv[1])
			if err == nil {
				result[token] = limit
			}
		}
	}
	return result
}
