# README

## Pre-sales Deal Tracker

### Prerequisites to run the server on your local computer

- Clone the repo to your local computer, cd into the repo then cd into the backend folder.
- To successfully start the server you need to the following:
    1. Install Docker. Visit [Docker docs](https://docs.docker.com/engine/install/)  for instructions on how to install

    2. Create and run a postgres docker container, the command has been defined in the make file

        ```bash
        make postgres
        ```

    3. Create a database in that postgres container

        ```bash
        make createdb
        ```

    4. Install go migrate. Visit [Go migrate cli docs](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) for a proper documentation

        ```bash
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

        ```

    5. Migrate db schema to database: to create tables in the database

    ```bash
    make migrateup1
    ```

- Check the Makefile for other relevant commands

- To start the server, run:

    ```bash
    go run main.go
    ```

The postgres container must be running for the server to connect to the database.

A `.env.example` file is provided to show how your env should look like
