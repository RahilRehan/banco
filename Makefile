copy-env:
	@cp .env-sample .env

db-setup:
	chmod +x ./scripts/db.sh
	./scripts/db.sh