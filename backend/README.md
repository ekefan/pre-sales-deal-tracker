# API for Pre-sales Deal Tracker

## APPLICATION CONCEPT

With the amount of deals comming in every month at Vas Technologies, Vas Deal Tracker is to help the pre-sales team optimise their responses and workflow to customer requests.

### The workflow

- The application is started with a predefined admin user, who is a pre-sales engineer
- This user, just like every other user should update their password to continue using the application
- An admin can create other users, update them, reset their password and even delete them when necessary
- A sales user can create a pitch request, update that pitch request and delete.
- A sales user can only interact with pitch requests they create.
- An admin, can initiate the action to create a deal, after they have viewed a pitch request, they can update the deal and delete it

API end points have been designed in the swagger documentation, view it on swagger UI or import on POSTMAN.

All end points require an Authorization header with BearerAuth tokens except the /auth/login end point from which the auth tokens are generated.

<!-- IVAN: What do you think of the the endpoints now???  -->
<!-- // FIXME: we've to discern about the "users" resource. We have three options: -->
<!-- // 1. "users" are only used to keep track of logins. Stick to addresses like "/register", "login", "logout", "password-reset", and so on. -->
<!-- // 2. "users" are part of the domain of our application. In this case use "GET /users", "POST /users", "PUT /users/:id", and so on. -->
<!-- // 3. Mix & Match: clearly separate those endpoints: the one used for signing-up, login, logout, ecc. from the ones used as domain in our application. Potentially, we could also have separated table. -->