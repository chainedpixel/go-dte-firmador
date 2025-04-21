package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chainedpixel/go-dte-signer/internal/infrastructure/adapters"
	"github.com/chainedpixel/go-dte-signer/pkg/logs"
)

// Server represents an HTTP server
type Server struct {
	httpServer *http.Server
	router     *adapters.Router
}

// NewServer creates a new server
func NewServer(router *adapters.Router, port string, readTimeout, writeTimeout int) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      router.GetHTTPHandler(),
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		router: router,
	}
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	// Create a channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		logs.Info("Server listening on", map[string]interface{}{
			"port": s.httpServer.Addr,
		})
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	// Create a channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive an error or an interrupt.
	select {
	case err := <-serverErrors:
		// Error starting or closing listener
		return fmt.Errorf("server error: %w", err)

	case <-shutdown:
		logs.Info("Server shutdown initiated")

		// Give outstanding requests a deadline for completion
		ctxShutdown, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Shutdown the server gracefully
		if err := s.httpServer.Shutdown(ctxShutdown); err != nil {
			// Error from closing listeners, or context timeout
			logs.Error("Could not gracefully shutdown the server:", map[string]interface{}{
				"error": err.Error(),
			})
			if err := s.httpServer.Close(); err != nil {
				return fmt.Errorf("could not stop server: %w", err)
			}
		}

		// Wait for all connections to finish
		logs.Info("Server gracefully stopped")
		return nil

	case <-ctx.Done():
		logs.Info("Server shutdown initiated by context")
		return s.Stop(ctx)
	}
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
	// Create a deadline to wait for
	ctxShutdown, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := s.httpServer.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("could not gracefully shutdown the server: %w", err)
	}

	logs.Info("Server stopped gracefully")
	return nil
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.httpServer.Handler.ServeHTTP(w, r)
}

// RunWithGracefulShutdown runs the server with graceful shutdown
func RunWithGracefulShutdown(ctx context.Context, server *Server) error {
	// Start server in background
	errChan := make(chan error, 1)
	go func() {
		if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	// Wait for server to finish
	if err := <-errChan; err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
