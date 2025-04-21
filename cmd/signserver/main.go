package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chainedpixel/go-dte-signer/configs"
	"github.com/chainedpixel/go-dte-signer/pkg/logs"
)

func main() {
	// Create a context that will be canceled on termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Bootstrap the application
	app, err := configs.Bootstrap()
	if err != nil {
		fmt.Println(err)
	}

	// Log startup information
	logs.Info("Signer server initialized successfully")

	// Start the application
	if err := app.Start(ctx); err != nil {
		logs.Fatal("Application failed", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
