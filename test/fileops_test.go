package test

import (
	comm "github.com/qiangyt/go-comm/v2"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/afero"
)

func TestFileOps_CopyFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	// 创建源文件
	srcFile := filepath.Join(tmpDir, "source.txt")
	testContent := []byte("test content")
	if err := afero.WriteFile(fs, srcFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// 复制文件
	destFile := filepath.Join(tmpDir, "subdir", "dest.txt")
	ops.CopyFile(srcFile, destFile)

	// 验证目标文件存在且内容正确
	destContent, err := afero.ReadFile(fs, destFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(destContent) != string(testContent) {
		t.Errorf("Destination content = %s, want %s", destContent, testContent)
	}
}

func TestFileOps_MoveFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	// 创建源文件
	srcFile := filepath.Join(tmpDir, "source.txt")
	testContent := []byte("test content")
	if err := afero.WriteFile(fs, srcFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// 移动文件
	destFile := filepath.Join(tmpDir, "dest.txt")
	ops.MoveFile(srcFile, destFile)

	// 验证源文件不存在
	if exists, _ := afero.Exists(fs, srcFile); exists {
		t.Error("Source file should not exist after move")
	}

	// 验证目标文件存在且内容正确
	destContent, err := afero.ReadFile(fs, destFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(destContent) != string(testContent) {
		t.Errorf("Destination content = %s, want %s", destContent, testContent)
	}
}

func TestFileOps_EnsureDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	// 创建多级目录
	testDir := filepath.Join(tmpDir, "a", "b", "c")
	ops.EnsureDir(testDir)

	// 验证目录存在
	if !ops.DirExists(testDir) {
		t.Error("Directory should exist after EnsureDir()")
	}

	// 再次调用应该不会出错
	ops.EnsureDir(testDir)
}

func TestFileOps_FileExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	// 创建文件
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := afero.WriteFile(fs, testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 测试文件存在
	if !ops.FileExists(testFile) {
		t.Error("FileExists() should return true for existing file")
	}

	// 测试文件不存在
	if ops.FileExists(filepath.Join(tmpDir, "nonexistent.txt")) {
		t.Error("FileExists() should return false for non-existent file")
	}
}

func TestFileOps_DirExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	// 创建目录
	fs.MkdirAll(tmpDir, 0755)

	// 测试目录存在
	if !ops.DirExists(tmpDir) {
		t.Error("DirExists() should return true for existing directory")
	}

	// 测试目录不存在
	if ops.DirExists(filepath.Join(tmpDir, "nonexistent")) {
		t.Error("DirExists() should return false for non-existent directory")
	}

	// 测试文件（不是目录）
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := afero.WriteFile(fs, testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if ops.DirExists(testFile) {
		t.Error("DirExists() should return false for files")
	}
}

func TestFileOps_GetFileSize(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	// 创建文件
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("test content")
	if err := afero.WriteFile(fs, testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 获取文件大小
	size := ops.GetFileSize(testFile)
	expectedSize := int64(len(testContent))

	if size != expectedSize {
		t.Errorf("GetFileSize() = %d, want %d", size, expectedSize)
	}
}

func TestFileOps_SetPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	// 使用真实文件系统进行权限测试
	fs := afero.NewOsFs()
	ops := comm.NewFileOps(fs)
	tmpDir := t.TempDir()

	// 创建文件
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := afero.WriteFile(fs, testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 设置权限
	ops.SetPermissions(testFile, "755")

	// 验证权限
	info, err := fs.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	expectedMode := uint32(0755)
	actualMode := uint32(info.Mode().Perm())
	if actualMode != expectedMode {
		t.Errorf("File mode = %o, want %o", actualMode, expectedMode)
	}
}

func TestFileOps_SetPermissionsEmpty(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := afero.WriteFile(fs, testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 空权限字符串应该不做任何事
	ops.SetPermissions(testFile, "")
}

func TestFileOps_CreateSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping symlink test on Windows")
	}

	// 符号链接需要使用真实文件系统
	fs := afero.NewOsFs()
	ops := comm.NewFileOps(fs)
	tmpDir := t.TempDir()

	// 创建目标文件
	targetFile := filepath.Join(tmpDir, "target.txt")
	if err := afero.WriteFile(fs, targetFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	// 创建符号链接
	linkFile := filepath.Join(tmpDir, "link.txt")
	ops.CreateSymlink(targetFile, linkFile, false, "")

	// 验证符号链接存在
	linkInfo, err := fs.Stat(linkFile)
	if err != nil {
		t.Fatalf("Failed to stat link: %v", err)
	}

	// 读取链接内容验证
	content, err := afero.ReadFile(fs, linkFile)
	if err != nil {
		t.Fatalf("Failed to read link: %v", err)
	}
	if string(content) != "test" {
		t.Errorf("Link content = %s, want test", content)
	}

	_ = linkInfo
}

func TestFileOps_CopyFilePanic(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	// 测试复制不存在的文件应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("CopyFile() with non-existent source should panic")
		}
	}()

	ops.CopyFile(filepath.Join(tmpDir, "nonexistent.txt"), filepath.Join(tmpDir, "dest.txt"))
}

func TestFileOps_GetFileSizePanic(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)

	// 测试获取不存在文件的大小应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetFileSize() with non-existent file should panic")
		}
	}()

	ops.GetFileSize("/nonexistent/file.txt")
}

func TestFileOps_RenameFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	// 创建测试文件
	oldPath := filepath.Join(tmpDir, "old.txt")
	if err := afero.WriteFile(fs, oldPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 重命名文件
	newPath := ops.RenameFile(oldPath, "new.txt")
	expectedPath := filepath.Join(tmpDir, "new.txt")

	if newPath != expectedPath {
		t.Errorf("RenameFile() = %s, want %s", newPath, expectedPath)
	}

	// 验证旧文件不存在
	if exists, _ := afero.Exists(fs, oldPath); exists {
		t.Error("Old file should not exist after rename")
	}

	// 验证新文件存在
	if exists, _ := afero.Exists(fs, newPath); !exists {
		t.Error("New file should exist after rename")
	}
}

func TestFileOps_CreateDirectory(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	testDir := filepath.Join(tmpDir, "testdir")
	ops.CreateDirectory(testDir, "", "", false, "")

	if !ops.DirExists(testDir) {
		t.Error("Directory should exist after CreateDirectory()")
	}
}

func TestFileOps_RemoveFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	ops := comm.NewFileOps(fs)
	tmpDir := "/tmp"

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := afero.WriteFile(fs, testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 删除文件
	ops.RemoveFile(testFile, false, "")

	// 验证文件不存在
	if exists, _ := afero.Exists(fs, testFile); exists {
		t.Error("File should not exist after RemoveFile()")
	}
}

func TestNewFileOpsWithNilFs(t *testing.T) {
	ops := comm.NewFileOps(nil)
	if ops.GetFs() == nil {
		t.Error("comm.NewFileOps(nil) should create OsFs")
	}
}
