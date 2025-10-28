include .env
export

DB_STRING=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

migrate-up:
	goose -dir ./migrations postgres "$(DB_STRING)" up
migrate-status:
	goose -dir ./migrations postgres "$(DB_STRING)" status