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
- Bring up postgres docker container and configure it
```bash
make db-setup
```
#### TODO
- split migration.sql to multiple migration scripts
- organize the migrations in a folder
- move some db env variables to .env file

- add colors to make commands
- all scripts in one folder