<!-- TODO: there are some warnings about the Markdown syntax violations. -->
<!-- Consider installing a Markdown linter. -->
## Pre-sales Deal Tracker
## README

### Prerequisite to run the server on your local computer
- Clone the repo to your local computer, cd into the repo the cd into the backend folder.
- Follow this steps to successfully start the server
    1. Docker installed. Visit [Docker docs](https://docs.docker.com/engine/install/)  for instructions on how to install
    2. Create and run a postgres docker container, the command has been defined in the make file
    ```bash
    make postgres
    ```
    3. Create a database in that postgres container
    ```bash
    make createdb
    ```
    4. Migrate db schema to database: to create tables in the database
    <!-- FIXME: this is not working. "migrate" might not be installed on the machine.-->
    <!-- Mention it and link a webpage to download it. -->
    ```bash
    make migrateup1
    ```

- Check the Makefile for other relevant commands

- To start the server, run:

    ```
    go run main.go
    ```

The postgres container must be running for the server to connect to the database.

A `.env.example` file is provided to show how your env should look like