# Build stage
FROM golang:1.25-alpine AS builder

# Instala dependências do sistema
RUN apk add --no-cache git

# Define diretório de trabalho
WORKDIR /app

# Copia go.mod e go.sum e baixa dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código-fonte
COPY . .

# Compila o binário
RUN go build -o consumer main.go

# Runtime stage
FROM alpine:3.18

# Define diretório de trabalho
WORKDIR /app

# Copia binário e .env
COPY --from=builder /app/consumer .
COPY .env .env

# Expõe nenhuma porta (consumer é cliente)
# Define comando padrão
CMD ["./consumer"]
