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
