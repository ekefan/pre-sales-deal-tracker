# Step By Step process taken to write this application

1. Create migrations using go migrate in package db
2. Populate up migration script with init schema
3. Populate down migration script based on SQL queries in up migration

3. Create makefile to keep reqular commands
4. Write SQL Queries
5. Use SQLC to generate db querier and sql connection
6. Setup postgres database on docker
7. Run migration to create public schema on database
8. Write unit test on db query functions to ensure queries work as intended
9. Setup gin server, with db engine for router, the db interface to make queries and other needed fields within a server instance or context
10. Write handler functions for each endpoint
11. Test handler functions work as intended
12. Optimise handler functions
13. Add viper to marshall envfiles to app server
