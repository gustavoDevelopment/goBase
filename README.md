# 🚀 Go API Template with Clean Architecture

This is a Go project template that follows the principles of Clean Architecture, designed to be the starting point for scalable and maintainable applications.

## 🌟 Features

- 🏗️ **Clean Architecture** with clear separation of layers.
- 🗄️ **Multiple Database Support**:
    - **MongoDB**: Main database with a generic repository pattern.
    - **PostgreSQL**: Ready to be used.
    - **Oracle**: Ready to be used.
- 📝 **Detailed Documentation**
- 🔄 **GitHub Actions** for CI/CD
- 🐳 **Docker** & Docker Compose ready for production

## 🚦 Prerequisites

- Go 1.25 or higher
- Docker
- Docker Compose
- Access to MongoDB, PostgreSQL, and Oracle databases.

## ⚡ Quick Start

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/api-template-go-ms.git
    cd api-template-go-ms
    ```

2.  **Initialize the project:**
    ```bash
    chmod +x init.sh
    ./init.sh
    ```

3.  **Set up the environment:**
    Create a `.env` file in the root of the project with the following variables:
    ```env
    # HTTP Server
    BASE_PATH=/api/template-go-ms/v1
    
    # MongoDB
    MONGO_URI=mongodb://user:password@host:port
    MONGO_DATABASE=your-db
    TIMEOUT=10s
    
    # PostgreSQL
    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_USER=user
    POSTGRES_PASSWORD=password
    POSTGRES_DBNAME=your-db
    POSTGRES_SSLMODE=disable
    POSTGRES_MAXOPENCONNS=10
    POSTGRES_MAXIDLECONNS=5
    
    # Oracle
    ORACLE_HOST=localhost
    ORACLE_PORT=1521
    ORACLE_USER=user
    ORACLE_PASSWORD=password
    ORACLE_DBNAME=your-db
    ORACLE_MAXOPENCONNS=10
    ORACLE_MAXIDLECONNS=5
    
    # JWT
    JWT_SECRET=your-secret-key
    
    # JSON Config
    JSON_CONFIG_PATH=./configs/parameters.json
    ```

4.  **Run the application:**
    ```bash
    go run ./cmd/main.go
    ```

## 🏗️ Project Structure

```
.
├── cmd/
│   ├── app/              # Application startup and shutdown logic
│   └── main.go           # Entry point
├── configs/              # Configuration files (config.yaml)
├── internal/
│   ├── application/      # Business logic
│   ├── client/           # External API clients
│   ├── config/           # Configuration loading
│   ├── domain/           # Domain models and interfaces
│   ├── infrastructure/   # Concrete implementations (database, repositories)
│   │   ├── database/
│   │   │   ├── mongo/
│   │   │   ├── oracle/
│   │   │   └── postgrest/
│   │   └── repository/
│   ├── interfaces/       # Controllers and routes
│   └── pkg/              # Shared packages (logger, utils)
├── go.mod
├── go.sum
├── init.sh
└── README.md
```

## ⚙️ Configuration

The application is configured through `configs/config.yaml` and environment variables. Environment variables take precedence.

### Databases

The application can connect to MongoDB, PostgreSQL, and Oracle. The connections are initialized at startup based on the configuration provided.

- **MongoDB**: Used as the main database for the generic repository.
- **PostgreSQL**: Connection available.
- **Oracle**: Connection available.

## 🧪 Testing

To run the tests, use the following command:
```bash
go test -v ./...
```

## 🐳 Docker

To build and run the application with Docker, use the following commands:

1.  **Build the image:**
    ```bash
    docker build -t api-template-go-ms .
    ```

2.  **Run the container:**
    ```bash
    docker run -p 8426:8426 --env-file .env api-template-go-ms
    ```

## 🚀 Usage Tips

This section provides tips on how to use the different components of the application, including database connections.

### Accessing Database Connections

All database connections are initialized at startup and are available through the `models.Application` struct, which is passed to your handlers and services.

-   **MongoDB**: The main database, used with a generic repository pattern.
    ```go
    // In your service or handler
    mongoDBClient := app.MongoDB() // Returns *mongoDb.Database
    // You can now use mongoDBClient to interact with MongoDB
    // For example, using a repository:
    userRepo := repository.NewMongoUserRepository(mongoDBClient)
    ```

-   **PostgreSQL**: The PostgreSQL connection pool is available for use.
    ```go
    // In your service or handler
    pgPool := app.PostgreSQLPool() // Returns *pgxpool.Pool
    // You can now use pgPool to execute queries against PostgreSQL
    rows, err := pgPool.Query(context.Background(), "SELECT name FROM users WHERE id=$1", 1)
    ```

-   **Oracle**: The Oracle database connection is also available.
    ```go
    // In your service or handler
    oracleDB := app.OraclePool() // Returns *sql.DB
    // You can now use oracleDB to execute queries against Oracle
    rows, err := oracleDB.QueryContext(context.Background(), "SELECT product_name FROM products WHERE id = :1", 101)
    ```

## 📝 How to Add a New Handler

Follow these steps to add a new API endpoint to the application. We'll use a "Product" entity as an example.

### 1. Define the Domain Model

Create your model in the `internal/domain/` directory (e.g., `product.go`).

```go
// internal/domain/product.go
package domain

