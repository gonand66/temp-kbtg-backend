# Temp KBTG Backend

Backend API built with Go and Fiber framework.

## 🚀 Features

- **Fast**: Built on top of Fasthttp, the fastest HTTP engine for Go
- **Express-like**: Familiar API for developers coming from Node.js/Express
- **Middleware**: CORS, Logger, Recover middleware included
- **Clean Structure**: Organized project structure for scalability

## 📋 Prerequisites

- Go 1.17 or higher
- Git

## 🛠️ Installation

1. Install dependencies:

```bash
go mod download
```

2. Copy environment file:

```bash
cp .env.example .env
```

3. Update `.env` with your configuration

## 🏃 Running the Application

### Development

```bash
go run main.go
```

### Production Build

```bash
go build -o bin/app main.go
./bin/app
```

## 📁 Project Structure

```
.
├── main.go           # Application entry point
├── config/           # Configuration files
├── handlers/         # Request handlers
├── middleware/       # Custom middleware
├── models/           # Data models
├── routes/           # Route definitions
├── utils/            # Utility functions
├── .env             # Environment variables
├── .env.example     # Example environment variables
├── .gitignore       # Git ignore file
└── README.md        # This file
```

## 🔌 API Endpoints

### Health Check

```
GET /api/v1/health
```

### Welcome

```
GET /
```

### Users

```
GET    /api/v1/users      # Get all users
GET    /api/v1/users/:id  # Get user by ID
POST   /api/v1/users      # Create new user
```

## 📝 Example Request

### Create User

```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

## 🔧 Technologies

- [Go](https://golang.org/) - Programming language
- [Fiber v2](https://gofiber.io/) - Web framework
- [Fiber Middleware](https://docs.gofiber.io/) - Built-in middleware

## 📖 Documentation

- [Fiber Documentation](https://docs.gofiber.io/)
- [Go Documentation](https://golang.org/doc/)

## 👨‍💻 Development

To add new routes:

1. Create handler in `handlers/` directory
2. Define routes in `routes/` directory
3. Register routes in `main.go`

## 📄 License

MIT
