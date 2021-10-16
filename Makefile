db-setup:
	@chmod +x ./scripts/db.sh
	./scripts/db.sh

cli-migrate:
	migrate create -dir db/migrations -ext sql -seq $(MIGRATION_NAME)

cli-migrate-up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up

cli-migrate-down:
	migrate -database ${POSTGRESQL_URL} -path db/migrations down

sqlc-generate:
	sqlc -f db/sqlc/sqlc.yaml generate

test:
	go clean -testcache
	go test -v -cover ./...

server:
	go run main.go

.PHONY: test, server