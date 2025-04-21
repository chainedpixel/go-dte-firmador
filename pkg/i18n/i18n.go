package i18n

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"

	"github.com/chainedpixel/go-dte-signer/pkg/logs"
)

// Translator handles message translations
type Translator struct {
	messages map[string]map[string]string
	mutex    sync.RWMutex
	locale   string
}

// NewTranslator creates a new translator instance
func NewTranslator(localesDir string, defaultLocale string) (*Translator, error) {
	logs.Debug(fmt.Sprintf("Initializing translator with localesDir=%s, defaultLocale=%s", localesDir, defaultLocale))

	t := &Translator{
		messages: make(map[string]map[string]string),
		locale:   defaultLocale,
	}

	// Verify the locales directory exists
	if _, err := os.Stat(localesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read locales directory: %w", err)
	}

	// Load all locale files from the specified directory
	files, err := os.ReadDir(localesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read locales directory: %w", err)
	}

	logs.Debug(fmt.Sprintf("Found %d files in locales directory", len(files)))

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		locale := file.Name()[:len(file.Name())-5]
		filePath := filepath.Join(localesDir, file.Name())

		logs.Debug(fmt.Sprintf("Loading locale file: %s", filePath))

		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read locale file %s: %w", file.Name(), err)
		}

		messages := make(map[string]string)
		if err := yaml.Unmarshal(data, &messages); err != nil {
			return nil, fmt.Errorf("failed to parse locale file %s: %w", file.Name(), err)
		}

		t.messages[locale] = messages
		logs.Debug(fmt.Sprintf("Loaded %d messages for locale %s", len(messages), locale))
	}

	// Verify at least the default locale exists
	if _, exists := t.messages[defaultLocale]; !exists {
		return nil, fmt.Errorf("default locale %s not found", defaultLocale)
	}

	logs.Debug(fmt.Sprintf("Translator initialized successfully with %d locales", len(t.messages)))
	return t, nil
}

// SetLocale sets the current locale
func (t *Translator) SetLocale(locale string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exists := t.messages[locale]; exists {
		t.locale = locale
		logs.Debug(fmt.Sprintf("Locale set to %s", locale))
	} else {
		logs.Warn(fmt.Sprintf("Attempted to set unknown locale: %s", locale))
	}
}

// GetLocale returns the current locale
func (t *Translator) GetLocale() string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.locale
}

// Translate returns a translated message by key for the current locale
func (t *Translator) Translate(key string, args ...interface{}) string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	msg, ok := t.messages[t.locale][key]
	if !ok {
		// Fallback to default message or just return the key if not found
		logs.Debug(fmt.Sprintf("Translation key not found: %s", key))
		return key
	}

	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// T is a shorthand for Translate
func (t *Translator) T(key string, args ...interface{}) string {
	return t.Translate(key, args...)
}
