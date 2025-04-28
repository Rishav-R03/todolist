📄 TodoList API - Documentation
Overview
The TodoList API is a RESTful service built using Go, Chi Router, GORM ORM, and PostgreSQL as the backend database.
It provides full CRUD operations for managing tasks (todos), aiming to demonstrate clean backend architecture, database modeling, and API design best practices.

This project is ready for local development, production deployments, and further expansion into a full-stack application.

🚀 Features
Create, Read, Update, and Delete todos

Structured code with separation of concerns

Database connection using GORM ORM

JSON-based API responses

Graceful shutdown handling

Environment variable support for configuration (to be added)

🏗️ Tech Stack

Technology	Purpose
Go (Golang)	Main backend server
Chi	Lightweight router for API endpoints
GORM	ORM for PostgreSQL
PostgreSQL	Relational database
Postman	API Testing
📂 Project Structure
go
Copy
Edit
/todolist
│
├── main.go          // Application entry point
├── go.mod, go.sum   // Go modules for dependency management
├── models/          // (future) Folder to hold DB models
├── handlers/        // (future) Folder for HTTP handlers
└── README.md        // Documentation
⚙️ Setup Instructions
1. Prerequisites
Go 1.20+

PostgreSQL 13+

(Optional) Postman for API testing

2. Local Setup
Clone the Repository
bash
Copy
Edit
git clone https://github.com/yourusername/todolist-go.git
cd todolist-go
Install Dependencies
bash
Copy
Edit
go mod tidy
Configure PostgreSQL
Make sure you have a database called todolistdb created in PostgreSQL.

Example SQL:

sql
Copy
Edit
CREATE DATABASE todolistdb;
Make sure you have a user (postgres) with correct privileges or adjust according to your environment.

Run the Application
bash
Copy
Edit
go run main.go
Server will start listening at http://localhost:9000/

🌐 API Endpoints

Method	Endpoint	Description
GET	/todo/	Fetch all todos
POST	/todo/	Create a new todo
PUT	/todo/{id}	Update an existing todo
DELETE	/todo/{id}	Delete a todo by ID
📦 Request and Response Formats
1. Create Todo
POST /todo/

Request Body:

json
Copy
Edit
{
  "title": "Learn Go",
  "completed": false
}
Response:

json
Copy
Edit
{
  "message": "Todo created successfully",
  "todo_id": "generated-id-here"
}
2. Fetch All Todos
GET /todo/

Response:

json
Copy
Edit
{
  "data": [
    {
      "id": "todo-id",
      "title": "Learn Go",
      "completed": false,
      "createdAt": "2025-04-29T10:00:00Z"
    }
  ]
}
3. Update Todo
PUT /todo/{id}

Request Body:

json
Copy
Edit
{
  "title": "Learn Go Lang",
  "completed": true
}
Response:

json
Copy
Edit
{
  "message": "Todo updated successfully"
}
4. Delete Todo
DELETE /todo/{id}

Response:

json
Copy
Edit
{
  "message": "Todo deleted successfully"
}
🛡️ Error Handling
The API responds with appropriate HTTP status codes and error messages in case of failures:


Status Code	Meaning
400 Bad Request	Invalid input or missing fields
404 Not Found	Resource not found
500 Internal Server Error	Server-side error
Example error response:

json
Copy
Edit
{
  "message": "The id is invalid"
}
🏗 Future Improvements
Add JWT Authentication (Login/Register)

Pagination and Filtering for Todo list

Dockerize the application

Add Swagger/OpenAPI documentation

Deploy on Cloud (AWS/GCP)

Integrate CI/CD pipelines

Add Unit and Integration Tests

🙌 Acknowledgments
Go Programming Language

Chi Router

GORM ORM

PostgreSQL

Postman

📢 Conclusion
This TodoList API project demonstrates a clean, scalable approach to building Go REST APIs backed by PostgreSQL.
It is ready for local development, production deployment, and future extension into more complex systems.
