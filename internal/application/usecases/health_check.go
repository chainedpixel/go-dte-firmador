package usecases

import (
	"context"
	"runtime"
	"time"

	"github.com/chainedpixel/go-dte-signer/pkg/response"
)

// HealthCheckUseCase handles health check operations
type HealthCheckUseCase struct {
	startTime time.Time
}

// HealthCheckOutput represents the health check response data
type HealthCheckOutput struct {
	Status    string    `json:"status"`
	Uptime    string    `json:"uptime"`
	Timestamp time.Time `json:"timestamp"`
	GoVersion string    `json:"goVersion"`
}

// NewHealthCheckUseCase creates a new health check use case
func NewHealthCheckUseCase() *HealthCheckUseCase {
	return &HealthCheckUseCase{
		startTime: time.Now(),
	}
}

// Execute performs the health check
func (uc *HealthCheckUseCase) Execute(ctx context.Context) (*response.Response, error) {
	// 1. Get current time
	now := time.Now()

	// 2. Calculate uptime
	uptime := now.Sub(uc.startTime).String()

	// 3. Create health check output
	output := &HealthCheckOutput{
		Status:    "UP",
		Uptime:    uptime,
		Timestamp: now,
		GoVersion: runtime.Version(),
	}

	return response.NewSuccessResponse(output), nil
}
