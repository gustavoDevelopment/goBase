package main

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/application"
	"api-ptf-core-business-orchestrator-go-ms/internal/infrastructure"
	"context"
	"fmt"
)

func main() {
	fmt.Println("Starting application...")

	// 1. Create repository (infrastructure layer)
	userRepo := infrastructure.NewInMemoryUserRepository()

	// 2. Create service (application layer), injecting the repository
	userService := application.NewUserService(userRepo)

	// 3. Use the service to get a user
	user, err := userService.GetUserByID(context.Background(), 1)
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
		return
	}

	fmt.Printf("Found user: %+v\n", user)
}
