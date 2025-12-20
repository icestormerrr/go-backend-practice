install:
	go mod tidy

run:
	go run ./cmd/server

build:
	go build ./cmd/server

test:
	go test ./... -v

gen-oas:
	swag init -g cmd/server/main.go -o docs

validate-oas:
	npx swagger-cli validate docs/swagger.yaml