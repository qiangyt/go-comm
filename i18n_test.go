package comm

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitI18n(t *testing.T) {
	a := require.New(t)

	// Test English initialization
	InitI18n("en")
	a.NotNil(bundle)
	a.NotNil(localizer)

	// Test Chinese initialization
	InitI18n("zh")
	a.NotNil(bundle)
	a.NotNil(localizer)

	// Test invalid language defaults to English
	InitI18n("invalid")
	a.NotNil(bundle)
	a.NotNil(localizer)
}

func TestSetLanguage(t *testing.T) {
	a := require.New(t)

	// Test setting English
	SetLanguage("en")
	a.NotNil(localizer)

	// Test setting Chinese
	SetLanguage("zh")
	a.NotNil(localizer)
}

func TestGetLanguage(t *testing.T) {
	a := require.New(t)

	// Test default language when no env vars
	os.Unsetenv("LANG")
	os.Unsetenv("LANGUAGE")
	InitI18n(DetectLanguage())
	lang := GetLanguage()
	a.Equal("en", lang)

	// Test LANG environment variable
	os.Setenv("LANG", "zh_CN.UTF-8")
	InitI18n(DetectLanguage())
	lang = GetLanguage()
	a.Equal("zh", lang)

	// Test LANGUAGE environment variable
	os.Unsetenv("LANG")
	os.Setenv("LANGUAGE", "en_US.UTF-8")
	InitI18n(DetectLanguage())
	lang = GetLanguage()
	a.Equal("en", lang)

	// Cleanup
	os.Unsetenv("LANG")
	os.Unsetenv("LANGUAGE")
}

func TestT_English(t *testing.T) {
	a := require.New(t)

	// Initialize with English
	InitI18n("en")

	// Test type conversion errors
	msg := T("error.required", map[string]interface{}{
		"Hint": "config",
		"Key":  "name",
	})
	a.Contains(msg, "config.name")
	a.Contains(msg, "required")

	msg = T("error.type.string", map[string]interface{}{
		"Hint":  "value",
		"Type":  "int",
		"Value": "123",
	})
	a.Contains(msg, "value")
	a.Contains(msg, "string")

	// Test plugin errors
	msg = T("error.plugin.version_mismatch", map[string]interface{}{
		"Namespace": "test",
		"Name":      "plugin1",
		"Expected":  2,
		"Actual":    1,
	})
	a.Contains(msg, "test/plugin1")
	a.Contains(msg, "version")

	// Test file errors
	msg = T("error.file.not_found", map[string]interface{}{
		"Path": "/tmp/test.txt",
	})
	a.Contains(msg, "/tmp/test.txt")
	a.Contains(msg, "not found")

	// Test network errors
	msg = T("error.net.interface_not_found", map[string]interface{}{
		"Interface": "eth0",
	})
	a.Contains(msg, "eth0")
}

func TestT_Chinese(t *testing.T) {
	a := require.New(t)

	// Initialize with Chinese
	InitI18n("zh")

	// Test type conversion errors
	msg := T("error.required", map[string]interface{}{
		"Hint": "config",
		"Key":  "name",
	})
	a.Contains(msg, "config.name")
	a.Contains(msg, "必需")

	msg = T("error.type.string", map[string]interface{}{
		"Hint":  "value",
		"Type":  "int",
		"Value": "123",
	})
	a.Contains(msg, "value")
	a.Contains(msg, "字符串")

	// Test plugin errors
	msg = T("error.plugin.version_mismatch", map[string]interface{}{
		"Namespace": "test",
		"Name":      "plugin1",
		"Expected":  2,
		"Actual":    1,
	})
	a.Contains(msg, "test/plugin1")
	a.Contains(msg, "版本")

	// Test file errors
	msg = T("error.file.not_found", map[string]interface{}{
		"Path": "/tmp/test.txt",
	})
	a.Contains(msg, "/tmp/test.txt")
	a.Contains(msg, "未找到")

	// Test network errors
	msg = T("error.net.interface_not_found", map[string]interface{}{
		"Interface": "eth0",
	})
	a.Contains(msg, "eth0")
	a.Contains(msg, "接口")
}

