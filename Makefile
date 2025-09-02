# Nome da imagem Docker
IMAGE_NAME=consumer-rabbitmq-to-mongo

.PHONY: build docker-run clean

# Build local do Go
build:
	go build -o consumer main.go

# Build da imagem Docker
docker-build:
	docker build -t $(IMAGE_NAME) .

# Executa container Docker em modo interativo
docker-run:
	docker run --env-file .env --name $(IMAGE_NAME)-container --rm $(IMAGE_NAME)

# Remove binário local
clean:
	rm -f consumer
