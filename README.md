# Heli Task - Todo Service (Golang)

A simple **Todo Service** built with **Golang**, **PostgreSQL**, and **AWS SQS (via LocalStack)** following **Clean Architecture principles**.

---

## 🚀 Features

- Create and persist `TodoItem` in **PostgreSQL**
- Publish new Todo messages to **AWS SQS** (simulated with LocalStack)
- Clean architecture with separation of **Domain**, **Repository**, **Service**, and **Adapters**
- **Unit and Integration tests**
- **Mocked SQS** in unit tests to avoid external dependencies

---

## 🛠 Prerequisites

- **Go** >= 1.23.4
- **Docker & Docker Compose**
- **AWS CLI** (for LocalStack interaction)
- **Make** (optional but recommended for simplified commands)

---

## ⚡ Project Setup

### 0️⃣ Prepare dependencies (Important)

After cloning the repository, run the following commands inside the project directory to download and prepare all dependencies:

```bash
make sync-vendor
```
### 1️⃣ Start the project

```bash
make run
```

- Spins up Docker containers (PostgreSQL + LocalStack)
- Create migrations
- Runs the Go service on `http://localhost:8080`

---

### 2️⃣ Stop all containers

```bash
make stop
```

Stops all Docker containers and cleans up the environment.

---

### 3️⃣ API Endpoints

| Method | Endpoint         | Description               |
|--------|-----------------|---------------------------|
| GET    | `/health/check`  | Health check endpoint      |
| GET    | `/todo`          | Retrieve all todos         |
| POST   | `/todo`          | Create a new todo item     |

**Sample Request (POST /todo)**

```bash
curl -X POST http://localhost:8080/todo \
  -H "Content-Type: application/json" \
  -d '{
        "description": "Finish Golang test task",
        "dueDate": "2025-08-05"
      }'
```

**Sample Response**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "description": "Finish Golang test task",
  "dueDate": "2025-08-05T00:00:00Z",
  "createdAt": "2025-08-04T10:21:00Z"
}
```

---

## 🧪 Running Tests

### Run all tests

```bash
make test
```

This runs:

- ✅ **Unit tests** for `TodoService` and `SQSAdapter` (with mocked SQS)
- ✅ **Integration tests** for PostgreSQL repository

### Run a specific test

```bash
make test-specific name={{name}}
```

---

## 📬 Working with SQS (LocalStack)

1️⃣ **List all queues**

```bash
make list-queues
```

---

2️⃣ **Receive messages from the queue**

```bash
make receive-messages
```

This shows all tasks published to the queue without deleting them.

---

## ✅ Notes

- Integration tests run on **real PostgreSQL** and clean up test data automatically.
- Unit tests mock SQS to avoid dependency on LocalStack.
- Use `make receive-messages` to inspect messages while testing locally.
- Use `make list-queues` to verify queues exist before receiving messages.
- Only valid todos with non-empty description and future due dates are persisted and published to SQS.

---

**Built with ❤️ in Golang**
