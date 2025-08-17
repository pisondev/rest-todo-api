# REST To-Do API

A simple RESTful API for managing a personal to-do list. Built with Go and MySQL, it supports full CRUD functionality and JWT for user authentication.

---

## Features

- **User Authentication:** Complete JWT-based authentication (Register & Login).
- **CRUD Operations:** Full Create, Read, Update, Delete functionality for tasks.
- **Soft Deletes:** Tasks are soft-deleted to preserve data integrity and allow for recovery.
- **Task Filtering:** Users can filter their tasks by status and due date.
- **Layered Architecture:** Follows a clean architecture (controller, service, repository) for separation of concerns.

---

## Tech Stack

- **Language**: Go
- **Web Framework**: Fiber
- **Database**: MySQL (w/ `go-sql-driver/mysql`)
- **Authentication**: JWT (w/ `golang-jwt/jwt/v5`)
- **Validation**: `go-playground/validator/v10`
- **Environment Config**: `.env` w/ `github.com/joho/godotenv`

---

## Folder Structure

```
rest-todo-api/
├── app/         # Database connection & config
├── controller/  # HTTP request handlers
├── exception/   # Central error handler
├── helper/      # Common utilities
├── middleware/  # JWT Auth middleware
├── model/       # Domain & web models (structs)
├── repository/  # DB query logic
├── service/     # Business logic
├── .env         # Config file
├── .gitignore
├── go.mod
├── go.sum
├── main.go      # Main func (router, dependency injection)
├── openapi.json # API spec documentation
└── README.md    # Project info
```

---

## Getting Started

### 1. Clone the repo
```bash
git clone [https://github.com/pisondev/rest-todo-api.git](https://github.com/pisondev/rest-todo-api.git)
cd rest-todo-api
```

### 2. Create a `.env` file
Create a file named `.env` in the root of the project and add your database credentials and a JWT secret.

```env
SERVER_PORT=":..."
ALLOWED_ORIGIN="..."

DB_USER="root"
DB_PASS="your_mysql_password"
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="rest_todo_api"
DB_PARAMS="parseTime=true&loc=UTC"

JWT_SECRET_KEY="your-super-secret-key"
```

> ⚠️ Don't commit your `.env` file — it should be in your `.gitignore`.

### 3. Start Dependencies
This project requires a running MySQL instance. You can use a local installation or start one easily with Docker.
```bash
# This command will start a MySQL 8 container
docker run -d --name mysql-todo -p 3306:3306 -e MYSQL_ROOT_PASSWORD=your_mysql_password -e MYSQL_DATABASE=rest_todo_api mysql:8
```

### 4. Install Go Modules
```bash
go mod tidy
```

### 5. Run the app
```bash
go run main.go
```

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/register` | Register a new user. |
| POST | `/api/login` | Login and receive a JWT. |
| POST | `/api/tasks` | Create a new task. (Auth required) |
| GET | `/api/tasks` | Get all tasks for a user (with filters). (Auth required) |
| GET | `/api/tasks/{taskId}` | Get a specific task by its ID. (Auth required) |
| PATCH | `/api/tasks/{taskId}` | Update a task by its ID. (Auth required) |
| DELETE | `/api/tasks/{taskId}` | Soft delete a task by its ID. (Auth required) |

> Filtering for `GET /api/tasks` is supported via query params: `status`, `due_date`

---

## TODO / Roadmap

- [ ] **Pagination:** Add pagination to the `GET /api/tasks` endpoint.
- [ ] **Advanced Filtering:** Add filtering by task title (search).
- [ ] **User Roles:** Implement user roles (e.g., standard user vs. admin).
- [ ] **Tests and CI/CD:** Add unit and integration tests for the service and handler layers.
- [ ] **Deployment:** Containerize the Go application with Docker for easy deployment.

---

## License

MIT — free to use and modify.