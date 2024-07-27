postgres:
	docker run --name postgresDealTracker -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=vasDealTracker -d postgres

createdb:
	docker exec postgresDealTracker createdb -username=root dealTrackerDB

dropdb:
	docker exec postgresDealTracker drop dealTrackerDB

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb
