package configs

import (
	"context"
	"fmt"

	"github.com/chainedpixel/go-dte-signer/internal/application/usecases"
	"github.com/chainedpixel/go-dte-signer/internal/domain/services"
	"github.com/chainedpixel/go-dte-signer/internal/infrastructure/adapters"
	"github.com/chainedpixel/go-dte-signer/internal/infrastructure/cypher"
	"github.com/chainedpixel/go-dte-signer/internal/infrastructure/handlers"
	"github.com/chainedpixel/go-dte-signer/internal/infrastructure/server"
	"github.com/chainedpixel/go-dte-signer/pkg/i18n"
	"github.com/chainedpixel/go-dte-signer/pkg/logs"
)

// Application holds all application components
type Application struct {
	Server *server.Server
	Config *Config
}

// Bootstrap initializes the application
func Bootstrap() (*Application, error) {
	// 1. Load configuration
	config, foundConfig, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	if !foundConfig {
		logs.Warn("Config file not found, using default values, if this is not intended, please check the config file")
	}

	logs.Info("Configuration loaded successfully")

	// 2. Initialize logging
	sv, err := initServerDependencies(config)
	if err != nil {
		return nil, err
	}

	// 3. Return bootstrapped application
	app := &Application{
		Server: sv,
		Config: config,
	}

	logs.Info("Application bootstrap completed successfully")
	return app, nil
}

// Start starts the application
func (a *Application) Start(ctx context.Context) error {
	logs.Info(fmt.Sprintf("Starting server on port %s", a.Config.Server.Port))
	return a.Server.Start(ctx)
}

func initServerDependencies(config *Config) (*server.Server, error) {
	// 1. Initialize translator
	logs.Debug("Initializing translator...")
	translator, err := i18n.NewTranslator(config.Locale.LocalesDir, config.Locale.DefaultLocale)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize translator: %w", err)
	}
	logs.Info("Translator initialized successfully")

	// 2. Initialize infrastructure components
	logs.Debug("Initializing infrastructure components...")
	keyProcessor := cypher.NewKeyProcessor()
	jwsSigner := cypher.NewJWSSigner()
	certificateRepository := adapters.NewFileCertificateRepository(
		config.Filesystem.CertificatesDir,
		keyProcessor,
	)
	logs.Info("Infrastructure components initialized successfully")

	// 3. Initialize domain services
	logs.Debug("Initializing domain services...")
	signingService := services.NewSigningService(certificateRepository, jwsSigner)
	logs.Info("Domain services initialized successfully")

	// 4. Initialize application use cases
	logs.Debug("Initializing application use cases...")
	documentSigningUseCase := usecases.NewDocumentSigningUseCase(signingService, translator)
	healthCheckUseCase := usecases.NewHealthCheckUseCase()
	logs.Info("Application use cases initialized successfully")

	// 5. Initialize HTTP handlers
	logs.Debug("Initializing HTTP handlers...")
	signHandler := handlers.NewSignHandler(documentSigningUseCase, config.Server.SignerRoute)
	healthHandler := handlers.NewHealthHandler(healthCheckUseCase, config.Server.HealthRoute)
	logs.Info("HTTP handlers initialized successfully")

	// 6. Initialize router and register routes
	logs.Debug("Initializing router...")
	router := adapters.NewRouter()
	router.RegisterHandler(signHandler)
	router.RegisterHandler(healthHandler)
	logs.Info("Router initialized successfully")

	// 7. Initialize server
	logs.Info("Initializing server...")
	httpServer := server.NewServer(
		router,
		config.Server.Port,
		config.Server.ReadTimeout,
		config.Server.WriteTimeout,
	)
	logs.Info("Server initialized successfully")

	return httpServer, nil
}
