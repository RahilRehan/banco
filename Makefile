-include app.env

run-postgres:
	docker run --name postgres-banco -p 5432:5432 -v banco-data:/var/lib/postgresql/data -e POSTGRES_USER="${POSTGRES_USER}" -e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" -d postgres

create-db:
	docker exec -it postgres-banco createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} bancodb

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