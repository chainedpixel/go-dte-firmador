package logs

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// Log level mapping
var logsLevel = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
}

// CustomFormatter defines a custom format for logs
type CustomFormatter struct {
	logrus.TextFormatter
}

// Format implements the logrus.Formatter interface
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColors := map[logrus.Level]string{
		logrus.DebugLevel: "\033[37m", // Gray for debug
		logrus.InfoLevel:  "\033[32m", // Green for info
		logrus.WarnLevel:  "\033[33m", // Yellow for warning
		logrus.ErrorLevel: "\033[31m", // Red for error
		logrus.FatalLevel: "\033[35m", // Purple for fatal
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	color := levelColors[entry.Level]

	baseMessage := fmt.Sprintf("%s %s [%s]\033[0m %s",
		color,
		timestamp,
		level,
		entry.Message,
	)

	var fields string
	if len(entry.Data) > 0 {
		fields = "\n	Details:\n"
		for k, v := range entry.Data {
			fields += fmt.Sprintf("	%-10s: %v\n", k, v)
		}
	}

	return []byte(fmt.Sprintf("%s%s\n", baseMessage, fields)), nil
}

// InitLogger initializes the logger with the specified log level and log directory
func InitLogger(logLevel string) error {
	Logger = logrus.New()

	formatter := &CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		},
	}

	Logger.SetFormatter(formatter)
	Logger.SetLevel(determineLogLevel(logLevel))

	// Set output in console
	writer := io.Writer(os.Stdout)
	Logger.SetOutput(writer)

	Debug("Logger initialized with level: %s", logLevel)

	return nil
}

// formatMessage handles different message formatting
func formatMessage(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	// If only one argument, just convert to string
	if len(args) == 1 {
		return fmt.Sprint(args[0])
	}

	// Check if first argument is a string that might be a format
	if format, ok := args[0].(string); ok {
		// If it contains format specifiers, use fmt.Sprintf
		if strings.Contains(format, "%") {
			return fmt.Sprintf(format, args[1:]...)
		}

		// Otherwise, join all arguments if they're strings
		allStrings := true
		parts := make([]string, len(args))
		for i, arg := range args {
			if _, ok := arg.(string); !ok {
				allStrings = false
				break
			}
			parts[i] = fmt.Sprint(arg)
		}

		if allStrings {
			return strings.Join(parts, " ")
		}
	}

	// Fall back to simple concatenation using fmt.Sprint
	return fmt.Sprint(args...)
}

// Helper function for logging with fields and variable arguments
func logWithFields(level logrus.Level, fields map[string]interface{}, args ...interface{}) {
	if Logger == nil {
		// Create a basic logger if not initialized
		Logger = logrus.New()
		Logger.SetFormatter(&CustomFormatter{})
	}

	message := formatMessage(args...)

	if fields != nil {
		Logger.WithFields(fields).Log(level, message)
	} else {
		Logger.Log(level, message)
	}
}

// Debug logs a debug level message with formatting support
func Debug(args ...interface{}) {
	var fields map[string]interface{}
	logWithFields(logrus.DebugLevel, fields, args...)
}

// DebugWithFields logs a debug level message with fields and formatting support
func DebugWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(logrus.DebugLevel, fields, args...)
}

// Info logs an info level message with formatting support
func Info(args ...interface{}) {
	var fields map[string]interface{}
	logWithFields(logrus.InfoLevel, fields, args...)
}

// InfoWithFields logs an info level message with fields and formatting support
func InfoWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(logrus.InfoLevel, fields, args...)
}

// Warn logs a warning level message with formatting support
func Warn(args ...interface{}) {
	var fields map[string]interface{}
	logWithFields(logrus.WarnLevel, fields, args...)
}

// WarnWithFields logs a warning level message with fields and formatting support
func WarnWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(logrus.WarnLevel, fields, args...)
}

// Error logs an error level message with formatting support
func Error(args ...interface{}) {
	var fields map[string]interface{}
	logWithFields(logrus.ErrorLevel, fields, args...)
}

// ErrorWithFields logs an error level message with fields and formatting support
func ErrorWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(logrus.ErrorLevel, fields, args...)
}

// Fatal logs a fatal level message with formatting support and exits
func Fatal(args ...interface{}) {
	var fields map[string]interface{}
	logWithFields(logrus.FatalLevel, fields, args...)
	os.Exit(1)
}

// FatalWithFields logs a fatal level message with fields and formatting support and exits
func FatalWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(logrus.FatalLevel, fields, args...)
	os.Exit(1)
}

// determineLogLevel converts a string log level to a logrus.Level
func determineLogLevel(logLevel string) logrus.Level {
	if level, ok := logsLevel[strings.ToLower(logLevel)]; ok {
		return level
	}
	return logrus.InfoLevel
}
