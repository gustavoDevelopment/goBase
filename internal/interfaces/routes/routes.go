package routes

import (
	uD "api-ptf-core-business-orchestrator-go-ms/internal/interfaces/routes/domain"
	uR "api-ptf-core-business-orchestrator-go-ms/internal/interfaces/routes/utils"
	"api-ptf-core-business-orchestrator-go-ms/internal/models"

	"github.com/gorilla/mux"
)

// SetupRoutes configura todas las rutas
func SetupRoutes(router *mux.Router, a *models.Application) {
	uR.RegisterInfoRoutes(router, a)
	uR.RegisterRysncRoutes(router)
	uD.RegisterUserRoutes(router, a)
}
