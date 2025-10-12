# Go CRUD - My First RESTful API with Golang

A clean and modern RESTful API built with Go, featuring a view-based architecture with structured input/output schemas and complete CRUD operations for blog posts.

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. All user management operations require authentication except user creation and retrieval.

### Authentication Flow

1. **Create user** via `POST /users` (public)
2. **Login** via `POST /auth/login` with email/password
3. **Receive JWT token** in response
4. **Use token** in `Authorization: Bearer <token>` header for protected endpoints
5. JWT tokens expire after **24 hours**

### Error Codes
- `401 Unauthorized`: Missing/invalid/expired token
- `403 Forbidden`: Attempting to modify another user's account

## 📋 API Endpoints

### 🔓 Public Endpoints

| Method | Endpoint | Description | Request Body | Response |
|--------|----------|-------------|--------------|----------|
| GET | `/health` | Health check | - | `{"status": "healthy", "service": "go-crud-api"}` |
| POST | `/users` | Create user account | `CreateUserInput` | `UserResponse` |
| GET | `/users/:id` | Get user details | - | `UserResponse` |
| POST | `/auth/login` | Authenticate user | `LoginInput` | `AuthResponse` |

### 🔒 Protected Endpoints (Require JWT)

| Method | Endpoint | Description | Request Body | Response |
|--------|----------|-------------|--------------|----------|
| PATCH | `/users/:id` | Update your account | `PartialUpdateUserInput` | `UserResponse` |
| DELETE | `/users/:id` | Delete your account | - | `MessageResponse` |
| POST | `/posts` | Create a new post | `CreatePostRequest` | `PostResponse` |
| GET | `/posts?page=1&limit=10` | Get posts with pagination | Query params | `ListPostsResponse` |
| GET | `/posts/:id` | Get post by ID | - | `PostResponse` |
| PUT | `/posts/:id` | Update entire post | `UpdatePostRequest` | `PostResponse` |
| PATCH | `/posts/:id` | Partial update post | `PatchPostRequest` | `PostResponse` |
| DELETE | `/posts/:id` | Delete post | - | `MessageResponse` |

## 🏗️ Project Structure

```
go-crud/ 
├── services/             # Business logic layer
│   └── post_service.go   
├── views/                # API endpoint handlers
│   └── post_views.go     # Post-specific CRUD endpoints
├── schemas/              # Input/output schemas
│   └── post_schemas.go   # Post request/response schemas
├── models/              # Data models
│   └── postModel.go     
├── initializers/        # Application initialization
│   ├── initPostgres.go  # Database connection
│   └── loadEnvVariables.go # Environment loader
├── migration/           # Database migration
│   └── migration.go     
├── examples/            # Usage examples
│   └── user_example.go  
├── main.go             # Application entry point
├── go.mod              # Go module file
└── go.sum              # Dependency checksums
```

### Architecture Layers

1. **Models** (`models/`): Define data structures and database schemas
2. **Schemas** (`schemas/`): Input/output data transfer objects with validation
3. **Services** (`services/`): Contain business logic and data validation
4. **Views** (`views/`): Handle HTTP requests/responses for specific models
5. **Initializers** (`initializers/`): Handle app startup and configuration

## 🛠️ Technologies Used

- **Go 1.24.5** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library for Go
- **PostgreSQL** - Database
- **godotenv** - Environment variable management
- **go-playground/validator** - Request validation

## ⚡ Getting Started

### Prerequisites

- Go 1.24+ installed
- PostgreSQL database running
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd go-crud
   ```

2. **Set up environment variables**
   ```bash
   # Create .env file in root directory
   touch .env
   ```
   
   Add your required configuration:
   ```env
   DB_DSN=postgres://username:password@localhost:5432/database_name?sslmode=disable
   JWT_SECRET=your-secret-key-here-change-in-production
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Run database migration**
   ```bash
   go run migration/migration.go
   ```

5. **Start the server**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

## Access Swagger Docs (After start server)

   `http://localhost:8080/swagger/index.html`

<img src="docs/swagger_docs.png" width="500" />

## 🧪 Running Tests

This project includes unit tests for the   API endpoints using a test suite pattern.

### Test Setup

Tests are located in the `test/` directory and use:
- **TestSuite pattern** for setup and cleanup
- **Real database** for integration testing
- **Structured test organization** similar to Django TestCase

