# Go Gin API Clean Architecture

This project is a RESTful API boilerplate built with [Go](https://golang.org/) and [Gin Web Framework](https://github.com/gin-gonic/gin), following the principles of Clean Architecture. It includes User Management, Authentication (JWT), and Swagger documentation.

## Features

- **Clean Architecture**: Separation of concerns (Domain, Usecase, Repository, Delivery).
- **User Management**: registration, profile management, and admin listing with pagination.
- **Authentication**: Secure password hashing (bcrypt) and JWT-based authentication.
- **Database**: MySQL integration with migration support.
- **Documentation**: Swagger API documentation.
- **Testing**: Unit tests with mocks.
- **CI/CD**: GitHub Actions workflow.

## Prerequisites

- [Go](https://go.dev/dl/) 1.23+
- [MySQL](https://www.mysql.com/)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (CLI tool for migrations)
- [Make](https://www.gnu.org/software/make/) (Optional, aimed at Linux/Mac/WSL users)

## Setup

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd go-gin-api
    ```

2.  **Environment Variables:**
    Copy `.env.example` to `.env` and configure your database credentials and JWT secret.
    ```bash
    cp .env.example .env
    ```
    Update the `.env` file:
    ```properties
    DB_HOST=localhost
    DB_PORT=3306
    DB_USER=root
    DB_PASS=yourpassword
    DB_NAME=go_gin_db
    JWT_SECRET=your-secure-secret-key
    ```

3.  **Database Setup:**
    Create the database in MySQL:
    ```sql
    CREATE DATABASE go_gin_db;
    ```

4.  **Database Migrations:**
    
    This project uses `golang-migrate` for managing database schemas.

    **Using Makefile (Recommended):**
    These commands automatically load your `.env` file for connection details.

    | Command | Description |
    | :--- | :--- |
    | `make migrate-up` | Apply all pending migrations |
    | `make migrate-down` | Rollback the last migration group |
    | `make migrate-force version=N` | Force the database version to N (use if dirty) |

    **Troubleshooting (Windows):**
    If `make` is not available, you can install it via Chocolatey (`choco install make`) or run the command manually:
    ```powershell
    # Apply migrations manually
    migrate -database "mysql://root:pass@tcp(localhost:3306)/go_gin_db" -path migrations up
    ```

## Running the Application
1.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

2.  **Start the Server:**
    ```bash
    # Using Makefile
    make run

    # Or standard Go command
    go run cmd/api/main.go
    ```
    The server will start on port `8080` (or as defined in `.env`).

## Development

### Common Commands

| Description | Make Command | Go Command |
| :--- | :--- | :--- |
| **Run App** | `make run` | `go run cmd/api/main.go` |
| **Run Tests** | `make test` | `go test -v ./...` |
| **Build App** | `make build` | `go build -o bin/api cmd/api/main.go` |

### Creating Migrations

To create a new migration pair (up/down):

- **Makefile:** `make migrate-create name=create_users`
- **Windows (Powershell):** `.\migrate-create.ps1 -name create_users`

### Running Migrations

**Using Makefile (Linux/Mac/WSL):**
`make` is a build automation tool standard on Unix-like systems. It is not installed by default on Windows.

- `make migrate-up`
- `make migrate-down`

**Using PowerShell (Windows):**
We have provided a native PowerShell script `migrate.ps1` for Windows users.

- **Up:** `.\migrate.ps1 up`
- **Down:** `.\migrate.ps1 down`
- **Force:** `.\migrate.ps1 force -Version <N>`

> [!IMPORTANT]
> **Production Safeguard**: If `APP_ENV=production` is set in your `.env`, the `down` command will fail by default to prevent data loss. Use `-ConfirmInput` to override.

## API Documentation

Access Swagger UI at: http://localhost:8181/swagger/index.html

## Project Structure

```
.
├── cmd/api/            # Main entry point
├── internal/
│   ├── bootstrap/      # App initialization
│   ├── config/         # Configuration loading
│   ├── delivery/       # HTTP Handlers & Middleware
│   ├── domain/         # Entities & Interfaces
│   ├── infrastructure/ # Database implementations
│   ├── pkg/            # Shared packages (JWT, etc)
│   └── usecase/        # Business Logic
├── migrations/         # SQL Migration files
└── docs/               # Documentation files
```
