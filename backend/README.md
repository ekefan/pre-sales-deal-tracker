# API for Pre-sales Deal Tracker

## Prerequisites to Run the Server on Your Local Computer

1. **Install Docker**  
   Visit [Docker Docs](https://docs.docker.com/engine/install/) for clear instructions on how to install Docker.

2. **Clone the repository**  
   Clone this repo to your local computer, `cd` into the project directory, then `cd` into the backend folder.

### Starting the Server

To spin up the server along with the Postgres database server, use Docker Compose to build and start the server containers:

```bash
docker-compose up
```

### API Documentation (Swagger)

The API is documented using **Swagger**. You can use the Swagger file by:

- Uploading the `Swagger.yml` file located in the backend folder to the [Swagger Editor](https://editor.swagger.io/).  
This will allow you to view the API specification in a user-friendly interface and test the API endpoints using the interactive UI.
