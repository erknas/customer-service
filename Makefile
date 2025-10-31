include .env
export

DB_STRING=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

build:
	@go build -o bin/customer-service cmd/customer-service/main.go
run: build
	@./bin/customer-service
	
migrate-up:
	@goose -dir ./migrations postgres "$(DB_STRING)" up
migrate-status:
	@goose -dir ./migrations postgres "$(DB_STRING)" status

proto:
	@protoc -I . -I ./plugins \
  		--go_out=pkg --go_opt=paths=source_relative \
  		--go-grpc_out=pkg --go-grpc_opt=paths=source_relative \
  		--grpc-gateway_out=pkg --grpc-gateway_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:pkg" \
  		api/customer/customer.proto

infra:
	@docker compose up -d
exec:
	@docker exec -it customer-service_db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: proto