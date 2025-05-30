# Go Food Delivery Backend Service

A comprehensive food delivery backend service built with Go, following hexagonal architecture principles and designed for microservices deployment.

## 🚀 Overview

This backend service powers a modern food delivery application with features including user authentication, restaurant management, order processing, payment integration (TBD), and real-time delivery tracking (TBD). The system is designed with a modular architecture that supports scalability and maintainability.

## 🏗️ Architecture

The project follows **Hexagonal Architecture** (Ports and Adapters) pattern with clear separation of concerns:

- **Domain Layer**: Core business logic and entities
- **Application Layer**: Use cases and business workflows
- **Infrastructure Layer**: External integrations (database, APIs, etc.)
- **Interface Layer**: HTTP controllers and API endpoints

### Microservices Structure

The application is organized into the following modules:

- **User Service**: Authentication, registration, profile management
- **Restaurant Service**: Restaurant listings, menu management
- **Food Service**: Food items, categories, inventory
- **Cart Service**: Shopping cart operations
- **Order Service**: Order processing and tracking
- **Payment Service**: Payment processing and verification
- **Media Service**: File upload and media management

## 🛠️ Tech Stack

- **Language**: Go 1.24.1
- **Web Framework**: Gin
- **Database**: MySQL with GORM
- **Cache**: Redis
- **Search**: Elasticsearch
- **Object Storage**: MinIO (S3-compatible)
- **Authentication**: JWT, OAuth2 (Google)
- **Email**: SMTP with Gomail
- **Containerization**: Docker
- **Orchestration**: Kubernetes (planned)

## 📋 Features

### Core Functionality
- ✅ User registration and authentication (Email, Google OAuth)
- ✅ Restaurant management and listings
- ✅ Food/menu item management with categories
- ✅ Shopping cart operations (add, update, delete, list)
- ✅ Order processing and management
- ✅ Payment card management (create, list, update status)
- ✅ Media upload and management
- ✅ Search functionality with Elasticsearch (foods & restaurants)
- ✅ Review and rating system (foods & restaurants)
- ✅ Favorites system (foods & restaurants)
- ✅ User address management
- ✅ Email verification with Redis-based code generation

### Technical Features
- 🔐 JWT-based authentication with token introspection
- 📧 Email verification system with SMTP
- 🔍 Advanced Elasticsearch search with facets and filtering
- 📱 RESTful API design with proper error handling
- 🚀 Redis caching for verification codes
- 📊 Structured logging with GORM
- 🐳 Docker containerization with multi-stage builds
- 🔄 Elasticsearch index synchronization
- 📈 Health check endpoints (/ping)
- 🏗️ RPC communication between modules
- 🔒 Middleware-based authentication and recovery

## 🚦 Getting Started

### Prerequisites

- Go 1.24.1 or higher
- Docker and Docker Compose
- MySQL 8.0+
- Redis 7.0+
- MinIO
- Elasticsearch 8.12+

### Environment Variables

Create a `.env` file in the root directory:

```env
# Database
DB_DSN=user:password@tcp(localhost:3306)/food_delivery?charset=utf8mb4&parseTime=True&loc=Local

# Server
PORT=3000
GIN_MODE=release

# JWT
JWT_SECRET_KEY=your-jwt-secret-key

# OAuth2 Google
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:3000/v1/google/callback

# Email SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# MinIO S3
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false

# Elasticsearch
ELASTICSEARCH_URL=http://localhost:9200
ELASTICSEARCH_USERNAME=
ELASTICSEARCH_PASSWORD=

# Service URLs (for RPC communication)
USER_SERVICE_URL=http://localhost:3000/v1
FOOD_SERVICE_URL=http://localhost:3000/v1
RESTAURANT_SERVICE_URL=http://localhost:3000/v1
CAT_SERVICE_URL=http://localhost:3000/v1
```

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/ntttrang/go-food-delivery-backend-service.git
   cd go-food-delivery-backend-service
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start infrastructure services**
   ```bash
   cd docs
   docker-compose up -d
   ```

4. **Run database migrations**
   ```bash
   # Import the SQL schema
   mysql -u root -p food_delivery < docs/food_delivery.sql
   ```

