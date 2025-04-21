package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/chainedpixel/go-dte-signer/internal/application/usecases"
	"github.com/chainedpixel/go-dte-signer/pkg/logs"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	path               string
	healthCheckUseCase *usecases.HealthCheckUseCase
}

// RegisterRoutes registers the handler routes with the router
func (h *HealthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(h.path, h.Handle).Methods(http.MethodGet)
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(healthCheckUseCase *usecases.HealthCheckUseCase, path string) *HealthHandler {
	return &HealthHandler{
		path:               path,
		healthCheckUseCase: healthCheckUseCase,
	}
}

// Handle handles health check requests
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// 1: Set response headers
	w.Header().Set("Content-Type", "application/json")

	// 2: Execute health check use case
	resp, err := h.healthCheckUseCase.Execute(r.Context())
	if err != nil {
		logs.Error("ERROR: Unexpected error in health check use case: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
		})
		return
	}

	// 3: Write successful response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logs.Error("ERROR: Failed to encode response: %v", err)
	}
}
