package qio

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v3/q18n"
	"github.com/qiangyt/go-comm/v3/qerr"
	"github.com/qiangyt/go-comm/v3/qlang"
	"github.com/spf13/afero"
)

var AppFs = afero.NewOsFs()

func DefaultEtcHostsP() string {
	r, err := DefaultEtcHosts()
	if err != nil {
		panic(qerr.NewSystemError("get default etc hosts", err))
	}
	return r
}

func CopyFileP(fs afero.Fs, path string, newPath string) int64 {
	r, err := CopyFile(fs, path, newPath)
	if err != nil {
		panic(qerr.NewSystemError("copy file", err))
	}
	return r
}

func CopyFile(fs afero.Fs, path string, newPath string) (int64, error) {
	err := EnsureFileExists(fs, path)
	if err != nil {
		return 0, err
	}

	src, err := fs.Open(path)
	if err != nil {
		return 0, errors.Wrapf(err, "read file %s", path)
	}
	defer src.Close()

	dst, err := fs.Create(newPath)
	if err != nil {
		return 0, errors.Wrapf(err, "create file %s", newPath)
	}
	defer dst.Close()

	nBytes, err := io.Copy(dst, src)
	if err != nil {
		return 0, errors.Wrapf(err, "copy file %s to %s", path, newPath)
	}
	return nBytes, nil
}

func RenameP(fs afero.Fs, path string, newPath string) {
	if err := Rename(fs, path, newPath); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func Rename(fs afero.Fs, path string, newPath string) error {
	err := fs.Rename(path, newPath)
	if err != nil {
		return errors.Wrapf(err, "move file %s to %s", path, newPath)
	}
	return nil
}

func StatP(fs afero.Fs, path string, ensureExists bool) os.FileInfo {
	r, err := Stat(fs, path, ensureExists)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

// Stat ...
func Stat(fs afero.Fs, path string, ensureExists bool) (os.FileInfo, error) {
	r, err := fs.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "stat file: %s", path)
		}
		if ensureExists {
			return nil, errors.Wrap(err, q18n.T("error.file.not_found", map[string]any{
				"Path": path,
			}))
		}
		return nil, nil
	}

	return r, nil
}

