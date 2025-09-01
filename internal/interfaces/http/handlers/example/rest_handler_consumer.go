package handlers

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/domain/examples"
	"net/http"
)

type ExampleHandler struct {
}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

func (h *ExampleHandler) GetAllCharacters(w http.ResponseWriter, r *http.Request) {
	examples.GetAllCharacters(w, r)
}

func (h *ExampleHandler) GetAllPlanets(w http.ResponseWriter, r *http.Request) {
	examples.GetAllPlanets(w, r)
}
