## Pre-sales Deal Tracker
**Pre-sales Deal Tracker** is a web application designed to assist the pre-sales engineering team at Vas Technologies(A small networking company, where I am a network intern) in managing deal processes efficiently. The application features a simple monolithic architecture based on a client-server model and communicates using RESTful APIs.
## Architecture
- Monolithic Application: The application is built as a single unified codebase, where all components are tightly integrated and deployed as a single unit(PS. The frontend and backend could be in different repos, but I wanted it this way)
- Client-Server Model: The frontend (client) interacts with the backend(server) through a well-defined set of RESTful APIs. The backend handles the applications logic, processes client requests and interacts with the backend to get the desired data.
- RESTful API: The backend exposes endpoints for various operations, allowing the frontend to fetch, create, update, delete and process data as needed.
## Backend:
- The backend is responsible for processing requests from the frontend. Key responsibilities include:
    - Data Management: handles the CRUD operations for deals, customers, users and any related entities
    - Business Logi: implements the logic to manage deals and track their progress through various stages, generating a means to make data driven decisions about deal negotiations in the future.
    - API Endpoints: Provides RESTful API endpoints to supply data to the frontend, ensuring seamles communication between the client and server
## Technologies Used:
- Programming Languages: Go, Node.js(typescript)
- Frameworks: Gin, Next.js
- Database: PostgreSQL
- UI & components: Tailwind, Schdcn UI