# Go Clean Architecture Example

This project is a starting point for a Go application using a layered architecture based on Clean Architecture principles.

## Project Structure

- `cmd/`: Entry points of the application. For example, `cmd/server/main.go` for the main application binary.

- `internal/`: Contains the core application code, separated into layers.

  - `internal/domain/`: Contains the core domain entities and business rules. These are the most important objects in the application and have no dependencies on other layers.

  - `internal/application/`: Contains the application logic (use cases or services). It orchestrates the flow of data between the domain and the infrastructure layers.

  - `internal/infrastructure/`: Contains the implementation details, such as database access, external API clients, etc. This layer depends on the application and domain layers.

- `configs/`: Configuration files for the application (e.g., `config.yaml`).

- `go.mod`, `go.sum`: Go module files for dependency management.

## How to Run

Once the environment issue is resolved, you can run the application with:

```bash
go run cmd/server/main.go
```
