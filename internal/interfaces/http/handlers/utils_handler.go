package handlers

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/utils"
	"net/http"
)

// Health handles the health check endpoint
func Health(w http.ResponseWriter, r *http.Request) {
	utils.SendSuccess(w, "SUCCESS", "Service is healthy", http.StatusOK, nil)
}

// Health handles the health check endpoint
func Rysnc(w http.ResponseWriter, r *http.Request) {
	utils.SendSuccess(w, "SUCCESS", "Rsync started", http.StatusOK, nil)
}
