package utilsRoutes

import (
	"net/http"

	"api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/handlers"
	"api-ptf-core-business-orchestrator-go-ms/internal/models"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/constants"

	"github.com/gorilla/mux"
)

func RegisterInfoRoutes(router *mux.Router, a *models.Application) {
	subrouter := router.PathPrefix(constants.UTILS_GROUP).Subrouter()
	subrouter.HandleFunc(constants.HEALTH_CHECK, func(w http.ResponseWriter, r *http.Request) {
		handlers.Health(w, r, a)
	}).Methods(constants.GET)
}
