postgres:
	docker run --name dealTracker -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=vasDealTracker -d postgres

createdb:
	docker exec dealTracker createdb --username=root  --owner=root dealTrackerDB

dropdb:
	docker exec dealTracker dropdb dealTrackerDB

migrateup:
	migrate -path db/migrations -database "postgresql://root:vasDealTracker@localhost:5432/dealTrackerDB?sslmode=disable" -verbose up

migrateforce:
	migrate -path db/migrations -database "postgresql://root:vasDealTracker@localhost:5432/dealTrackerDB?sslmode=disable" force 1

sqlc:
	sqlc generate

docker:
	docker compose up -d --build
.PHONY: postgres createdb dropdb
