# 🚀 Go Starter Boilerplate - Clean Layered Architecture

[![Go Version](https://img.shields.io/badge/Go-1.25.4+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin-008080?style=flat)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/ORM-GORM-blue?style=flat)](https://gorm.io/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A high-performance, scalable, and professional Go boilerplate built with **Clean Layered Architecture**. This starter kit is designed to provide a solid foundation for enterprise-grade backend applications, featuring dependency injection, modularity, and robust CI/CD pipelines.

---

## 📑 Table of Contents
- [Core Architecture & Theory](#-core-architecture--theory)
- [Tech Stack](#-tech-stack)
- [Key Features](#-key-features)
- [Project Structure](#-project-structure)
- [CLI Tools & Automated DI](#-cli-tools--automated-di)
- [Setup & Installation](#-setup--installation)
- [Environment Variables](#-environment-variables)
- [Documentation (Swagger)](#-documentation-swagger)
- [CI/CD & Deployment](#-cicd--deployment)
- [SOLID Principles](#-solid-principles)

---

## 🏗️ Core Architecture & Theory

This project follows the **Clean Layered Architecture** pattern, ensuring high maintainability and testability by decoupling business logic from external dependencies (DB, API, Frameworks).

### The Four Layers:
1.  **Transport / Router**: Defines API endpoints and handles HTTP protocol specifics using Gin.
2.  **Controller**: Entry point for requests. Validates input, calls services, and formats JSON responses.
3.  **Service (Usecase)**: The heart of the application. Contains business logic and orchestrates data flow between repositories.
4.  **Repository (Data Access)**: Handles database interactions using GORM. Isolated from business logic.

### 🔌 Dependency Injection (Provider Pattern)
We use a **Provider Pattern** for Dependency Injection, managed through the `provider/` directory. While the wiring is handled in Go, this project features an **Automated DI Engine** that generates this wiring for you.

Benefits:
-   **Singleton Management**: Ensures single instances of services/repositories.
-   **Automated Wiring**: No more manual `NewService(repo1, repo2, ...)` maintenance.
-   **Cycle Detection**: The induction engine automatically detects and prevents circular dependencies.
-   **Decoupling**: Layers remain isolated through interfaces.

---

## 🛠️ Tech Stack

-   **Language**: [Go (Golang)](https://golang.org/)
-   **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
-   **ORM**: [GORM](https://gorm.io/) with PostgreSQL
-   **Database**: [PostgreSQL](https://www.postgresql.org/)
-   **Authentication**: [JWT (JSON Web Token)](https://jwt.io/)
-   **Payment Gateway**: [Xendit](https://www.xendit.co/)
-   **Object Storage**: [Supabase Storage](https://supabase.com/storage)
-   **Documentation**: [Swagger (swaggo)](https://github.com/swaggo/swag)
-   **Containerization**: [Docker](https://www.docker.com/)

---

## ✨ Key Features

-   **Lock-tight Auth**: Complete Authentication flow (Login, Register, Email Verification, Forgot Password).
-   **Role-Based Access**: Specialized routes for Admin and User roles.
-   **Payment Ready**: Integrated with Xendit for seamless payment processing and webhooks.
-   **Cloud Storage**: Native support for Supabase storage for file uploads.
-   **Automated Migrations**: Database schema automatically synchronizes on startup.
-   **Standardized Responses**: Unified JSON response structure for success and error handling.
-   **Modular Routing**: Cleanly separated route definitions per module.

---

## 📁 Project Structure

```text
├── cmd/                # Custom CLI tools (e.g., Swagger compiler)
├── config/             # Configuration loaders & schema definitions
├── controllers/        # Request handlers (Input validation, Response formatting)
├── middleware/         # Gin middlewares (Auth, Logging, Gzip)
├── models/             
│   ├── dto/            # Data Transfer Objects (Request/Response structs)
│   ├── entity/         # GORM database models
│   └── error/          # Custom error definitions
├── provider/           # Dependency Injection & Bootstrapping
├── repositories/       # Database access layer (SQL logic)
├── router/             # Route definitions & grouping
├── services/           # Business logic layer
├── swagger/            # Automatically generated Swagger UI files
├── utils/              # Cross-cutting helpers (Response helpers, string utils)
├── main.go             # Application entry point
├── Dockerfile          # Multi-stage production container
└── .github/workflows/  # CI/CD Pipelines
```

---

## 🛠️ CLI Tools & Automated DI

This boilerplate features a unique **PowerShell-based Dependency Injection Engine** located in the `cmd/` directory. These scripts scan your code, resolve dependencies, and automatically generate the necessary boilerplate code in the `provider/` package.

### 💉 Automated Injection Scripts
| Script | Purpose |
| :--- | :--- |
| `do_inject_config.ps1` | Scans configuration files and wires them into the Config Provider. |
| `do_inject_repository.ps1` | Discovers GORM repositories and updates `repositories_provider.go`. |
| `do_inject_services.ps1` | Performs **topological sorting** on services to resolve their dependencies in the correct order. |
| `do_inject_controllers.ps1` | Wires services into controllers and updates the Controller Provider. |
| `do_inject_middleware.ps1` | Manages registration of custom Gin middlewares. |

### 🧠 How the DI Engine Works:
1.  **Scanning**: The scripts use regex to find constructor functions (e.g., `NewAccountService`).
2.  **Resolution**: It identifies the required parameters (Repositories, other Services, or Configs).
3.  **Topological Sort**: (For Services) It builds a dependency graph and sorts them so dependencies are initialized before they are needed.
4.  **Codegen**: It writes a clean, standardized `provider/*.go` file with all the wiring logic.

### 📝 Swagger Compilation
- **`compile_swagger.ps1`**: This script wraps the `swag init` command with optimized flags (`--parseDependency`, `--parseInternal`) to ensure your Swagger documentation is always complete and up-to-date.

---

## 🚀 Setup & Installation

### Prerequisites
-   Go 1.25.4 or higher
-   PostgreSQL
-   Docker (Optional, for containerization)

### Step 1: Clone the Repository
```bash
git clone <your-repo-url>
cd <your-directory-name>
```

### Step 2: Environment Setup
Copy the example environment file and fill in your credentials:
```bash
cp .env.example .env
```

### Step 3: Install Dependencies
```bash
go mod download
```

### Step 4: Run the Application
```bash
go run main.go
```

---

## � Environment Variables

| Variable | Description |
| :--- | :--- |
| `DB_HOST` | PostgreSQL Host address |
| `DB_USER` | Database username |
| `DB_PASSWORD` | Database password |
| `DB_PORT` | Database port (default 5432) |
| `DB_NAME` | Name of the database |
| `JWT_SECRET_KEY` | Secret key for signing JWT tokens |
| `XENDIT_API_KEY` | Your Xendit Secret Key |
| `HOST_PORT` | Port for the Go server to listen on |

---

## 📖 Documentation (Swagger)

This project uses `swaggo` to generate interactive API documentation.

### How to Compile Swagger:
Run the provided script to synchronize your code comments with the documentation:
```bash
./cmd/compile_swagger
```

### Accessing Swagger UI:
Once the app is running, visit:
`http://localhost:<PORT>/swagger/index.html`

---

## 🚢 CI/CD & Deployment

### GitHub Workflows
The project includes automated pipelines in `.github/workflows/`:
-   **`go-build-test.yml`**: Automatically runs tests and builds the binary on every push.
-   **`deploy-production.yml` / `deploy-development.yml`**: CD pipelines for automated shipping to target environments.
-   **`uptime-check.yml`**: Specialized workflow to monitor service health.

### Docker
We use a **Multi-Stage Dockerfile** to ensure the production image is as small and secure as possible.

**Build the image:**
```bash
docker build -t go-starter .
```

**Run the container:**
```bash
docker run -p 8080:8080 --env-file .env go-starter
```

---

## 💎 SOLID Principles

-   **S**: Single Responsibility. Each layer (Controller, Service, Repo) does exactly one thing.
-   **O**: Open/Closed. Services are open for extension via interfaces but closed for modification.
-   **L**: Liskov Substitution. Implementations are interchangeable via Provider interfaces.
-   **I**: Interface Segregation. Slim, specific interfaces for each capability.
-   **D**: Dependency Inversion. High-level modules don't depend on low-level modules; both depend on abstractions.

---

Developed with ❤️ by Abdan Hafidz.

