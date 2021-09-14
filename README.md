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

#### INITIAL SETUP
- export all secret environment variables
```bash
make export-vars
```
#### DATABASE DESIGN & SETUP
![db design](dbdesign.png)
- Bring up postgres docker container and create user and db for banco
```bash
make db-setup
```
- Create a new migration file (both up and down)
```bash
make MIGRATION_NAME=some_name cli-migrate
```