func TestT_MissingTranslation(t *testing.T) {
	a := require.New(t)

	InitI18n("en")

	// Test with missing message ID - should return the messageID
	msg := T("non.existent.message", map[string]interface{}{
		"Param": "value",
	})
	a.Equal("non.existent.message", msg)
}

func TestTf(t *testing.T) {
	a := require.New(t)

	InitI18n("en")

	// Test with simple format string (no translation)
	msg := Tf("Hello %s, you have %d messages", "John", 5)
	a.Contains(msg, "John")
	a.Contains(msg, "5")
}

func TestTf_withNilLocalizer(t *testing.T) {
	a := require.New(t)

	// Temporarily set localizer to nil to test fallback
	mutex.Lock()
	savedLocalizer := localizer
	localizer = nil
	mutex.Unlock()

	// Should use fmt.Sprintf directly
	msg := Tf("Count: %d", 42)
	a.Contains(msg, "42")

	// Restore localizer
	mutex.Lock()
	localizer = savedLocalizer
	mutex.Unlock()
}

func TestTf_nonExistentMessage(t *testing.T) {
	a := require.New(t)

	InitI18n("en")

	// Test with non-existent message ID, should fallback to format string
	msg := Tf("nonexistent.message.%s", "test")
	a.Contains(msg, "test")
}

func TestLocalizeError(t *testing.T) {
	a := require.New(t)

	InitI18n("en")

	// Test LocalizeError
	err := LocalizeError("error.required", map[string]interface{}{
		"Hint": "config",
		"Key":  "name",
	})
	a.NotNil(err)
	a.Contains(err.Error(), "config.name")
	a.Contains(err.Error(), "required")
}

func TestLocalizeErrorf(t *testing.T) {
	a := require.New(t)

	InitI18n("en")

	// Test LocalizeErrorf with simple format
	err := LocalizeErrorf("Error: %s at line %d", "syntax error", 42)
	a.NotNil(err)
	a.Contains(err.Error(), "syntax error")
	a.Contains(err.Error(), "42")
}

