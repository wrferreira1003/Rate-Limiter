package middleware

import (
	"fmt"
	"net/http"
	"strings"

	strategy "github.com/wrferreira1003/Rate-Limiter/limiter"
)

func RateLimiterMiddleware(rl *strategy.RateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrair IP
			clientIP := r.RemoteAddr

			//Extrair apenas o IP sem a porta
			ip := strings.Split(clientIP, ":")[0]
			fmt.Println("IP:", ip)

			// Extrair token
			token := r.Header.Get("API_KEY")

			// Verificar se passou no rate limiting
			if err := rl.CheckRequest(ip, token); err != nil {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(err.Error()))
				return
			}

			// Se passou, chama o próximo
			next.ServeHTTP(w, r)
		})
	}
}
