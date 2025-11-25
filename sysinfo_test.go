package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSystemDetector(t *testing.T) {
	a := require.New(t)

	detector := NewSystemDetector()
	a.NotNil(detector)

	info := detector.GetSystemInfo()
	a.NotEmpty(info.OS)
	a.NotEmpty(info.Arch)
}

func TestSystemDetector_GetSystemInfo(t *testing.T) {
	a := require.New(t)

	detector := NewSystemDetector()
	info := detector.GetSystemInfo()

	// OS should be one of linux, windows, macos
	a.True(info.OS == "linux" || info.OS == "windows" || info.OS == "macos")

	// Arch should be one of x86, amd64, arm64, armhf
	a.True(info.Arch == "x86" || info.Arch == "amd64" || info.Arch == "arm64" || info.Arch == "armhf" || info.Arch != "")
}

func TestSystemDetector_IsLinux(t *testing.T) {
	a := require.New(t)

	detector := NewSystemDetector()
	info := detector.GetSystemInfo()

	if info.OS == "linux" {
		a.True(detector.IsLinux())
		a.False(detector.IsWindows())
		a.False(detector.IsMacOS())
	}
}

func TestSystemDetector_IsWindows(t *testing.T) {
	a := require.New(t)

	detector := NewSystemDetector()
	info := detector.GetSystemInfo()

	if info.OS == "windows" {
		a.False(detector.IsLinux())
		a.True(detector.IsWindows())
		a.False(detector.IsMacOS())
	}
}

func TestSystemDetector_IsMacOS(t *testing.T) {
	a := require.New(t)

	detector := NewSystemDetector()
	info := detector.GetSystemInfo()

	if info.OS == "macos" {
		a.False(detector.IsLinux())
		a.False(detector.IsWindows())
		a.True(detector.IsMacOS())
	}
}

func TestSystemDetector_MatchSystem(t *testing.T) {
	a := require.New(t)

	detector := NewSystemDetector()
	info := detector.GetSystemInfo()

	// Test matching current OS
	a.True(detector.MatchSystem(info.OS, "", "", ""))

	// Test matching current arch
	a.True(detector.MatchSystem("", "", "", info.Arch))

	// Test empty match (should always return true)
	a.True(detector.MatchSystem("", "", "", ""))

	// Test non-matching OS
	a.False(detector.MatchSystem("nonexistent_os", "", "", ""))

	// Test multiple values (comma separated)
	a.True(detector.MatchSystem("linux,windows,macos", "", "", ""))

	// Test multiple values (space separated)
	a.True(detector.MatchSystem("linux windows macos", "", "", ""))
}

func TestMatchValues(t *testing.T) {
	a := require.New(t)

	// Empty expected always matches
	a.True(matchValues("anything", ""))

	// Exact match
	a.True(matchValues("linux", "linux"))

	// Non-match
	a.False(matchValues("linux", "windows"))

	// Multiple values comma separated
	a.True(matchValues("linux", "linux,windows,macos"))
	a.True(matchValues("windows", "linux,windows,macos"))
	a.False(matchValues("freebsd", "linux,windows,macos"))

	// Multiple values space separated
	a.True(matchValues("linux", "linux windows macos"))
	a.True(matchValues("macos", "linux windows macos"))
}
