# TODO: dockerize your web server by writing a Dockerfile for it.
# TODO: put together a docker-compose.yml file. For now, you'll use it to start the application locally, then we'll use it for deploy or end2end tests.
postgres:
	docker run --name postgresDealTracker -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=vasDealTracker -d postgres

createdb:
	docker exec postgresDealTracker createdb --username=root  --owner=root dealTrackerDB

dropdb:
	docker exec postgresDealTracker dropdb dealTrackerDB

migrateup1:
	migrate -path internal/db/migrations -database "postgresql://root:vasDealTracker@localhost:5432/dealTrackerDB?sslmode=disable" -verbose up
migrateforce:
	migrate -path internal/db/migrations -database "postgresql://root:vasDealTracker@localhost:5432/dealTrackerDB?sslmode=disable" force 1

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb
