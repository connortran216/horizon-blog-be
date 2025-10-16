# Go CRUD API

A learning project demonstrating modern Go development practices through a production-ready RESTful API with JWT authentication, PostgreSQL, and clean architecture. Perfect for learning Golang patterns and best practices.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose

### Run with Docker (Recommended)
```bash
git clone https://github.com/connortran216/go-crud.git
cd go-crud
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Manual Setup
```bash
go mod tidy
go run main.go
```

## 🔐 Authentication

JWT-based authentication with 24-hour token expiration.

**Public endpoints:**
- `POST /users` - Create account
- `POST /auth/login` - Authenticate
- `GET /users/:id` - Get user profile

**Protected endpoints** (require `Authorization: Bearer <token>`):
- All other user and post operations

## 📋 API Endpoints

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users` | Create user account |
| GET | `/users/:id` | Get user details |
| PATCH | `/users/:id` | Update account |
| DELETE | `/users/:id` | Delete account |

### Posts
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/posts` | Create post |
| GET | `/posts` | List posts (paginated) |
| GET | `/posts/:id` | Get post |
| PUT | `/posts/:id` | Update post |
| PATCH | `/posts/:id` | Partial update |
| DELETE | `/posts/:id` | Delete post |

### Tags
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/tags` | Create tag |
| GET | `/tags` | List all tags |
| GET | `/tags/popular` | Popular tags |
| GET | `/tags/search` | Search tags |
| GET | `/tags/:id` | Get tag |
| PUT | `/tags/:id` | Update tag |
| DELETE | `/tags/:id` | Delete tag |

### System
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |

## 🏗️ Architecture

Clean architecture with SOLID principles:

- **Models** - GORM entities with relationships
- **Schemas** - Request/response DTOs with validation
- **Services** - Business logic layer
- **Views** - HTTP handlers organized by model
- **Middleware** - Authentication, logging, rate limiting

## 🛠️ Tech Stack

- **Go 1.24.5** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM with PostgreSQL
- **JWT** - Authentication
- **Docker** - Containerization

## 📖 Learning Focus

This project demonstrates key Go concepts and patterns:

### 🏛️ Clean Architecture
- **SOLID principles** in Go implementation
- **Separation of concerns** across layers
- **Dependency injection** patterns
- **Interface-based design** for testability

### 🛠️ Go Best Practices
- **Error handling** strategies and patterns
- **Context usage** for request lifecycle management
- **Middleware patterns** for cross-cutting concerns
- **Environment-based configuration**

### 🧪 Testing Strategies
- **Table-driven tests** for comprehensive coverage
- **Integration testing** with real database
- **Test factories** for consistent test data
- **Test suite patterns** for setup/teardown

### 🔒 Production-Ready Features
- **JWT authentication** with proper security practices
- **Rate limiting** for API protection
- **Structured logging** for observability
- **Docker containerization** for deployment

### 📊 Database Patterns
- **GORM modeling** with relationships and constraints
- **Migration strategies** for schema evolution
- **Connection pooling** and performance optimization
- **Transaction management** where appropriate

## 🧪 Testing

```bash
# Run all tests
go test -v ./test/

# With enhanced output
gotestsum --format=short-verbose ./test/
```

## 📚 Documentation

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **API Examples**: See test files in `/test/`

## 🔒 Security Features

- **Rate Limiting**: 100ms intervals, burst capacity of 5
- **Input Validation**: Multi-layer validation with go-playground/validator
- **Error Handling**: Structured error responses
- **Logging**: Selective logging for errors and performance

## 🔮 Roadmap

- [x] Add authentication and authorization ✅
- [x] Implement pagination for list endpoints ✅
- [x] Include API documentation with Swagger ✅
- [x] Add unit and integration tests ✅
- [x] Implement logging middleware ✅
- [x] Add rate limiting ✅
- [x] Docker containerization ✅
- [x] Category/Tags System with Many-to-Many Relationships ✅
- [ ] File Upload System (Images/Media)
- [ ] Search and Filtering Engine
- [ ] Comments and Reactions System
- [ ] Caching Layer with Redis Integration

---
*Built with Go and modern software architecture principles*