func TestT_AllMessageIDs(t *testing.T) {
	a := require.New(t)

	testCases := []struct {
		messageID string
		lang      string
		data      map[string]interface{}
	}{
		// Type conversion errors
		{"error.required", "en", map[string]interface{}{"Hint": "test", "Key": "key"}},
		{"error.type.string", "en", map[string]interface{}{"Hint": "test", "Type": "int", "Value": "123"}},
		{"error.type.string_array", "en", map[string]interface{}{"Hint": "test", "Type": "int", "Value": "123"}},
		{"error.type.string_map", "en", map[string]interface{}{"Hint": "test", "Type": "int", "Value": "123"}},
		{"error.type.int", "en", map[string]interface{}{"Hint": "test", "Type": "string", "Value": "abc"}},
		{"error.type.int_array", "en", map[string]interface{}{"Hint": "test", "Type": "string", "Value": "abc"}},
		{"error.type.int_map", "en", map[string]interface{}{"Hint": "test", "Type": "string", "Value": "abc"}},
		{"error.type.float", "en", map[string]interface{}{"Hint": "test", "Type": "string", "Value": "abc"}},
		{"error.type.float_array", "en", map[string]interface{}{"Hint": "test", "Type": "string", "Value": "abc"}},
		{"error.type.float_map", "en", map[string]interface{}{"Hint": "test", "Type": "string", "Value": "abc"}},
		{"error.type.bool", "en", map[string]interface{}{"Hint": "test", "Type": "string", "Value": "abc"}},
		{"error.type.bool_array", "en", map[string]interface{}{"Hint": "test", "Type": "string", "Value": "abc"}},
		{"error.type.bool_map", "en", map[string]interface{}{"Hint": "test", "Type": "string", "Value": "abc"}},

		// Plugin errors
		{"error.plugin.version_mismatch", "en", map[string]interface{}{"Namespace": "ns", "Name": "plugin", "Expected": 2, "Actual": 1}},
		{"error.plugin.unsupported_kind", "en", map[string]interface{}{"Namespace": "ns", "Name": "plugin", "Kind": "unknown"}},
		{"error.plugin.namespace_not_specified", "en", map[string]interface{}{"Loader": "loader"}},
		{"error.plugin.namespace_already_registered", "en", map[string]interface{}{"Namespace": "ns", "Loader": "loader"}},
		{"error.plugin.duplicated_kind", "en", map[string]interface{}{"Namespace": "ns", "Kind": "kind", "Loader": "loader"}},
		{"error.plugin.start_failed", "en", map[string]interface{}{"PluginId": "ns/plugin", "Version": "1.0", "Cause": "error"}},
		{"error.plugin.stop_failed", "en", map[string]interface{}{"PluginId": "ns/plugin", "Version": "1.0", "Cause": "error"}},

		// Plugin log messages
		{"log.plugin.loader.starting", "en", nil},
		{"log.plugin.loader.started", "en", nil},
		{"log.plugin.loader.start_failed", "en", nil},
		{"log.plugin.loader.stopping", "en", nil},
		{"log.plugin.loader.stopped", "en", nil},
		{"log.plugin.loader.stop_failed", "en", nil},

		// File errors
		{"error.file.already_exists", "en", map[string]interface{}{"Path": "/tmp/test"}},
		{"error.file.not_found", "en", map[string]interface{}{"Path": "/tmp/test"}},
		{"error.dir.already_exists", "en", map[string]interface{}{"Path": "/tmp/test"}},
		{"error.dir.not_found", "en", map[string]interface{}{"Path": "/tmp/test"}},
		{"error.file.expect_file_but_dir", "en", map[string]interface{}{"Path": "/tmp/test"}},
		{"error.file.expect_dir_but_file", "en", map[string]interface{}{"Path": "/tmp/test"}},
		{"error.file.cannot_remove_root", "en", nil},

		// Env errors
		{"error.env.zero_length_string", "en", nil},
		{"error.env.cannot_separate_key_value", "en", nil},

		// Network errors
		{"error.net.interface_not_found", "en", map[string]interface{}{"Interface": "eth0"}},
	}

	for _, tc := range testCases {
		t.Run(tc.messageID+"_"+tc.lang, func(t *testing.T) {
			InitI18n(tc.lang)
			msg := T(tc.messageID, tc.data)
			a.NotEmpty(msg)
			// Should not return the message ID itself (means translation was found)
			if tc.data != nil {
				a.NotEqual(tc.messageID, msg, "Translation not found for %s", tc.messageID)
			}
		})
	}

	// Test all messages in Chinese too
	for _, tc := range testCases {
		tc.lang = "zh"
		t.Run(tc.messageID+"_"+tc.lang, func(t *testing.T) {
			InitI18n(tc.lang)
			msg := T(tc.messageID, tc.data)
			a.NotEmpty(msg)
			if tc.data != nil {
				a.NotEqual(tc.messageID, msg, "Translation not found for %s", tc.messageID)
			}
		})
	}
}

func TestConcurrentLanguageSwitch(t *testing.T) {
	a := require.New(t)

	// Test concurrent language switching
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			SetLanguage("en")
			_ = T("error.required", map[string]interface{}{"Hint": "test", "Key": "key"})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			SetLanguage("zh")
			_ = T("error.required", map[string]interface{}{"Hint": "test", "Key": "key"})
		}
		done <- true
	}()

	<-done
	<-done

	// Should not crash
	a.NotNil(localizer)
}
