package handlers

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/models"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/utils"
	"net/http"
)

type info struct {
	AppName     string `json:"appName"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Uuid        string `json:"uuid"`
}

// Health handles the health check endpoint
func Health(w http.ResponseWriter, r *http.Request, app *models.Application) {
	appInformation := info{
		AppName:     app.Configs().AppName,
		Version:     app.Configs().Version,
		Environment: app.Configs().Environment,
		Uuid:        app.Configs().Uuid,
	}
	utils.SendSuccess(w, "SUCCESS", "Service is healthy", http.StatusOK, appInformation)
}

// Health handles the health check endpoint
func Rysnc(w http.ResponseWriter, r *http.Request) {
	utils.SendSuccess(w, "SUCCESS", "Rsync started", http.StatusOK, nil)
}
