package example

import (
	handlers "api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/handlers/example"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/constants"

	"github.com/gorilla/mux"
)

func RegisterExampleRoutes(router *mux.Router) {
	// Initialize handlers
	exampleHandler := handlers.NewExampleHandler()

	subrouter := router.PathPrefix(constants.REST_CLIENT_GROUP).Subrouter()
	subrouter.HandleFunc("/characters", exampleHandler.GetAllCharacters).Methods(constants.GET)
	subrouter.HandleFunc("/planets", exampleHandler.GetAllPlanets).Methods(constants.GET)
}
