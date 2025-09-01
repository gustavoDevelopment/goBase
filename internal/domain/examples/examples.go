package examples

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/client"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func GetAllCharacters(w http.ResponseWriter, r *http.Request) {
	logger.Logger().Info("GetAllObjects")
	integrationDomain, _ := utils.GetIntegrationPath("examples.one.domain")
	integrationContextPath, _ := utils.GetIntegrationPath("examples.one.contextPath")
	integrationPort, _ := utils.GetIntegrationPath("examples.one.port")
	integrationPortInt, _ := strconv.Atoi(integrationPort)
	integrationPath, _ := utils.GetIntegrationPath("examples.one.characters")

	logger.Logger().Info("GetAllObjects", zap.String("integrationDomain", integrationDomain), zap.Int("integrationPort", integrationPortInt), zap.String("integrationPath", integrationPath))

	rc := client.NewRestClient(10 * time.Second)
	req := &client.RequestData{
		Host:        integrationDomain,
		Port:        integrationPortInt,
		UseHTTPS:    true,
		ContextPath: "",
		Path:        integrationContextPath + integrationPath,
		PathVars:    nil,
		QueryParams: nil,
		Headers:     map[string]string{"Accept": "application/json"},
		Body:        nil,
	}

	respBody, status, err := rc.Get(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making request: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Status: %d\n", status)

	if status != http.StatusOK {
		http.Error(w, fmt.Sprintf("Unexpected status code: %d", status), http.StatusInternalServerError)
		return
	}

	var jsonData interface{}
	if err := json.Unmarshal(respBody, &jsonData); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON response: %v", err), http.StatusInternalServerError)
		return
	}
	_ = utils.SendSuccess(w, "SUCCESS", "Get all data of Restful API", http.StatusOK, jsonData)
}

func GetAllPlanets(w http.ResponseWriter, r *http.Request) {
	logger.Logger().Info("GetAllPlanets")
	integrationDomain, _ := utils.GetIntegrationPath("examples.one.domain")
	integrationContextPath, _ := utils.GetIntegrationPath("examples.one.contextPath")
	integrationPort, _ := utils.GetIntegrationPath("examples.one.port")
	integrationPortInt, _ := strconv.Atoi(integrationPort)
	integrationPath, _ := utils.GetIntegrationPath("examples.one.planets")

	logger.Logger().Info("GetAllPlanets", zap.String("integrationDomain", integrationDomain), zap.Int("integrationPort", integrationPortInt), zap.String("integrationPath", integrationPath))

	rc := client.NewRestClient(10 * time.Second)
	req := &client.RequestData{
		Host:        integrationDomain,
		Port:        integrationPortInt,
		UseHTTPS:    true,
		ContextPath: "",
		Path:        integrationContextPath + integrationPath,
		PathVars:    nil,
		QueryParams: nil,
		Headers:     map[string]string{"Accept": "application/json"},
		Body:        nil,
	}

	respBody, status, err := rc.Get(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making request: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Status: %d\n", status)

	if status != http.StatusOK {
		http.Error(w, fmt.Sprintf("Unexpected status code: %d", status), http.StatusInternalServerError)
		return
	}

	var jsonData interface{}
	if err := json.Unmarshal(respBody, &jsonData); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON response: %v", err), http.StatusInternalServerError)
		return
	}
	_ = utils.SendSuccess(w, "SUCCESS", "Get all data of Restful API", http.StatusOK, jsonData)
}
