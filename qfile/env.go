package qfile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/qiangyt/go-comm/v2"
	"github.com/qiangyt/go-comm/v2/qshell"
	"github.com/spf13/afero"
)

func SysEnvFileNames(fs afero.Fs, shell string) []string {
	r := []string{}

	if len(shell) == 0 {
		shell = os.Getenv("SHELL")
	}

	home, _ := ExpandHomePath("~")
	hasHome := (len(home) > 0)

	pth := filepath.Join("/etc/profile")
	if exists, _ := FileExists(fs, pth); exists {
		r = append(r, pth)
	}

	pth = filepath.Join("/etc/paths")
	if exists, _ := FileExists(fs, pth); exists {
		r = append(r, pth)
	}

	if !strings.Contains(shell, "zsh") {
		if exists, _ := FileExists(fs, "/etc/bashrc"); exists {
			r = append(r, pth)
		}

		if hasHome {
			pth = filepath.Join(home, ".bashrc")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			}

			pth = filepath.Join(home, ".bash_profile")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			} else {
				pth = filepath.Join(home, ".bash_login")
				if exists, _ := FileExists(fs, pth); exists {
					r = append(r, pth)
				}
				pth = filepath.Join(home, ".profile")
				if exists, _ := FileExists(fs, pth); exists {
					r = append(r, pth)
				}
			}
		}
	} else {
		if exists, _ := FileExists(fs, "/etc/zshrc"); exists {
			r = append(r, pth)
		}

		if hasHome {
			pth = filepath.Join(home, ".zshrc")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			}

			pth = filepath.Join(home, ".zshenv")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			}

			pth = filepath.Join(home, ".zprofile")
			if exists, _ := FileExists(fs, pth); exists {
				r = append(r, pth)
			} else {
				pth = filepath.Join(home, ".zsh_login")
				if exists, _ := FileExists(fs, pth); exists {
					r = append(r, pth)
				}
				pth = filepath.Join(home, ".profile")
				if exists, _ := FileExists(fs, pth); exists {
					r = append(r, pth)
				}
			}
		}
	}

	pth = ".env"
	if exists, _ := FileExists(fs, pth); exists {
		r = append(r, pth)
	}

	return r
}

func LoadEnvScripts(fs afero.Fs, vars map[string]string, filenames ...string) (map[string]string, error) {
	errs := comm.NewErrorGroup(false)

	if len(filenames) == 0 {
		filenames = SysEnvFileNames(fs, "")
	}

	for _, filename := range filenames {
		var err error
		vars, err = LoadEnvScript(fs, vars, filename)
		errs.Add(err)
	}

	return vars, errs.MayError()
}

func LoadEnvScript(fs afero.Fs, vars map[string]string, filename string) (map[string]string, error) {
	if filename == "/etc/paths" {
		paths, err := ReadFileLines(fs, filename)
		if err == nil {
			if len(vars["PATH"]) >= 0 {
				paths = append([]string{vars["PATH"]}, paths...)
			}
			vars["PATH"] = strings.Join(paths, ":")
		}
		return vars, err
	}

	output, err := qshell.RunGoshCommand(vars, "", filename, nil)
	if err != nil {
		return vars, err
	}
	return output.Vars, nil
}
