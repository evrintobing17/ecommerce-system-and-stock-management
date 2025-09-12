# Ecommerce System
## Architecture

### Services
- **User Service**: Handles user authentication and profile management  
- **Product Service**: Manages product catalog and inventory  
- **Order Service**: Processes orders, payments, and stock reservations  
- **Shop Service**: Manages shops and their relationships with owners  
- **Warehouse Service**: Handles stock management and transfers between warehouses  

### Technology Stack
- **Language**: Go 1.19+  
- **Web Framework**: Gin  
- **RPC Framework**: gRPC  
- **Database**: PostgreSQL with GORM  
- **Containerization**: Docker and Docker Compose  
- **Authentication**: JWT tokens  

### Features
- User authentication with email/phone and JWT  
- Product catalog with stock availability  
- Order processing with stock reservation  
- Payment processing simulation  
- Automatic stock release for expired orders  
- Warehouse management with stock transfers  
- Shop management with multiple warehouses  
- Concurrency control for stock management  
- Comprehensive logging and error handling  

---

## Getting Started

### Prerequisites
- Go 1.19+  
- Docker and Docker Compose  
- PostgreSQL (optional, Docker version included)  

### Local Development
```bash
git clone <repository-url>
cd ecommerce-system
go mod download
docker-compose up -d postgres

go run cmd/user-service/main.go
go run cmd/product-service/main.go
go run cmd/order-service/main.go
go run cmd/shop-service/main.go
go run cmd/warehouse-service/main.go
```

# Build and start all services
docker-compose up -d --build

# Stop services
docker-compose down
