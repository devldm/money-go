# money-go

> A simple backend service for money transfer between users, written in Go.

---

## Features

- **User Management**:  
  - Retrieve user information by ID.
  - Maintain user balances.

- **Money Transfer**:  
  - Send money between users.
  - Automatic balance updates for sender and receiver.
  - Keep a record of transactions with status and timestamps.

- **Transaction History**:  
  - Query a user's sent and received transactions.
  - Paginated transaction history support.

- **gRPC API**:  
  - Exposes gRPC endpoints for all main operations.
  - Service definitions in `api/v1/` (see proto files for details).

- **Persistence**:  
  - Uses SQLite for lightweight storage.
  - Database connection and migration handled automatically.

- **Dockerized**:  
  - Ready-to-run via Docker (see below).

---

## Technologies Used

- **Language:** Go (`golang`)
- **gRPC**: Communication protocol & service interface
- **SQLite**: Embedded database for local development
- **sqlx**: SQL toolkit for Go
- **Protocol Buffers**: API schema definition
- **Docker**: Containerized deployment

---

## Getting Started

### Prerequisites

- Go 1.24+ installed, or Docker
- `protoc` and Go plugins for proto compilation (optional, for API edits)

### Running Locally

#### 1. Clone the repository

```sh
git clone https://github.com/devldm/money-go.git
cd money-go
```

#### 2. Build & Run with Docker

```sh
docker build -t money-go .
docker run -p 50051:50051 money-go
```

#### 3. Or Run Directly (requires Go and SQLite)

```sh
go run ./cmd/server
```

#### 4. gRPC API

- The service listens on port **50051** by default.
- See proto files in `api/v1/` for gRPC service definitions.
- Use tools like [grpcurl](https://github.com/fullstorydev/grpcurl) or postman with gRPC support to test endpoints.

---

## Project Structure

```
money-go/
├── api/v1/          # Protobuf definitions for User & Transaction services
├── cmd/server       # Main entrypoint for the gRPC server
├── internal/
│   ├── database/    # Database connection and setup
│   ├── models/      # Data models for users and transactions
│   ├── repository/  # Data access layer (CRUD for users & transactions)
│   └── server/      # Business logic & gRPC service implementations
├── Dockerfile
├── justfile         # Justfile for common dev tasks (build, run, proto)
└── db.sqlite        # SQLite database (auto-created)
```

---

## Example gRPC Methods

- `SendMoney(from_user_id, to_user_id, amount, currency)`
- `GetUser(user_id)`
- `GetTransaction(transaction_id)`
- `GetTransactionsByUser(user_id)`

---

## Justfile Tasks

For local development convenience:

- `just server` — Run the gRPC server
- `just build`  — Build the server binary to `bin/`
- `just proto`  — Compile proto files

---

## Known Limitations & TODOs

- No authentication or authorization (demo only)
- No advanced error handling or input validation
- No web or REST API (gRPC only)
- Minimal test coverage
- Not yet scalable or production-ready

---
