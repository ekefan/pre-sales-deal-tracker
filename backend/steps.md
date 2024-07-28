# Step By Step process taken to write this application

1. Create migrations using go migrate in package db
2. Populate up migration script with init schema
3. Populate down migration script based on SQL queries in up migration

3. Create makefile to keep reqular commands
4. Write SQL Queries
5. Setup postgres database on docker