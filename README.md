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
  - Unit testing
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
### Setup Database
- before setting running below commands, make sure docker is running
- Script creates root user for postgres and then creates a `bancoadmin` user and `bancodb` db which the app will use
- Enter details for postgres root user and password. And enter password for `bancoadmin` user
- `bancoadmin` and `bancodb` is defined in `app.env`
  ```bash
  make db-setup
  ```
### Run Migrations
- To run migrations, `POSTGRES_URL` variables must be exported
- replace {PASSWORD} with the password you entered while creating bancoadmin in previous step
```
export POSTGRESQL_URL='postgres://bancoadmin:{PASSWORD}@localhost:5432/bancodb?sslmode=disable'
cli-migrate-up
```
### Run tests
```
make test
```
### Run banco app
- make sure token is 32 in length
```
export TOKEN_SYMMETRIC_KEY=qwertyuiopasdfghjklzxcvbnmqwerty
export DB_PASSWORD={PASSWORD}
make server
```

### Use API
- Export the below collection in postman or thunder-client(vscode)
- After `create-user`, run user `user-login` to get the token to perform any other api action
[API Collection](./collection_banco.json)
