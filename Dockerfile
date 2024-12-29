# Etapa 1: Dependências
FROM golang:1.23.1-alpine AS dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Etapa 2: Builder (para produção)
FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY --from=dependencies /app /app
COPY . . 
COPY ./cmd/server/.env /app/cmd/server/.env
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server -ldflags="-s -w" ./cmd/server

# Etapa 3: Desenvolvimento com CompileDaemon
FROM golang:1.23.1-alpine AS dev
WORKDIR /app/src
RUN go install github.com/githubnemo/CompileDaemon@latest
COPY --from=dependencies /app /app
COPY . .
CMD ["CompileDaemon", "--directory=/app/src", "--build=go build -o /app/server ./cmd/server", "--command=/app/server"]


# Etapa 4: Imagem Final para Produção
FROM alpine:latest AS final
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
