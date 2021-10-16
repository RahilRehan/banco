# BANCO - A crude bank backend API

## FUNCTIONALITY
- Create User in the banco system
  - Each user can create multiple accounts, but accounts must have different currency
  - Only user, authenticated into banco system can manage their accounts(create, list, update, delete - CRUD)
- Transactions - money can be transferred from one user account to other
  - To perform transaction, user must be authenticated into banco system
  - User can only send money from their account 
  - Transaction can only take place between accounts of same currency
  - Each transaction is consistent

## REQUIREMENTS
- Go
- Docker
  - Postgres
- Make

## DB design/architecture
![Image](db.png)

## TECHNICAL DETAILS
- Database setup scripts in `scripts/db.sh`
- SQLC to generate models and crud code
- golang-migrate for migrations
- environment variable setup in makefile (best practice)
- Nice way to implements DB transaction for money transfer
	- transaction lock and deadlock handling while updating user account balances
	- For each transaction
  	- transfer details are stored in transfer table(fromAccount, toAccount, amount)
  	- Two entries created in entry table, how much money got added/deducted from toAccount and fromAccount
  	- Update account balance of fromAccount and toAccount
- Testing
  - Unit tests
  - Integration tests
  - testing via mocking (dependency injection - db layer is injected into api layer)
  - go-mockery is used for mocking
  - Test containers are used to run integration tests
  - There is no service layer, as it seems to be a little overkill for this project.
- In api request - custom param validator (used reflection)
- User password encryption using bcrypt 
- Use Paseto based user authentication
  - JWT authentication code is also present
  - Interface is used for Token based authentication
  - So, you can easily replace Paseto with JWT
- Github Actions is used as Pipeline 


## SETUP

### Export variables
```bash
  export POSTGRES_USER=postgres
  export POSTGRES_PASSWORD=anything
  export DB_HOST=localhost # for local setup
  export DB_HOST=host.docker.internal # for docker setup
  export TOKEN_SYMMETRIC_KEY=blahblahblahblahblahblahblahblah # 32 in length

```

### local setup
- `bancoadmin` and `bancodb` is defined in `app.env`
```bash
  make run-postgres
  make create-db
  make server
```

### Run Migrations
- Migrations are run by default when app restarts
- But, if you want to run migrations from cli, `POSTGRES_URL` variables must be exported
```
export POSTGRESQL_URL='postgres://${POSTGRES_USER}:{$POSTGRES_PASSWORD}@localhost:5432/bancodb?sslmode=disable'
cli-migrate-up
```

### Run tests
```
make test
```

### Docker setup
```bash
  export DB_HOST=host.docker.internal
  make run-postgres
  docker network create banco
  docker network  connect banco postgres-banco
  docker build -t banco .
  docker run --name banco -p 8080:8080 --network banco -e DB_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -e TOKEN_SYMMETRIC_KEY=${TOKEN_SYMMETRIC_KEY} -e DB_HOST=${DB_HOST} -e GIN_MODE=release banco
```

## Docker compose setup
```bash
export DB_HOST=host.docker.internal
```
- If you want to run any command line snippets for docker compose setup, use concept of start.sh !(more details)[https://www.youtube.com/watch?v=jf6sQsz0M1M&list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE&index=25]

- All the above instructions can be bundled together into makefile target or some script, but I like it this way as it will become very easier to debug if you know what is happening under the hood.
  
### Use API
- Export the below collection in postman or thunder-client(vscode)
- After `create-user`, run user `user-login` to get the token to perform any other api action
[API Collection](./collection_banco.json)