5. **Start the application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:3000`

### Docker Deployment

1. **Build the Docker image**
   ```bash
   docker build -t food-delivery-backend .
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose -f docs/docker-compose.yml up -d
   ```

## 📚 API Documentation

### Health Check

```http
GET /ping
```

### API Endpoints

All API endpoints are prefixed with `/v1`:

#### User Management

- `POST /v1/register` - User registration
- `POST /v1/authenticate` - User login
- `POST /v1/google/signup` - Google OAuth signup
- `GET /v1/google/callback` - Google OAuth callback
- `GET /v1/profile` - Get user profile (🔒 Auth required)
- `GET /v1/generate-code` - Generate verification code (🔒 Auth required)
- `GET /v1/verify/{code}` - Verify email code (🔒 Auth required)
- `POST /v1/users` - Create user (🔒 Auth required)
- `GET /v1/users` - List users
- `GET /v1/users/{id}` - Get user details
- `PATCH /v1/users/{id}` - Update user (🔒 Auth required)
- `POST /v1/users/address` - Create user address (🔒 Auth required)
- `GET /v1/users/address` - List user addresses

#### Categories

- `POST /v1/categories` - Create category (🔒 Auth required)
- `GET /v1/categories` - List categories
- `GET /v1/categories/{id}` - Get category details
- `PATCH /v1/categories/{id}` - Update category (🔒 Auth required)
- `DELETE /v1/categories/{id}` - Delete category (🔒 Auth required)

#### Restaurant Management

- `POST /v1/restaurants` - Create restaurant (🔒 Auth required)
- `GET /v1/restaurants` - List restaurants
- `GET /v1/restaurants/{id}` - Get restaurant details
- `PATCH /v1/restaurants/{id}` - Update restaurant
- `DELETE /v1/restaurants/{id}` - Delete restaurant
- `POST /v1/restaurants/favorites` - Add/remove favorite restaurant (🔒 Auth required)
- `GET /v1/restaurants/favorites` - List favorite restaurants (🔒 Auth required)
- `POST /v1/restaurants/comments` - Create restaurant comment (🔒 Auth required)
- `GET /v1/restaurants/comments` - List restaurant comments
- `DELETE /v1/restaurants/comments/{id}` - Delete restaurant comment
- `POST /v1/restaurants/menu-item` - Create menu item
- `GET /v1/restaurants/menu-item/{restaurantId}` - List menu items
- `DELETE /v1/restaurants/menu-item` - Delete menu item
- `POST /v1/restaurants/search` - Search restaurants

#### Food & Menu

- `POST /v1/foods` - Create food item
- `GET /v1/foods` - List food items
- `GET /v1/foods/{id}` - Get food details
- `PATCH /v1/foods/{id}` - Update food item
- `DELETE /v1/foods/{id}` - Delete food item
- `POST /v1/foods/favorites` - Add/remove favorite food (🔒 Auth required)
- `GET /v1/foods/favorites` - List favorite foods (🔒 Auth required)
- `POST /v1/foods/comments` - Create food comment (🔒 Auth required)
- `GET /v1/foods/comments` - List food comments
- `DELETE /v1/foods/comments/{id}` - Delete food comment

#### Cart Operations

- `POST /v1/carts` - Add/update cart item (🔒 Auth required)
- `GET /v1/carts` - List cart items
- `GET /v1/carts/cart-item` - List detailed cart items
- `GET /v1/carts/{userId}/{foodId}` - Get specific cart item
- `PATCH /v1/carts/{id}` - Update cart item
- `DELETE /v1/carts/{userId}/{foodId}` - Remove cart item

#### Order Management

- `POST /v1/orders` - Create order (🔒 Auth required)
- `GET /v1/orders` - List orders
- `GET /v1/orders/{id}` - Get order details
- `PATCH /v1/orders/{id}` - Update order
- `DELETE /v1/orders/{id}` - Delete order

#### Payment Management

- `POST /v1/cards` - Create payment card (🔒 Auth required)
- `GET /v1/cards/{id}` - Get card details
- `PATCH /v1/cards/{id}` - Update card status
- `GET /v1/cards/user/{userId}` - Get user's cards

#### Media Management

- `PUT /v1/medias` - Upload media file (🔒 Auth required)

## 🧪 Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## 📁 Project Structure

```text
├── main.go                 # Application entry point
├── middleware/             # HTTP middleware (auth, recovery)
├── modules/               # Business modules (hexagonal architecture)
│   ├── user/             # User management & authentication
│   │   ├── infras/       # Infrastructure layer
│   │   │   ├── controller/http-gin/  # HTTP controllers
│   │   │   └── repository/gorm-mysql/ # Data repositories
│   │   ├── model/        # Domain models
│   │   ├── service/      # Business logic
│   │   └── module.go     # Module setup
│   ├── restaurant/       # Restaurant operations
│   ├── food/            # Food items & categories
│   ├── cart/            # Shopping cart
│   ├── order/           # Order processing
│   ├── payment/         # Payment card management
│   ├── media/           # Media upload
│   └── category/        # Food categories
├── shared/               # Shared utilities
│   ├── component/       # Reusable components (JWT, Redis, Email, etc.)
│   ├── datatype/        # Common data types & errors
│   ├── infras/          # Infrastructure setup (DB, context)
│   └── model/           # Shared models & utilities
├── deployments/         # Deployment configurations
├── docs/               # Documentation & docker-compose
│   ├── docker-compose.yml
│   ├── food_delivery.sql
│   └── Note.md
└── uploads/            # File uploads directory
```

### Module Architecture

Each module follows hexagonal architecture:

- **`infras/controller/`** - HTTP handlers and route setup
- **`infras/repository/`** - Data access layer (MySQL, Elasticsearch, RPC)
- **`service/`** - Business logic and use cases
- **`model/`** - Domain entities and DTOs
- **`module.go`** - Dependency injection and module setup

## 🔍 Advanced Features

### Search Functionality

The application includes advanced search capabilities powered by Elasticsearch:

- **Restaurant Search**: Full-text search with location-based filtering
- **Food Search**: Multi-field search with faceted navigation
- **Auto-indexing**: Automatic synchronization between MySQL and Elasticsearch
- **Admin Endpoints**: Manual index management and synchronization

### RPC Communication

Modules communicate via internal RPC calls:

- **User Service**: Token introspection and user data lookup
- **Food Service**: Food information retrieval for cart operations
- **Restaurant Service**: Restaurant data for order processing

### Authentication & Security

- **JWT Tokens**: Stateless authentication with configurable expiration
- **Google OAuth2**: Social login integration
- **Email Verification**: Redis-based verification code system
- **Middleware Protection**: Route-level authentication enforcement

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Trang Nguyen** - *Initial work* - [@ntttrang](https://github.com/ntttrang)

## 🙏 Acknowledgments

- Built with [Gin](https://gin-gonic.com/) web framework
- Database ORM powered by [GORM](https://gorm.io/)
- Search functionality by [Elasticsearch](https://www.elastic.co/)
- Object storage with [MinIO](https://min.io/)
