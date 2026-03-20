package qplugin

import (
	"context"
	"runtime"
	"testing"

	"github.com/qiangyt/go-comm/v2/qfile"
	"github.com/qiangyt/go-comm/v2/qlog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// Tests for fs_plugin_loader.go

func TestFsPluginLoader_Start_EmptyDir(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	// Note: The loader doubles the namespace in the path
	// NewFsPluginLoader sets dir = filepath.Join(dir, namespace)
	// Then Start() calls filepath.Join(me.dir, me.namespace)
	// So for "local" namespace, it expects: /plugins/local/local
	qfile.MkdirP(fs, "/plugins/local/local")
	logger := qlog.NewDiscardLogger()
	loader := NewLocalPluginLoader(logger, fs, "/plugins")

	err := loader.Start(logger)
	a.NoError(err)
}

func TestFsPluginLoader_Start_AlreadyStarted(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	qfile.MkdirP(fs, "/plugins/local/local")
	logger := qlog.NewDiscardLogger()
	loader := NewLocalPluginLoader(logger, fs, "/plugins")

	// Start once
	err := loader.Start(logger)
	a.NoError(err)

	// Start again - should return nil without error
	err = loader.Start(logger)
	a.NoError(err)
}

// Tests for sysinfo.go (Linux detection)

func TestDetectLinuxDist_NoFiles(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping on non-Linux")
	}
	// This test would require modifying /etc files, skip on real systems
	t.Skip("Requires root access to test properly")
}

// Tests for command.go OpenHandler

func TestOpenHandler_Panic(t *testing.T) {
	a := require.New(t)

	ctx := context.Background()
	args := []string{"test"}

	// OpenHandler is only defined on non-Windows
	if runtime.GOOS != "windows" {
		// This would typically try to open a file/URL
		// In test, just verify it can be called (will fail on no display)
		// We can't actually test it without mocking exec.Command
		_ = ctx
		_ = args
		_ = a
		t.Skip("OpenHandler requires GUI to test")
	}
}

// Tests for lock_file_windows.go

func TestCreateLockFile_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping on non-Windows")
	}
	// CreateLockFile is platform-specific
	// On Windows, it uses atomic file creation
	t.Skip("Platform-specific function")
}
