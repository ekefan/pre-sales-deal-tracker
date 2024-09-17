<!-- TODO: there are some errors in the Markdown syntax. Please fix them. Maybe you can also install a Markdown linter to help you catching issues. -->
## Pre-sales Deal Tracker

**Pre-sales Deal Tracker** is a web application designed to assist the pre-sales engineering team at Vas Technologies (a small networking company where I am a network intern) in managing deal processes efficiently. The application features a monolithic architecture based on a client-server model and communicates using RESTful APIs.

The entire application revolves around the pre-sales engineer (the admin). But the flow of the application starts after a sales rep delivers a customer (pitch) request. However, all sales rep users must be created by the admin first.

- A Sales rep logs in into the app and submit the pitch request to the pre-sales engineer.
- The engineer views the request and creates an ongoing deal based on the pitch request, then continues to update the deal based on actions taken by the engineer or any stakeholder involved.
- Based on meetings with the customer, the sales rep can update the pitch request, which would lead to updating the deal status.
- the admin can create, read/get, update, and delete users; create deals; update and delete too; and also update pitch requests.
- A manager's role is to view/monitor deals and check through deals based on certain predefined queries.
- A sales rep can view deals and pitch requests for which sales rep_id or name is the same.

For proper API documentation, check [here](/backend/swagger.yml), and use [Swagger Editor](https://editor.swagger.io/) to view the documentation with a good UI.

### Architecture

Monolithic Application: The application is built with a unified codebase, where all components are integrated. (PS. The frontend and backend could be in different repos, but I wanted it this way.)
Client-Server Model: The frontend (client) interacts with the backend (server) through a well-defined set of RESTful APIs. The backend handles the application logic, processes client requests, and interacts with the backend to get the desired data.
RESTful API: The backend exposes endpoints for various operations, allowing the frontend to fetch, create, update, delete, and process data as needed.

### Backend:

The backend is responsible for processing requests from the client. Key responsibilities include:
  Data Management: handles the CRUD operations for deals, customers, users, and any related entities.
  Business Logi: implements the logic to manage deals and track their progress through various stages, generating a means to make data-driven decisions about deal negotiations in the future.
  API Endpoints: Provides RESTful API endpoints to supply data to the frontend, ensuring seamless communication between the client and server. This endpoint includes tasks for creating any of these models: a user, pitch request and deal, reading or getting any of the models, updating them, and deleting them where necessary.

### Frontend:

Working on it...

### Technologies Used:

Programming Languages: Go, Node.js (typescript)
Frameworks: Gin, Next.js
Database: PostgreSQL
UI & components: Tailwind, Shadcn UI
