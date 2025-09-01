package utilsRoutes

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/handlers"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/constants"

	"github.com/gorilla/mux"
)

func RegisterRysncRoutes(router *mux.Router) {
	subrouter := router.PathPrefix(constants.UTILS_GROUP).Subrouter()
	subrouter.HandleFunc(constants.RSYNC, handlers.Rysnc).Methods(constants.GET)
}
