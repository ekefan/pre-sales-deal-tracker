## README

To spin up the server in your local machine, you need docker,
Run:

1. `make postgres` to create and run a postgres docker container
2. `make createdb` to create a database in that postgres container
3. `make migrateup1` to migrate db schema to database

Before you start the server
`go run main.go`

The postgres container must be running for the application to run.
