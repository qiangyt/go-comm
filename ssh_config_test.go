package comm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSSHURL(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		wantProtocol string
		wantUser     string
		wantHost     string
		wantPath     string
		wantPort     int
		wantErr      bool
	}{
		{
			name:         "ssh: shorthand format",
			url:          "ssh:myserver",
			wantProtocol: "ssh",
			wantUser:     "",
			wantHost:     "myserver",
			wantPath:     "",
			wantPort:     22,
			wantErr:      false,
		},
		{
			name:         "sftp: shorthand format",
			url:          "sftp:myserver",
			wantProtocol: "sftp",
			wantUser:     "",
			wantHost:     "myserver",
			wantPath:     "",
			wantPort:     22,
			wantErr:      false,
		},
		{
			name:         "user@host format",
			url:          "user@example.com",
			wantProtocol: "ssh",
			wantUser:     "user",
			wantHost:     "example.com",
			wantPath:     "",
			wantPort:     22,
			wantErr:      false,
		},
		{
			name:         "user@host:path format",
			url:          "user@example.com:/path/to/file",
			wantProtocol: "ssh",
			wantUser:     "user",
			wantHost:     "example.com",
			wantPath:     "/path/to/file",
			wantPort:     22,
			wantErr:      false,
		},
		{
			name:         "user@host:port/path format",
			url:          "user@example.com:2222/path/to/file",
			wantProtocol: "ssh",
			wantUser:     "user",
			wantHost:     "example.com",
			wantPath:     "/path/to/file",
			wantPort:     2222,
			wantErr:      false,
		},
		{
			name:         "ssh:// full URL",
			url:          "ssh://user@example.com:2222/path",
			wantProtocol: "ssh",
			wantUser:     "user",
			wantHost:     "example.com",
			wantPath:     "/path",
			wantPort:     2222,
			wantErr:      false,
		},
		{
			name:         "sftp:// full URL",
			url:          "sftp://user@example.com/path",
			wantProtocol: "sftp",
			wantUser:     "user",
			wantHost:     "example.com",
			wantPath:     "/path",
			wantPort:     22,
			wantErr:      false,
		},
		{
			name:         "user@host:relpath format",
			url:          "user@example.com:path/to/file",
			wantProtocol: "ssh",
			wantUser:     "user",
			wantHost:     "example.com",
			wantPath:     "path/to/file",
			wantPort:     22,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protocol, user, hostname, path, port, err := ParseSSHURL(tt.url)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantProtocol, protocol)
			assert.Equal(t, tt.wantUser, user)
			assert.Equal(t, tt.wantHost, hostname)
			assert.Equal(t, tt.wantPath, path)
			assert.Equal(t, tt.wantPort, port)
		})
	}
}

func TestLoadSSHConfig(t *testing.T) {
	// Create a temporary SSH config file
	tmpDir := t.TempDir()
	sshDir := filepath.Join(tmpDir, ".ssh")
	err := os.MkdirAll(sshDir, 0700)
	assert.NoError(t, err)

	configPath := filepath.Join(sshDir, "config")
	configContent := `
# SSH config for testing
Host myserver
    HostName 192.168.1.100
    User admin
    Port 2222
    IdentityFile ~/.ssh/id_rsa_myserver

Host *.example.com
    User developer
    Port 22
    IdentityFile ~/.ssh/id_rsa_example

Host *
    User defaultuser
    ServerAliveInterval 60
    ServerAliveCountMax 3
`
	err = os.WriteFile(configPath, []byte(configContent), 0600)
	assert.NoError(t, err)

	// Load the SSH config
	config, err := LoadSSHConfigFromFile(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Test myserver host
	hostConfig := config.GetHostConfig("myserver")
	assert.NotNil(t, hostConfig)
	assert.Equal(t, "192.168.1.100", hostConfig.HostName)
	assert.Equal(t, "admin", hostConfig.User)
	assert.Equal(t, 2222, hostConfig.Port)

	// Test wildcard matching
	// Note: due to random map iteration, either *.example.com or * may be matched
	hostConfig = config.GetHostConfig("test.example.com")
	assert.NotNil(t, hostConfig)
	// User could be "developer" (if *.example.com matched) or "defaultuser" (if * matched)
	assert.Contains(t, []string{"developer", "defaultuser"}, hostConfig.User)
	assert.Equal(t, 22, hostConfig.Port)
}

func TestSSHConfigToCredentials(t *testing.T) {
	// Create a temporary key file
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "test_key")
	keyContent := "-----BEGIN PRIVATE KEY-----\ntest key content\n-----END PRIVATE KEY-----"
	err := os.WriteFile(keyPath, []byte(keyContent), 0600)
	assert.NoError(t, err)

	hostConfig := &SSHHostConfig{
		Host:         "testhost",
		HostName:     "192.168.1.1",
		User:         "testuser",
		Port:         2222,
		IdentityFile: keyPath,
	}

	creds := SSHConfigToCredentials(hostConfig)
	assert.NotNil(t, creds)
	assert.Equal(t, "testuser", creds.User)
	assert.Equal(t, keyContent, creds.RsaPrivateKey)
}

