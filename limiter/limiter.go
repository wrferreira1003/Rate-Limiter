package strategy

import (
	"fmt"
	"time"

	"github.com/wrferreira1003/Rate-Limiter/config"
)

type RateLimiter struct {
	strategy RateLimiterStrategy
	config   *config.Config
}

func NewRateLimiter(strategy RateLimiterStrategy, config *config.Config) *RateLimiter {
	return &RateLimiter{
		strategy: strategy,
		config:   config,
	}
}

// CheckRequest faz toda a lógica do rate limiting
func (rl *RateLimiter) CheckRequest(ip, token string) error {
	var key string
	var limit, blockDuration int

	// Se foi passado um token válido, sobrepõe a lógica de IP
	if token != "" {
		key = fmt.Sprintf("token:%s", token)

		// Se estiver configurado um limite customizado para esse token, usar.
		if customLimit, ok := rl.config.CustomTokenLimits[token]; ok {
			limit = customLimit
		} else {
			limit = rl.config.TokenRequestLimit
		}
		blockDuration = rl.config.TokenBlockDuration
	} else {
		// Caso contrário, usar o limite de IP
		key = fmt.Sprintf("ip:%s", ip)
		limit = rl.config.IpRequestLimit
		blockDuration = rl.config.IpBlockDuration
	}

	// 1. Verifica se está bloqueado
	blocked, err := rl.strategy.IsBlocked(key)
	if err != nil {
		return err
	}
	if blocked {
		return fmt.Errorf("you have reached the maximum number of requests or actions allowed within a certain time frame")
	}

	// 2. Incrementa o contador
	count, err := rl.strategy.IncrementCounter(key)
	if err != nil {
		return err
	}

	// 3. Se for a primeira vez (count == 1), define um TTL de 1 segundo (ou outro)
	//    para "descontar" esse request do bucket
	if count == 1 {
		rl.strategy.SetExpiration(key, 1*time.Second)
	}

	// 4. Se estourar o limite, bloquear
	if count > limit {
		// Bloqueia por blockDuration (segundos)
		rl.strategy.SetBlocked(key, time.Duration(blockDuration)*time.Second)
		return fmt.Errorf("you have reached the maximum number of requests or actions allowed within a certain time frame")
	}

	// Se não estourou, ok
	return nil
}
