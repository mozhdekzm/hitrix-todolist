# Todo List Service - GraphQL + GORM + Golang

This project is a **Todo List Service** built with:

- [Golang](https://golang.org/)
- [GraphQL (gqlgen)](https://gqlgen.com/)
- [GORM](https://gorm.io/) as ORM for MySQL
- MySQL and Redis for persistence and streaming

---

## 📦 Features

- Create, Read, Update, and Delete (CRUD) operations for Todos
- GraphQL API with gqlgen
- MySQL integration with GORM
- Redis Stream publisher for events
- Dockerized setup with `docker-compose`

---

## 🛠️ Requirements

- Go 1.23+
- Docker & Docker Compose
- Make

---

## ⚡ Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/your-username/todo-gql.git
cd todo-gql

2. Start dependencies (MySQL & Redis)
docker-compose up -d

3. Run the project
make run

4. Sync vendors (if needed)
make sync-vendor