# Auth Service

A microservice authentication system built with Go, providing comprehensive user authentication and authorization functionality through gRPC APIs.

## 🚀 Features

- **User Registration & Login**: Secure user registration and authentication
- **JWT Token Management**: Access and refresh token handling
- **Password Management**: Forgot password, reset password functionality
- **Account Verification**: Email-based account verification
- **Session Management**: User session tracking and management
- **Role-based Access Control**: User roles and permissions system
- **gRPC API**: High-performance RPC communication
- **Database Integration**: PostgreSQL with migrations
- **Caching**: Redis integration for performance
- **Queue System**: Asynchronous task processing
- **Service Discovery**: Automatic service registration and discovery

## 🏗️ Architecture

```
auth-service/
├── bootstrap/          # Application bootstrap and configuration
├── cmd/               # Application entry point
├── constants/         # Application constants
├── domain/           # Business logic layer
│   ├── entity/       # Domain entities
│   ├── repository/   # Data access interfaces
│   └── usecase/      # Business use cases
├── infrastructure/   # Infrastructure layer
│   ├── grpc_client/  # gRPC client implementations
│   ├── grpc_service/ # gRPC server implementations
│   └── repo/         # Repository implementations
└── migrations/       # Database migrations
```

## 📋 Prerequisites

- Go 1.24.6 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher
- Docker (optional, for development)

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd auth-service
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**
   ```bash
   # Create database
   make dev-create-db
   
   # Run migrations
   make migrate-dev-up
   ```

4. **Configure environment**
   - Copy `dev.config.yaml` and modify as needed
   - Update database connection string
   - Set JWT secrets and other security keys

## 🚀 Running the Application

### Development Mode
```bash
# Run the application
make run

# Or directly with go
go run cmd/main.go
```

### Build and Run
```bash
# Build the application
make build-grpc

# Run the built binary
./bin/grpc-server
```

## 📊 Database Management

### Migrations
```bash
# Apply migrations
make migrate-dev-up

# Rollback migrations
make migrate-dev-down

# Reset database
make migrate-dev-reset

# Drop database
make migrate-dev-drop

# Create new migration
make migrate-dev-create name=migration_name
```

### Database Operations
```bash
# Create database
make dev-create-db

# Drop database
make dev-drop-db

# Docker operations (if using Docker)
make dev-docker-create-db
make dev-docker-drop-db
```

## 🔧 Configuration

The application uses `dev.config.yaml` for configuration. Key settings include:

- **Database**: PostgreSQL connection string
- **gRPC**: Service port and host configuration
- **Redis**: Cache and queue configuration
- **JWT**: Secret keys for different token types
- **Mail Service**: gRPC client configuration for email service

### Environment Variables

The application supports the following environment variables:

- `MODE_ENV`: Environment mode (development/production)
- `URL_DB`: Database connection string
- `NAME_SERVICE`: Service name for discovery
- `PORT_GRPC`: gRPC server port
- `HOST_GRPC`: gRPC server host
- `SECRET_OTP`: OTP secret key
- `JWT_SECRET`: JWT secret keys
- `FRONTEND_URL`: Frontend application URL
- `MAIL_SERVICE_ADDR`: Mail service address

## 🔌 API Endpoints

The service provides the following gRPC endpoints:

### Authentication
- `Register`: User registration
- `Login`: User authentication
- `Logout`: User logout
- `RefreshToken`: Refresh access token

### Password Management
- `ForgotPassword`: Initiate password reset
- `ResetPasswordByToken`: Reset password using token
- `ResetPasswordByCode`: Reset password using code

### Account Management
- `VerifyAccount`: Verify user account
- `CheckToken`: Validate token
- `CheckCode`: Validate verification code

## 🏗️ Project Structure

### Domain Layer
- **Entities**: User, Session, Role models
- **Repositories**: Data access interfaces
- **Use Cases**: Business logic implementation

### Infrastructure Layer
- **gRPC Services**: API endpoint implementations
- **Repositories**: Database implementations
- **gRPC Clients**: External service clients

### Bootstrap
- **App**: Application initialization
- **Environment**: Configuration management

## 🔒 Security Features

- **Password Hashing**: Argon2id for secure password storage
- **JWT Tokens**: Secure token-based authentication
- **Session Management**: Secure session handling
- **Input Validation**: Comprehensive request validation
- **Rate Limiting**: Protection against abuse

## 🧪 Development

### Adding New Features
1. Define domain entities in `domain/entity/`
2. Create repository interfaces in `domain/repository/`
3. Implement business logic in `domain/usecase/`
4. Add gRPC service implementation in `infrastructure/grpc_service/`
5. Update database schema with migrations

### Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## 📝 Dependencies

### Core Dependencies
- `github.com/anhvanhoa/service-core`: Core service framework
- `github.com/anhvanhoa/sf-proto`: Protocol buffer definitions
- `github.com/go-pg/pg/v10`: PostgreSQL ORM
- `go.uber.org/zap`: Structured logging
- `google.golang.org/grpc`: gRPC framework

### Key Features
- **Authentication**: JWT-based authentication system
- **Database**: PostgreSQL with go-pg ORM
- **Caching**: Redis integration
- **Queue**: Asynq for background jobs
- **Logging**: Structured logging with Zap
- **Configuration**: Viper for configuration management

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions, please contact the development team or create an issue in the repository.
