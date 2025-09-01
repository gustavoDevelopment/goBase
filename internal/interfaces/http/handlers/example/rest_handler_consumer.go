package handlers

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/client"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIObject struct {
	ID   string                 `json:"id"`
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

type ExampleHandler struct {
}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

func (h *ExampleHandler) GetAllObjects(w http.ResponseWriter, r *http.Request) {
	rc := client.NewRestClient(10 * time.Second)
	req := &client.RequestData{
		Host:        "api.restful-api.dev",
		Port:        443,
		UseHTTPS:    true,
		ContextPath: "",
		Path:        "/objects",
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

	var items []APIObject
	if status != http.StatusOK {
		http.Error(w, fmt.Sprintf("Unexpected status code: %d", status), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(respBody, &items); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing response: %v", err), http.StatusInternalServerError)
		return
	}

	for _, obj := range items {
		fmt.Printf("ID: %s — Name: %s — Data: %+v\n", obj.ID, obj.Name, obj.Data)
	}
	_ = utils.SendSuccess(w, "SUCCESS", "Get all data of Restful API", http.StatusOK, items)
}
