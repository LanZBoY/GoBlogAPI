# GoBlogAPI

GoBlogAPI is a simple blog REST API built with [Gin](https://github.com/gin-gonic/gin) and MongoDB. The service exposes endpoints for user accounts, authentication and posts, and includes automatically generated Swagger documentation.

## Setup

1. Copy `.env_example` to `.env` and adjust the values as needed:
   ```
   SERVICE_NAME="Blog API"
   MONGO_URI="mongodb://WenTee:jp4wu6@localhost:27017"
   MOGNO_DATABASE="Blog"
   JWT_SECRET="xxxaaa"
   ```
2. Start MongoDB using Docker Compose:
   ```
   docker-compose up -d
   ```
   The compose file launches the `mongo:latest` image on port `27017`.
3. Run the application:
   ```
   go run app/main.go
   ```
   The server listens on `:8080` and serves Swagger docs at `/swagger/index.html`.

## Usage

With the application running you can browse the API docs at `http://localhost:8080/swagger/index.html` or interact with the endpoints using any HTTP client. Example routes include `/auth/login`, `/users`, and `/posts`.

To run tests and generate a coverage report you can use the provided make targets:
```bash
make test
make generate_test_report
```

