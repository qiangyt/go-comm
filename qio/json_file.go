package qio

import (
	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v3/qerr"
	"github.com/qiangyt/go-comm/v3/qjson"
	"github.com/spf13/afero"
)

func FromJsonFileP(fs afero.Fs, path string, envsubt bool, result any) {
	if err := FromJsonFile(fs, path, envsubt, result); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func FromJsonFile(fs afero.Fs, path string, envsubt bool, result any) error {
	yamlText, err := ReadFileText(fs, path)
	if err != nil {
		return err
	}

	if err := qjson.FromJson(yamlText, envsubt, result); err != nil {
		return errors.Wrapf(err, "parse json file: %s", path)
	}
	return nil
}

func MapFromJsonFileP(fs afero.Fs, path string, envsubt bool) map[string]any {
	r, err := MapFromJsonFile(fs, path, envsubt)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func MapFromJsonFile(fs afero.Fs, path string, envsubt bool) (map[string]any, error) {
	r := map[string]any{}
	if err := FromJsonFile(fs, path, envsubt, &r); err != nil {
		return nil, err
	}

	return r, nil
}
