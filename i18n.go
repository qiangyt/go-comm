package comm

import (
	"embed"
	"fmt"
	"os"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localesFS embed.FS

var (
	defaultBundle    *i18n.Bundle
	defaultLocalizer *i18n.Localizer
	i18nMutex        sync.RWMutex
)

func init() {
	// Initialize with English as default language
	InitI18n("en")
}

// InitI18n initializes the i18n system with the specified language.
// Supported languages: "en" (English), "zh" (Chinese).
// If an unsupported language is provided, it defaults to "en".
func InitI18n(lang string) {
	i18nMutex.Lock()
	defer i18nMutex.Unlock()

	// Validate language
	var defaultLang language.Tag
	switch lang {
	case "zh", "zh-CN", "zh-Hans":
		defaultLang = language.Chinese
	case "en", "en-US":
		defaultLang = language.English
	default:
		defaultLang = language.English
	}

	defaultBundle = i18n.NewBundle(defaultLang)
	defaultBundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// Load embedded locale files
	loadEmbeddedLocales()

	// Create localizer
	defaultLocalizer = i18n.NewLocalizer(defaultBundle, lang)
}

// SetLanguage changes the current language for i18n.
func SetLanguage(lang string) {
	InitI18n(lang)
}

// GetLanguage returns the current language setting from environment variable or default.
func GetLanguage() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LANGUAGE")
	}
	if lang == "" {
		return "en"
	}

	// Parse language code (e.g., "zh_CN.UTF-8" -> "zh")
	if len(lang) >= 2 {
		return lang[:2]
	}
	return "en"
}

func loadEmbeddedLocales() {
	// Load English translations
	if data, err := localesFS.ReadFile("locales/en.yaml"); err == nil {
		defaultBundle.MustParseMessageFileBytes(data, "en.yaml")
	}

	// Load Chinese translations
	if data, err := localesFS.ReadFile("locales/zh.yaml"); err == nil {
		defaultBundle.MustParseMessageFileBytes(data, "zh.yaml")
	}
}

// T translates a message ID with optional template data.
// Usage:
//   T("error.required", map[string]interface{}{"Hint": "config", "Key": "name"})
func T(messageID string, templateData map[string]interface{}) string {
	i18nMutex.RLock()
	defer i18nMutex.RUnlock()

	if defaultLocalizer == nil {
		return messageID
	}

	msg, err := defaultLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		// Fallback to messageID if translation is not found
		return messageID
	}
	return msg
}

// Tf translates a message ID with formatted arguments (printf-style).
// This is a convenience function for simple string formatting.
func Tf(messageID string, args ...interface{}) string {
	i18nMutex.RLock()
	defer i18nMutex.RUnlock()

	if defaultLocalizer == nil {
		return fmt.Sprintf(messageID, args...)
	}

	// Try to get translation first
	msg, err := defaultLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		// If translation not found, use messageID as format string
		return fmt.Sprintf(messageID, args...)
	}

	// Apply formatting to translated message
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// LocalizeError creates a localized error message.
func LocalizeError(messageID string, templateData map[string]interface{}) error {
	return fmt.Errorf("%s", T(messageID, templateData))
}

// LocalizeErrorf creates a localized error message with printf-style formatting.
func LocalizeErrorf(messageID string, args ...interface{}) error {
	return fmt.Errorf("%s", Tf(messageID, args...))
}
