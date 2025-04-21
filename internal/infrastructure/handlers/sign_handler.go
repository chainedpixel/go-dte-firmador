package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/chainedpixel/go-dte-signer/internal/application/usecases"
	"github.com/chainedpixel/go-dte-signer/pkg/logs"
)

// SignHandler handles document signing requests
type SignHandler struct {
	path                   string
	documentSigningUseCase *usecases.DocumentSigningUseCase
}

// RegisterRoutes registers the handler routes with the router
func (h *SignHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(h.path, h.Handle).Methods(http.MethodPost)
}

// NewSignHandler creates a new sign handler
func NewSignHandler(documentSigningUseCase *usecases.DocumentSigningUseCase, path string) *SignHandler {
	return &SignHandler{
		path:                   path,
		documentSigningUseCase: documentSigningUseCase,
	}
}

// Handle handles document signing requests
func (h *SignHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	// 1: Parse the request body
	var input usecases.DocumentSigningInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logs.Error("ERROR: Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	// 2: Execute the use case
	resp, err := h.documentSigningUseCase.Execute(r.Context(), input)
	if err != nil {
		logs.Error("ERROR: Unexpected error in document signing use case: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
		})
		return
	}

	// 3: Determine HTTP status code based on response
	statusCode := http.StatusOK
	if resp.Status != "OK" {
		statusCode = http.StatusBadRequest
	}

	// 4: Write response
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logs.Error("ERROR: Failed to encode response: %v", err)
	}
}
