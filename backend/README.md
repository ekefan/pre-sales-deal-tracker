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

### The code

To have a quick overview of the api requests and responses you can view export the swagger file to an api client or the swagger ui editor
To run the http server on your computer:

- create a fork and clone it to your local computer[find out how here](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo)
- navigate to the `backend` directory
- Build a docker image to have every service needed by the server, use the build tag to allow any pre-existing images of this server to be re-built, since there will be more features added

```bash
docker compose up -d --build
```

- You can read more about docker compose [here](https://docs.docker.com/get-started/docker-concepts/the-basics/what-is-docker-compose/)
- If the image build is successful you can send requests the server

<!-- # FIXME: -->
<!-- Done.... -->
<!-- And for the last statment, I get you man 100% -->
- The swagger UI should show a padlock close to the restricted routes. I don't remember how to do it but I know you can 😄
- In the request `models` you can mark the required fields as required. They will have a red star close to them
- The `pitch_requests` endpoint exposes a `PATCH` but should be a `PUT` since you're accepting the whole resource.
- The same as before applies for the `deals` endpoint
- The `error` response should have a `code` string field that is like a sentinel errors. Sentinel errors are the ones we're expecting to happen such as `NOT_FOUND`, `VALIDATION`, `NETWORK_CONNECTION`. Potentially, they could be shared with the client that relies logic against
- Keep this in mind, the `swagger.yml` is a contract we share with the FE team but it's not something written on the stones. It should seldomnly change.