func FileExistsP(fs afero.Fs, path string) bool {
	r, err := FileExists(fs, path)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

// FileExists ...
func FileExists(fs afero.Fs, path string) (bool, error) {
	fi, err := Stat(fs, path, false)
	if err != nil {
		return false, err
	}
	if fi == nil {
		return false, nil
	}
	if fi.IsDir() {
		return false, q18n.LocalizeError("error.file.expect_file_but_dir", map[string]any{
			"Path": path,
		})
	}
	return true, nil
}

func EnsureFileExistsP(fs afero.Fs, path string) {
	if err := EnsureFileExists(fs, path); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func EnsureFileExists(fs afero.Fs, path string) error {
	exists, err := FileExists(fs, path)
	if err != nil {
		return err
	}
	if !exists {
		return q18n.LocalizeError("error.file.not_found", map[string]any{
			"Path": path,
		})
	}
	return nil
}

func EnsureFileNotExistsP(fs afero.Fs, path string) {
	if err := EnsureFileNotExists(fs, path); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func EnsureFileNotExists(fs afero.Fs, path string) error {
	exists, err := FileExists(fs, path)
	if err != nil {
		return err
	}
	if exists {
		return q18n.LocalizeError("error.file.already_exists", map[string]any{
			"Path": path,
		})
	}
	return nil
}

func DirExistsP(fs afero.Fs, path string) bool {
	r, err := DirExists(fs, path)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

// DirExists ...
func DirExists(fs afero.Fs, path string) (bool, error) {
	fi, err := Stat(fs, path, false)
	if err != nil {
		return false, err
	}
	if fi == nil {
		return false, nil
	}
	if !fi.IsDir() {
		return false, q18n.LocalizeError("error.file.expect_dir_but_file", map[string]any{
			"Path": path,
		})
	}
	return true, nil
}

func EnsureDirExistsP(fs afero.Fs, path string) {
	if err := EnsureDirExists(fs, path); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func EnsureDirExists(fs afero.Fs, path string) error {
	exists, err := DirExists(fs, path)
	if err != nil {
		return err
	}
	if !exists {
		return q18n.LocalizeError("error.dir.not_found", map[string]any{
			"Path": path,
		})
	}
	return nil
}

func EnsureDirNotExistsP(fs afero.Fs, path string) {
	if err := EnsureDirNotExists(fs, path); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func EnsureDirNotExists(fs afero.Fs, path string) error {
	exists, err := DirExists(fs, path)
	if err != nil {
		return err
	}
	if exists {
		return q18n.LocalizeError("error.dir.already_exists", map[string]any{
			"Path": path,
		})
	}
	return nil
}

func RemoveFileP(fs afero.Fs, path string) {
	if err := RemoveFile(fs, path); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

// RemoveFile ...
func RemoveFile(fs afero.Fs, path string) error {
	found, err := FileExists(fs, path)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	if err := fs.Remove(path); err != nil {
		return errors.Wrapf(err, "delete file: %s", path)
	}
	return nil
}

func RemoveDirP(fs afero.Fs, path string) {
	if err := RemoveDir(fs, path); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

// RemoveDir ...
func RemoveDir(fs afero.Fs, path string) error {
	if path == "/" || path == "\\" {
		return q18n.LocalizeError("error.file.cannot_remove_root", nil)
	}
	found, err := DirExists(fs, path)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	if err := fs.RemoveAll(path); err != nil {
		return errors.Wrapf(err, "delete directory: %s", path)
	}
	return nil
}

func ReadFileBytesP(fs afero.Fs, path string) []byte {
	r, err := ReadFileBytes(fs, path)
	if err != nil {
		return r
	}
	return r
}

// ReadBytes ...
func ReadFileBytes(fs afero.Fs, path string) ([]byte, error) {
	r, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrapf(err, "read file: %s", path)
	}
	return r, nil
}

func ReadFileTextP(fs afero.Fs, path string) string {
	r, err := ReadFileText(fs, path)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func ReadFileText(fs afero.Fs, path string) (string, error) {
	bytes, err := ReadFileBytes(fs, path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ReadFileLinesP(fs afero.Fs, path string) []string {
	r, err := ReadFileLines(fs, path)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func ReadFileLines(fs afero.Fs, path string) ([]string, error) {
	if err := EnsureFileExists(fs, path); err != nil {
		return nil, err
	}

	f, err := fs.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "open file: %s", path)
	}
	defer f.Close()

	return ReadLines(f), nil
}

func WriteFileIfNotFoundP(fs afero.Fs, path string, content []byte) bool {
	r, err := WriteFileIfNotFound(fs, path, content)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

// WriteIfNotFound ...
func WriteFileIfNotFound(fs afero.Fs, path string, content []byte) (bool, error) {
	found, err := FileExists(fs, path)
	if err != nil {
		return false, err
	}
	if found {
		return false, nil
	}
	if err := WriteFile(fs, path, content); err != nil {
		return true, err
	}
	return true, nil
}

func WriteFileP(fs afero.Fs, path string, content []byte) {
	if err := WriteFile(fs, path, content); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func Mkdir4FileP(fs afero.Fs, filePath string) {
	if err := Mkdir4File(fs, filePath); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func Mkdir4File(fs afero.Fs, filePath string) error {
	dirPath := filepath.Dir(filePath)
	exists, err := DirExists(fs, dirPath)
	if err != nil {
		return nil
	}
	if !exists {
		if err = Mkdir(fs, dirPath); err != nil {
			return err
		}
	}
	return nil
}

// Write ...
func WriteFile(fs afero.Fs, path string, content []byte) error {
	if err := Mkdir4File(fs, path); err != nil {
		return err
	}

	if err := afero.WriteFile(fs, path, content, 0o640); err != nil {
		return errors.Wrapf(err, "write file: %s", path)
	}
	return nil
}

func WriteFileTextP(fs afero.Fs, path string, content string) {
	if err := WriteFileText(fs, path, content); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

// WriteText ...
func WriteFileText(fs afero.Fs, path string, content string) error {
	return WriteFile(fs, path, []byte(content))
}

// WriteTextIfNotFound ...
func WriteFileTextIfNotFoundP(fs afero.Fs, path string, content string) bool {
	r, err := WriteFileTextIfNotFound(fs, path, content)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func WriteFileTextIfNotFound(fs afero.Fs, path string, content string) (bool, error) {
	found, err := FileExists(fs, path)
	if err != nil {
		return found, err
	}
	if found {
		return false, nil
	}
	if err := WriteFileText(fs, path, content); err != nil {
		return false, err
	}
	return true, nil
}

// WriteLines ...
func WriteFileLinesP(fs afero.Fs, path string, lines ...string) {
	if err := WriteFileLines(fs, path, lines...); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func WriteFileLines(fs afero.Fs, path string, lines ...string) error {
	return WriteFileText(fs, path, qlang.JoinedLines(lines...))
}

func MkdirP(fs afero.Fs, path string) {
	if err := Mkdir(fs, path); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

// Mkdir ...
func Mkdir(fs afero.Fs, path string) error {
	if err := fs.MkdirAll(path, os.ModePerm); err != nil {
		return errors.Wrapf(err, "create directory: %s", path)
	}
	return nil
}

func ListSuffixedFilesP(fs afero.Fs, targetDir string, suffix string, skipEmptyFile bool) map[string]string {
	r, err := ListSuffixedFiles(fs, targetDir, suffix, skipEmptyFile)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func ListSuffixedFiles(fs afero.Fs, targetDir string, suffix string, skipEmptyFile bool) (map[string]string, error) {
	fiList, err := afero.ReadDir(fs, targetDir)
	if err != nil {
		return nil, errors.Wrapf(err, "read directory: %s", targetDir)
	}

	extLen := len(suffix)

	r := map[string]string{}
	for _, fi := range fiList {
		if fi.IsDir() {
			continue
		}
		if skipEmptyFile && fi.Size() == 0 {
			continue
		}

		fBase := filepath.Base(fi.Name())
		if !strings.HasSuffix(fBase, suffix) {
			continue
		}
		if len(fBase) == extLen {
			continue
		}

		fTitle := fBase[:len(fBase)-extLen]
		r[fTitle] = filepath.Join(targetDir, fi.Name())
	}

	return r, nil
}

func ExtractTitle(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	return base[:len(base)-len(ext)]
}

func TempFileP(fs afero.Fs, pattern string) string {
	r, err := TempFile(fs, pattern)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func TempFile(fs afero.Fs, pattern string) (string, error) {
	f, err := afero.TempFile(fs, "", pattern)
	if err != nil {
		return "", errors.Wrap(err, "create temporary file")
	}
	r := f.Name()
	f.Close()

	return r, nil
}

func TempTextFileP(fs afero.Fs, pattern string, content string) string {
	r, err := TempTextFile(fs, pattern, content)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func TempTextFile(fs afero.Fs, pattern string, content string) (string, error) {
	r, err := TempFile(fs, pattern)
	if err != nil {
		return "", err
	}
	if err := WriteFileText(fs, r, content); err != nil {
		return "", err
	}
	return r, nil
}

// EnsureDirWithSubdirsP 确保目录及其子目录存在（panic 版本）
func EnsureDirWithSubdirsP(fs afero.Fs, mainDir string, subdirs ...string) {
	if err := EnsureDirWithSubdirs(fs, mainDir, subdirs...); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

// EnsureDirWithSubdirs 确保目录及其子目录存在
func EnsureDirWithSubdirs(fs afero.Fs, mainDir string, subdirs ...string) error {
	// 创建主目录
	if err := Mkdir(fs, mainDir); err != nil {
		return err
	}
	// 创建子目录
	for _, subdir := range subdirs {
		fullPath := filepath.Join(mainDir, subdir)
		if err := Mkdir(fs, fullPath); err != nil {
			return err
		}
	}
	return nil
}

// EnsureFileWithContentP 确保文件存在，如果不存在则使用指定内容创建（panic 版本）
// 返回 true 表示文件是新创建的，false 表示文件已存在
func EnsureFileWithContentP(fs afero.Fs, path string, defaultContent []byte, subdirs ...string) bool {
	r, err := EnsureFileWithContent(fs, path, defaultContent, subdirs...)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

// EnsureFileWithContent 确保文件存在，如果不存在则使用指定内容创建
// 返回 true 表示文件是新创建的，false 表示文件已存在
// subdirs: 需要在文件父目录下创建的子目录列表
func EnsureFileWithContent(fs afero.Fs, path string, defaultContent []byte, subdirs ...string) (bool, error) {
	// 检查文件是否存在
	exists, err := FileExists(fs, path)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	// 确保父目录及其子目录存在
	parentDir := filepath.Dir(path)
	if err := EnsureDirWithSubdirs(fs, parentDir, subdirs...); err != nil {
		return false, err
	}

	// 写入默认内容
	if err := WriteFile(fs, path, defaultContent); err != nil {
		return false, err
	}
	return true, nil
}
