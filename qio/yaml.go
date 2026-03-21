package qio

import (
	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v3/qconfig"
	"github.com/qiangyt/go-comm/v3/qerr"
	"github.com/spf13/afero"
)

func FromYamlFileP(fs afero.Fs, path string, envsubt bool, result any) {
	if err := FromYamlFile(fs, path, envsubt, result); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func FromYamlFile(fs afero.Fs, path string, envsubt bool, result any) error {
	yamlText, err := ReadFileText(fs, path)
	if err != nil {
		return err
	}

	if err := qconfig.FromYaml(yamlText, envsubt, result); err != nil {
		return errors.Wrapf(err, "parse yaml file: %s", path)
	}
	return nil
}

func MapFromYamlFileP(fs afero.Fs, path string, envsubt bool) map[string]any {
	r, err := MapFromYamlFile(fs, path, envsubt)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func MapFromYamlFile(fs afero.Fs, path string, envsubt bool) (map[string]any, error) {
	r := map[string]any{}
	if err := FromYamlFile(fs, path, envsubt, &r); err != nil {
		return nil, err
	}

	return r, nil
}
