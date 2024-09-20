# API for Pre-sales Deal Tracker

## BRIEF SUMMARY FOR API CONCEPT

- For this application, an admin creates a user with a sales role.
- The sales user can submit pitch (customer) requests to the admin,
- The pitch requests contains mainly of a list of services the customer is requesting for, the task the sales users is asking of the pres-sales-team to meet up with the customer request, the id of the sales user, the deadline for the pre-sales engineer to work on the task.
- The admin, when he views the pitch request, can add a reminder to his calender basd on the pitch request dead, can create a deal based on this pitch request, then updated its status as the deal progress.
- When needed, the admin can delete users, and deals
- The sales user can delete pitch requests associated with their id.

## Prerequisites to Run the Server on Your Local Computer

1. **Install Docker**  
   Visit [Docker Docs](https://docs.docker.com/engine/install/) for clear instructions on how to install Docker.

2. **Clone the repository**  
   Clone this repo to your local computer, `cd` into the project directory, then `cd` into the backend folder.

   ```bash
   git clone https://github.com/ekefan/pre-sales-deal-tracker.git
   cd pre-sales-deal-tracker/backend
   ```

### Starting the Server

To spin up the server along with the Postgres database server, use Docker Compose to build and start the server's containers:

```bash
docker-compose up
```

### API Documentation (Swagger)

The API is documented using **Swagger**. You can use the Swagger file by:

- Uploading the `Swagger.yml` file located in the backend folder to the [Swagger Editor](https://editor.swagger.io/).  
This will allow you to view the API specification in a user-friendly interface and test the API endpoints using the interactive UI.
**The server's containers must be running to be able to test the endpoints on swagger**
