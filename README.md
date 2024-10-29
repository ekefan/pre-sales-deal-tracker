# Pre-sales Deal Tracker

**Pre-sales Deal Tracker** is deal tracker and management system. Designed to optimise sales-client processes while integration  pre-sales-management processes.

## Backend

The backend is built with golang and it interacts with a postgres database, Key responsibilities include:
    Data Management: handles the CRUD operations for deals, users and pitch requests.
    Business Logic: implements the logic to manage deals and track their progress through various stages, generating a means to make data-driven decisions about deal negotiations in the future.
    API Endpoints: Provides RESTful API endpoints to supply data to the frontend, ensuring seamless communication between the client and server. This end points include tasks for creating any of these resources: a user, pitch request and deal, reading them to visualize or perform calcaluations
    updating them, and deleting them where necessary.

View a more detailed documentation, [here](/backend/README.md).

## Frontend

A beautiful UI designed to allow the users to interact with the application seamlessly.
_Still under construction_

### Technologies Used

Programming Languages: Go, Node.js (typescript)
Frameworks: Gin, Next.js
Database: PostgreSQL
UI & components: Tailwind, Shadcn UI
