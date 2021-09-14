-include .env

copy-env:
	@cp .env-sample .env

db-setup:
	@chmod +x ./scripts/db.sh
	./scripts/db.sh

cli-migrate:
	migrate create -dir migrations -ext sql -seq $(MIGRATION_NAME)

cli-migrate-up:
	migrate -database ${POSTGRESQL_URL} -path migrations up

cli-migrate-down:
	migrate -database ${POSTGRESQL_URL} -path migrations down