# Rate Limiter em Go

Este projeto é um **rate limiter** desenvolvido em Go, com a finalidade de **limitar** o número de requisições feitas a um servidor web. A aplicação permite configurar diferentes limites de requisições por **IP** e por **Token**, além de suportar um tempo de bloqueio (ban) quando esses limites são excedidos.

---

## Índice

1. [Descrição do Projeto](#descrição-do-projeto)  
2. [Arquitetura e Tecnologias](#arquitetura-e-tecnologias)  
  
---

## Descrição do Projeto

Este projeto **controla o número de requisições** que podem ser feitas a um servidor Go em um determinado período de tempo. Ele foi projetado como um **middleware**, de forma que possa ser facilmente integrado a qualquer aplicação Go que utilize HTTP.

- **Limite por IP**: Se um IP exceder X requisições em um intervalo de 1 segundo, ele é bloqueado por um tempo configurável.  
- **Limite por Token**: Se um token (enviado no header `API_KEY`) exceder Y requisições em um intervalo de 1 segundo, ele também é bloqueado pelo tempo configurável.  

A aplicação faz uso de **Redis** para armazenar os contadores e as chaves de bloqueio, mas utiliza um pattern de **Strategy** para permitir a troca de repositório com facilidade (por exemplo, poderíamos futuramente usar Memcached ou um banco SQL).

---

## Arquitetura e Tecnologias

- **Go**: Linguagem principal do projeto.
- **Redis**: Utilizado para armazenar contadores e flags de bloqueio (ex.: `blocked:token` ou `blocked:ip`).  
- **Viper** ou **godotenv**: Para gerenciar variáveis de ambiente.  
- **Docker e Docker Compose**: Para orquestrar tanto o serviço Go quanto o Redis em containers.  
- **Middleware**: O Rate Limiter atua como middleware, interceptando requisições HTTP antes do handler final.

## Teste de carga

Para testar a aplicação, você pode usar o `ab` (Apache Benchmark) para gerar carga no servidor.

```bash
ab -n 50 -c 10 -H "API_KEY: abc123" http://localhost:8080/
```

```bash
ab -n 50 -c 10 http://localhost:8080/
```