import "time"

type Product struct {
    ID          string    `json:"id" bson:"_id,omitempty"`
    Name        string    `json:"name" bson:"name"`
    Price       float64   `json:"price" bson:"price"`
    DateCreated time.Time `json:"date_created" bson:"date_created"`
}
```

### 2. Create the Application Service

Define the business logic in the `internal/application/` directory (e.g., `product_service.go`).

```go
// internal/application/product_service.go
package application

import (
	"api-template-go-ms/internal/domain"
	"api-template-go-ms/internal/infrastructure/repository"
	"context"
)

type ProductService struct {
	productRepo *repository.GenericRepository[domain.Product]
}

func NewProductService(repo *repository.GenericRepository[domain.Product]) *ProductService {
	return &ProductService{productRepo: repo}
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}
```

### 3. Create the HTTP Handler

Create the handler that will process HTTP requests in `internal/interfaces/http/handlers/` (e.g., `product_handler.go`).

```go
// internal/interfaces/http/handlers/product_handler.go
package handlers

import (
	"api-template-go-ms/internal/application"
	"net/http"

	"github.com/gorilla/mux"
)

func GetProduct(service *application.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		product, err := service.GetProductByID(r.Context(), id)
		if err != nil {
			// Handle error (e.g., write a 404 or 500 response)
			return
		}

		// Write success response with the product
	}
}
```

### 4. Register the Route

Create a new file in `internal/interfaces/routes/domain/` to register the routes for your new domain (e.g., `product_routes.go`).

```go
// internal/interfaces/routes/domain/product_routes.go
package domain

import (
	"api-template-go-ms/internal/application"
	"api-template-go-ms/internal/domain"
	"api-template-go-ms/internal/infrastructure/repository"
	"api-template-go-ms/internal/interfaces/http/handlers"
	"api-template-go-ms/internal/models"

	"github.com/gorilla/mux"
)

func RegisterProductRoutes(router *mux.Router, app *models.Application) {
	// Initialize repository and service
	productRepo := repository.NewGenericRepository[domain.Product](app.MongoDB(), "products")
	productService := application.NewProductService(productRepo)

	// Define routes
	router.HandleFunc("/products/{id}", handlers.GetProduct(productService)).Methods("GET")
}
```

### 5. Add to Main Router Setup

Finally, call your new route registration function from `internal/interfaces/routes/routes.go`.

```go
// internal/interfaces/routes/routes.go
package routes

import (
	// ... other imports
)

func SetupRoutes(router *mux.Router, a *models.Application) {
    uR.RegisterInfoRoutes(router, a)
    uR.RegisterRysncRoutes(router)
    uD.RegisterUserRoutes(router, a)
    uD.RegisterProductRoutes(router, a) // Add this line
    eX.RegisterExampleRoutes(router)
}
```

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1.  Fork the project.
2.  Create a new branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.
