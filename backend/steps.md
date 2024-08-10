# Step By Step process taken to write this application

1. Create migrations using go migrate in package db
2. Populate up migration script with init schema
3. Populate down migration script based on SQL queries in up migration

4. Create makefile to keep reqular commands
5. Write SQL Queries
6. Use SQLC to generate db querier and sql connection
7. Setup postgres database on docker
8. Run migration to create public schema on database
9. Write unit test on db query functions to ensure queries work as intended
10. Setup gin server, with db engine for router, the db interface to make queries and other needed fields within a server instance or context
11. Write handler functions for each endpoint
12. Test handler functions work as intended
13. Optimise handler functions
14. Add viper to marshall envfiles to app server
15. Use bcrypt to hash and compare password hash
16. Create custom validator for role in http requests for user
17. Adjust query to filterdeals
18. Creating Payload file for JWT authentication
19. Define the payload data using a struct
