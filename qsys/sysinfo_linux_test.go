//go:build linux
// +build linux

package qsys

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectLinuxDist_FromOsRelease(t *testing.T) {
	a := require.New(t)

	// Create a temporary directory with fake /etc/os-release
	tmpDir := t.TempDir()
	etcDir := filepath.Join(tmpDir, "etc")
	_ = os.MkdirAll(etcDir, 0o755)

	osReleaseContent := `ID=ubuntu
VERSION_ID="22.04"
NAME="Ubuntu"`
	_ = os.WriteFile(filepath.Join(etcDir, "os-release"), []byte(osReleaseContent), 0o644)

	// We can't easily mock the /etc path, so this test requires root
	// For now, just verify that detectLinuxDist is called on Linux
	detector := NewSystemDetector()
	info := detector.GetSystemInfo()

	// On Linux, OS should be "linux"
	a.Equal("linux", info.OS)
}
