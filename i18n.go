package comm

import (
	"embed"
	"fmt"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localesFS embed.FS

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	mutex     sync.RWMutex
)

func init() {
	// Initialize with English as default language
	InitI18n("en")
}

// InitI18n initializes the i18n system with the specified language.
// Supported languages: "en" (English), "zh" (Chinese).
// If an unsupported language is provided, it defaults to "en".
func InitI18n(lang string) {
	mutex.Lock()
	defer mutex.Unlock()

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

	bundle = i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load embedded locale files
	loadEmbeddedLocales()

	// Create localizer
	localizer = i18n.NewLocalizer(bundle, lang)
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
	if data, err := localesFS.ReadFile("locales/active.en.toml"); err == nil {
		bundle.MustParseMessageFileBytes(data, "active.en.toml")
	}

	// Load Chinese translations
	if data, err := localesFS.ReadFile("locales/active.zh.toml"); err == nil {
		bundle.MustParseMessageFileBytes(data, "active.zh.toml")
	}
}

// T translates a message ID with optional template data.
// Usage:
//   T("error.required", map[string]interface{}{"Hint": "config", "Key": "name"})
func T(messageID string, templateData map[string]interface{}) string {
	mutex.RLock()
	defer mutex.RUnlock()

	if localizer == nil {
		return messageID
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
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
	mutex.RLock()
	defer mutex.RUnlock()

	if localizer == nil {
		return fmt.Sprintf(messageID, args...)
	}

	// Try to get translation first
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
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

// localize localizes a message with optional template data.
// This is the unified localization function used by FileOps and other components.
func localize(id string, args ...map[string]any) string {
	mutex.RLock()
	defer mutex.RUnlock()

	cfg := &i18n.LocalizeConfig{
		MessageID:   id,
		PluralCount: 2,
	}
	if len(args) > 0 {
		cfg.TemplateData = args[0]
	}

	r, err := localizer.Localize(cfg)
	if err != nil {
		return id
	}
	return r
}

// LocalizeFunc is the i18n localization function type
type LocalizeFunc func(id string, args ...map[string]any) string

var (
	// DefaultLocalizeFunc is the default localization function (returns id directly)
	DefaultLocalizeFunc LocalizeFunc = func(id string, args ...map[string]any) string {
		return id
	}
)

// NewLocalizedFileOps creates a FileOps with built-in i18n support
func NewLocalizedFileOps(fs afero.Fs) FileOps {
	ops := NewFileOps(fs)
	ops.SetLocalizeFunc(localize)
	return ops
}

// CommLocalize is an alias for localize, for backward compatibility
func CommLocalize(id string, args ...map[string]any) string {
	return localize(id, args...)
}

// SetCommLang is an alias for SetLanguage, for backward compatibility
func SetCommLang(lang string) {
	SetLanguage(lang)
}
