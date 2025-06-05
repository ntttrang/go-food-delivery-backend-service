# Go Food Delivery Backend Service

A comprehensive food delivery backend service built with Go, following hexagonal architecture principles and designed for microservices deployment.

Main features

- **User Service**: Authentication, registration, profile management
- **Restaurant Service**: Restaurant listings, menu management
- **Food Service**: Food items, categories, inventory
- **Cart Service**: Shopping cart operations
- **Order Service**: Order processing and *tracking (TBD)*
- **Media Service**: File upload and media management
- **Payment Service**: *Payment processing and verification (TBD)*
- **Real-time delivery tracking** *(TBD)*

## 🏗️ Project Structure

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

[https://github.com/ntttrang/go-food-delivery-backend-service/docs/food_delivery_api.json]

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Trang Nguyen** - *Initial work* - [@ntttrang](https://github.com/ntttrang)

## 🙏 Acknowledgments

- Built with [Gin](https://gin-gonic.com/) web framework
- Database ORM powered by [GORM](https://gorm.io/)
- Search functionality by [Elasticsearch](https://www.elastic.co/)
- Object storage with [MinIO](https://min.io/)
