## Pre-sales Deal Tracker

**Pre-sales Deal Tracker** is a web application designed to assist the pre-sales engineering team at Vas Technologies(A small networking company, where I am a network intern) in managing deal processes efficiently. The application features a simple monolithic architecture based on a client-server model and communicates using RESTful APIs.

The entire application revolves around the pre-sales engineer (the admin). But the flow of the application starts after a sales-rep delivers a pitch request. However all users with role: sales must be created by the admin first.

- They login into the app and submit the pitch request to the pre-sales engineer along with what the customer needs and what the sales-rep needs from the engineer
- The engineer views the request and creates an ongoing deal based on the pitch request
- Based on meetings with the customer the sales-rep can update the pitch-req and the user views and updates the deal if necessary
- the admin can create, read/get, update, and delete users, create deals, update and delete too, also update pitch requests
- A managers role is only to view/monitor deals and maybe in the future make comments on deals
- A sales-rep can view deals, and pitchreqs whose sales-rep_id or name is the-same.

### Architecture

- Monolithic Application: The application is built with a unified codebase, where all components are integrated but deployed as a single units(PS. The frontend and backend could be in different repos, but I wanted it this way)
- Client-Server Model: The frontend (client) interacts with the backend(server) through a well-defined set of RESTful APIs. The backend handles the applications logic, processes client requests and interacts with the backend to get the desired data.
- RESTful API: The backend exposes endpoints for various operations, allowing the frontend to fetch, create, update, delete and process data as needed.

### Backend:

- The backend is responsible for processing requests from the client. Key responsibilities include:
  - Data Management: handles the CRUD operations for deals, customers, users and any related entities
  - Business Logi: implements the logic to manage deals and track their progress through various stages, generating a means to make data driven decisions about deal negotiations in the future.
  - API Endpoints: Provides RESTful API endpoints to supply data to the frontend, ensuring seamles communication between the client and server. This endpoints includes tasks for creating any of these models: a user, pitch-request and deal, reading or getting any the models, updating them and deleting them where necessary.

### Frontend:

- Working on it...

### Technologies Used:

- Programming Languages: Go, Node.js(typescript)
- Frameworks: Gin, Next.js
- Database: PostgreSQL
- UI & components: Tailwind, Shadcn UI
