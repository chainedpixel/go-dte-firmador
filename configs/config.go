package configs

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chainedpixel/go-dte-signer/pkg/logs"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Locale     LocaleConfig     `mapstructure:"locale"`
	Filesystem FilesystemConfig `mapstructure:"filesystem"`
	Log        LogConfig        `mapstructure:"log"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string `mapstructure:"port"`
	SignerRoute  string `mapstructure:"signerroute"`
	HealthRoute  string `mapstructure:"healthroute"`
	ReadTimeout  int    `mapstructure:"readtimeout"`
	WriteTimeout int    `mapstructure:"writetimeout"`
}

// LocaleConfig holds localization configuration
type LocaleConfig struct {
	DefaultLocale string `mapstructure:"defaultlocale"`
	LocalesDir    string `mapstructure:"localesdir"`
}

// FilesystemConfig holds filesystem configuration
type FilesystemConfig struct {
	CertificatesDir string `mapstructure:"certificatesdir"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, bool, error) {
	v := viper.New()
	foundConfig := true

	// Set default values
	v.SetDefault("server.port", "8113")
	v.SetDefault("server.signerroute", "/signer")
	v.SetDefault("server.healthroute", "/health")
	v.SetDefault("server.readtimeout", 15)
	v.SetDefault("server.writetimeout", 15)
	v.SetDefault("locale.defaultlocale", "es")
	v.SetDefault("locale.localesdir", "./configs/locales")
	v.SetDefault("filesystem.certificatesdir", "./uploads/test/")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
	v.SetDefault("log.dir", "./logs")

	// Environment variables (APP_SERVER_PORT, APP_LOCALE_DEFAULTLOCALE, etc.)
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("../")
	v.AddConfigPath(".")

	// Read config file if it exists
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, false, fmt.Errorf("failed to read config file: %w", err)
		}
		foundConfig = false
		fmt.Println("Warning: No config file found, using defaults and environment variables")
	} else {
		fmt.Printf("Using config file: %s\n", v.ConfigFileUsed())
	}

	// Parse config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, true, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config and ensure directories exist
	if err := validateConfig(&config); err != nil {
		return nil, true, fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize logger
	if err := logs.InitLogger(config.Log.Level); err != nil {
		return nil, true, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Log configuration details at debug level
	logConfigDetails(&config)

	return &config, foundConfig, nil
}

// validateConfig validates the provided configuration and ensures required directories exist
func validateConfig(config *Config) error {
	// Validate server configuration
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	// Validate locale configuration
	if config.Locale.DefaultLocale == "" {
		return fmt.Errorf("default locale is required")
	}

	// Ensure locales directory exists
	if _, err := os.Stat(config.Locale.LocalesDir); os.IsNotExist(err) {
		return fmt.Errorf("locales directory does not exist: %s", config.Locale.LocalesDir)
	}

	// Validate filesystem configuration
	if config.Filesystem.CertificatesDir == "" {
		return fmt.Errorf("certificates directory is required")
	}

	// Ensure certificates directory exists
	if err := os.MkdirAll(config.Filesystem.CertificatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create certificates directory: %w", err)
	}

	return nil
}

// logConfigDetails logs the configuration details
func logConfigDetails(config *Config) {
	logs.Debug("Configuration loaded successfully")
	logs.Debug(fmt.Sprintf("Server configuration: port=%s, readTimeout=%d, writeTimeout=%d",
		config.Server.Port, config.Server.ReadTimeout, config.Server.WriteTimeout))
	logs.Debug(fmt.Sprintf("Locale configuration: defaultLocale=%s, localesDir=%s",
		config.Locale.DefaultLocale, config.Locale.LocalesDir))
	logs.Debug(fmt.Sprintf("Filesystem configuration: certificatesDir=%s",
		config.Filesystem.CertificatesDir))
}