### Running Tests

#### Standard Go Testing
```bash
# Run all tests with verbose output
go test -v ./test/

# Run specific test
go test -v ./test/ -run TestCreatePost_Success

# Run tests from test directory
cd test && go test -v
```

#### Enhanced Testing with gotestsum
For better test output formatting and reporting, use `gotestsum`:

```bash
# Install gotestsum (if not already installed)
go install gotest.tools/gotestsum@latest

# Run tests with enhanced formatting
gotestsum --format=short-verbose ./test -- -count=1 -v

# Run specific test with gotestsum
gotestsum --format=short-verbose ./test -- -count=1 -v -run TestCreatePost_Success

# Other useful gotestsum formats
gotestsum --format=pkgname ./test        # Package-level summary
gotestsum --format=testname ./test       # Show test names only
gotestsum --format=dots ./test           # Minimal dot output
```

### Test Structure

```
test/
├── auth_views_test.go     # Authentication endpoint tests
├── users_views_test.go    # User API endpoint tests
└── posts_views_test.go    # Post API endpoint tests
```

The test suite automatically:
- Sets up the router using the same configuration as production
- Initializes database connection and cleans test data
- Runs tests against real HTTP endpoints
- Cleans up after each test

## 📖 How to Generate go.mod and go.sum

### Creating go.mod

The `go.mod` file is the heart of Go modules. Here's how it's generated:

1. **Initialize a new module**:
   ```bash
   go mod init go-crud
   ```
   This creates a `go.mod` file with the module name.

2. **Add dependencies**:
   When you import packages in your Go files, use:
   ```bash
   go mod tidy
   ```
   This automatically adds required dependencies to `go.mod`.

3. **Manual dependency addition**:
   ```bash
   go get github.com/gin-gonic/gin
   go get gorm.io/gorm
   go get gorm.io/driver/postgres
   ```

### Understanding go.sum

The `go.sum` file contains cryptographic checksums of module dependencies:

- **Automatically generated** when you run `go mod tidy` or `go build`
- **Ensures integrity** of downloaded modules
- **Version verification** - prevents tampering
- **Should be committed** to version control

**Key commands**:
```bash
go mod tidy        # Add missing and remove unused modules
go mod verify      # Verify dependencies match go.sum
go mod download    # Download modules to local cache
```


## 🎯 Key Features Explained

### View-Based Architecture
The project uses a clean view-based pattern for organized API endpoints:

- **Model-Specific Views**: Dedicated view files for each model (e.g., `post_views.go`)
- **Standardized Responses**: Consistent API response format across all endpoints
- **Clean Structure**: Easy to maintain and extend with new models
- **Route Registration**: Automatic route setup for CRUD operations

### Schema-Driven Design
Structured input/output handling with validation:

- **Input Schemas**: `CreatePostRequest`, `UpdatePostRequest`, `PatchPostRequest`
- **Output Schemas**: `PostResponse`, `ListPostsResponse`, `ErrorResponse`
- **Data Transformation**: `ToModel()` methods convert requests to models
- **Validation Ready**: Built-in validation tags using go-playground/validator

### Service Layer
Clean separation of business logic:

- **Business Logic**: Input validation and business rules
- **Error Handling**: Proper error messages and HTTP status codes
- **Database Operations**: Abstracted database interactions
- **Reusability**: Services can be used across different views

### Database Integration
- **GORM ORM**: Powerful and developer-friendly ORM
- **Auto-migration**: Automatic database schema updates
- **Connection Pooling**: Efficient database connection management
- **Environment Configuration**: Database settings from environment variables

### Pagination Support
- **Query Parameters**: `page` and `limit` parameters for list endpoints
- **Default Values**: Automatic fallback to page=1, limit=10
- **Total Count**: Returns total records for proper pagination UI

## 🔮 Future Enhancements

- [x] ~~Add authentication and authorization~~ ✅ **Completed**
- [x] ~~Implement pagination for list endpoints~~ ✅ **Completed**
- [x] ~~Include API documentation with Swagger~~ ✅ **Completed**
- [x] ~~Add unit and integration tests~~ ✅ **Completed**
- [ ] Implement logging middleware
- [ ] Add rate limiting
- [ ] Docker containerization

## 📝 License

This project is for learning purposes.

---
*Built with ❤️ using Go and modern software architecture principles*
