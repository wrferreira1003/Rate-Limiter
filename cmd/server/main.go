package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wrferreira1003/Rate-Limiter/config"
	strategy "github.com/wrferreira1003/Rate-Limiter/limiter"
	"github.com/wrferreira1003/Rate-Limiter/middleware"
)

func main() {
	config, err := config.LoadConfig("/app/cmd/server")

	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}
	fmt.Println("Configuração carregada com sucesso:", config)

	redisStrategy := strategy.NewRedisStrategy(config.RedisHost, config.RedisPort)
	rateLimiter := strategy.NewRateLimiter(redisStrategy, config)

	//Criar o servidor
	r := mux.NewRouter()

	//Adicionar o middleware de rate limiting
	r.Use(middleware.RateLimiterMiddleware(rateLimiter))

	//Adicionar as rotas
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Status OK"))
	})

	//Sobe o servidor
	addr := fmt.Sprintf(":%s", config.Port)
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
