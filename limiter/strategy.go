package strategy

import "time"

type RateLimiterStrategy interface {
	IncrementCounter(key string) (int, error)               // Incrementa o contador para a chave especificada
	GetCounter(key string) (int, error)                     // Define o valor do contador para a chave especificada
	SetExpiration(key string, duration time.Duration) error // Define a expiração do contador para a chave especificada
	IsBlocked(key string) (bool, error)                     // Verifica se a chave está bloqueada
	SetBlocked(key string, duration time.Duration) error    // Bloqueia a chave por um período de tempo
}
