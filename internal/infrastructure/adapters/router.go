package adapters

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/chainedpixel/go-dte-signer/pkg/logs"
)

// Handler represents an HTTP handler with route registration capability
type Handler interface {
	RegisterRoutes(router *mux.Router)
}

// Router manages HTTP routing for the application
type Router struct {
	router *mux.Router
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

// RegisterHandler registers a handler with the router
func (r *Router) RegisterHandler(handler Handler) {
	handler.RegisterRoutes(r.router)
}

// GetHTTPHandler returns the HTTP handler for the router
func (r *Router) GetHTTPHandler() http.Handler {
	// Add common middleware
	handler := loggingMiddleware(r.router)
	handler = recoveryMiddleware(handler)
	handler = corsMiddleware(handler)

	return handler
}

// Router returns the underlying mux.Router
func (r *Router) Router() *mux.Router {
	return r.router
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response wrapper to capture the status code
		wrapper := newResponseWriter(w)

		// Process request
		next.ServeHTTP(wrapper, r)

		// Log after processing
		duration := time.Since(start)
		logs.Info("[%s] %s %s %d %s", r.Method, r.RequestURI, r.RemoteAddr, wrapper.status, duration)
	})
}

// recoveryMiddleware recovers from panics
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logs.Error("PANIC: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware handles CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