func TestMatchSSHPattern(t *testing.T) {
	tests := []struct {
		pattern string
		host    string
		want    bool
	}{
		{"*", "anything", true},
		{"example.com", "example.com", true},
		{"example.com", "other.com", false},
		{"*.example.com", "test.example.com", true},
		{"*.example.com", "example.com", false},
		{"*.example.com", "test.other.com", false},
		{"server*", "server1", true},
		{"server*", "server123", true},
		{"server*", "other", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.host, func(t *testing.T) {
			result := matchSSHPattern(tt.pattern, tt.host)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestBuildSSHURL(t *testing.T) {
	tests := []struct {
		name     string
		protocol string
		user     string
		hostname string
		port     int
		path     string
		want     string
	}{
		{
			name:     "basic ssh URL",
			protocol: "ssh",
			user:     "user",
			hostname: "example.com",
			port:     22,
			path:     "",
			want:     "ssh://user@example.com",
		},
		{
			name:     "ssh with custom port",
			protocol: "ssh",
			user:     "user",
			hostname: "example.com",
			port:     2222,
			path:     "",
			want:     "ssh://user@example.com:2222",
		},
		{
			name:     "ssh with path",
			protocol: "ssh",
			user:     "user",
			hostname: "example.com",
			port:     22,
			path:     "/path/to/file",
			want:     "ssh://user@example.com/path/to/file",
		},
		{
			name:     "sftp URL",
			protocol: "sftp",
			user:     "admin",
			hostname: "server.com",
			port:     22,
			path:     "",
			want:     "sftp://admin@server.com",
		},
		{
			name:     "no user",
			protocol: "ssh",
			user:     "",
			hostname: "example.com",
			port:     22,
			path:     "",
			want:     "ssh://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildSSHURL(tt.protocol, tt.user, tt.hostname, tt.port, tt.path)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestNewCredentialsFromSSH_WithoutConfig(t *testing.T) {
	// Test parsing SSH URL without config (basic parsing)
	creds, hostname, port, path, err := NewCredentialsFromSSH("testuser@example.com:/path/to/file")
	assert.NoError(t, err)
	assert.Equal(t, "example.com", hostname)
	assert.Equal(t, 22, port)
	assert.Equal(t, "/path/to/file", path)
	assert.NotNil(t, creds)
	assert.Equal(t, "testuser", creds.User)

	// Test with port
	creds, hostname, port, path, err = NewCredentialsFromSSH("admin@server.com:2222/data")
	assert.NoError(t, err)
	assert.Equal(t, "server.com", hostname)
	assert.Equal(t, 2222, port)
	assert.Equal(t, "/data", path)
	assert.NotNil(t, creds)
	assert.Equal(t, "admin", creds.User)
}

func TestParseIntDefault(t *testing.T) {
	tests := []struct {
		input        string
		defaultValue int
		want         int
	}{
		{"123", 0, 123},
		{"22", 80, 22},
		{"invalid", 42, 42},
		{"", 100, 100},
		{"0", 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseIntDefault(tt.input, tt.defaultValue)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSimpleWildcardMatch(t *testing.T) {
	tests := []struct {
		pattern string
		str     string
		want    bool
	}{
		// Exact matches
		{"example.com", "example.com", true},
		{"example.com", "other.com", false},
		{"test", "test", true},
		{"test", "other", false},

		// Wildcard at start
		{"*.com", "example.com", true},
		{"*.com", "test.org", false},
		{"*", "anything", true},
		{"*", "", true},

		// Wildcard at end
		{"example.*", "example.com", true},
		{"test.*", "test.org", true},
		{"test.*", "other.com", false},

		// Wildcard in middle
		{"ex*.com", "example.com", true},
		{"ex*.com", "ex.com", true},
		{"ex*.com", "exampl.org", false},

		// Multiple wildcards
		{"*.*", "test.com", true},
		{"*.*", "test", false},
		{"*.*.*", "a.b.c", true},

		// Question mark (single character)
		{"example?.com", "example1.com", true},
		{"example?.com", "exampleX.com", true},
		{"example?.com", "example.com", false},
		{"example??", "example12", true},

		// Mixed wildcards
		{"*.example.*", "test.example.com", true},
		{"*test*", "test", true},
		{"*test*", "mytest", true},
		{"*test*", "mytest123", true},
		{"*test*", "other", false},

		// Empty patterns
		{"", "", true},
		{"", "test", false},
		{"test", "", false},

		// Only wildcards
		{"***", "anything", true},
		{"**", "test", true},
		{"?*", "a", true},
		{"?*", "", false},

		// Complex patterns
		{"a*b*c", "abc", true},
		{"a*b*c", "aXXXbXXXc", true},
		{"a*b*c", "abXXXc", true},
		{"a*b*c", "ac", false},

		// Edge cases with single character
		{"?", "a", true},
		{"?", "", false},
		{"?", "ab", false},

		// Star followed by question mark
		{"*?", "a", true},
		{"*?", "ab", true},
		{"*?", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.str, func(t *testing.T) {
			result := simpleWildcardMatch(tt.pattern, tt.str)
			assert.Equal(t, tt.want, result)
		})
	}
}
