# BANCO - A crude bank backend API

## FUNCTIONALITY
- Create and manage user accounts
- Record changes to bank balance - entry in db for each change
- Money transaction - maintain consistency in a transaction

### REQUIREMENTS
- Go
- Docker
  - Postgres
- Make

### 1. DATABASE SETUP
![db design](dbdesign.png)
- Bring up postgres docker container and create user and db for banco
```bash
make db-setup
```
- Create a new migration file (both up and down)
```bash
make MIGRATION_NAME=some_name cli-migrate
```
- Run migrations
```bash
make cli-migrate-up
```
- Revert migrations
```bash
make cli-migrate-down
```
- Use sqlc to generate crud code, specify migrations and queries rest will be taken care
```
make sqlc-generate
```

### Making money transaction
- Transaction in databases is a very small unit of a program and it may contain several lowlevel tasks.
For this application, the transaction includes:
1. Create a transfer record with amount 10
2. Create an entry for account1 with -10
3. Create an entry for account2 with +10
4. Subtract 10 from balance of account1
5. Add 10 to balance of account2