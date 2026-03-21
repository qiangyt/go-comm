package qsys

import (
	"os"
	"strings"

	"github.com/a8m/envsubst/parse"
	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v3/qerr"
	"github.com/spf13/cast"
)

func EnvironMapP(overrides map[string]string) map[string]string {
	r, err := EnvironMap(overrides)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func EnvironMap(overrides map[string]string) (map[string]string, error) {
	envs := strings.Join(os.Environ(), "\n")
	r, err := UnmarshalEnv(envs)
	if err != nil {
		return nil, errors.Wrapf(err, "parse OS environments")
	}

	if len(overrides) > 0 {
		for k, v := range overrides {
			r[k] = cast.ToString(v)
		}
	}
	return r, nil
}

func EnvironListP(overrides map[string]string) []string {
	r, err := EnvironList(overrides)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func EnvironList(overrides map[string]string) ([]string, error) {
	envs, err := EnvironMap(overrides)
	if err != nil {
		return nil, err
	}

	r := make([]string, 0, len(envs)+len(overrides))
	for k, v := range envs {
		r = append(r, k+"="+cast.ToString(v))
	}
	return r, nil
}

func EnvSubstP(input string, env map[string]string) string {
	r, err := EnvSubst(input, env)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func EnvSubst(input string, env map[string]string) (string, error) {
	restr := parse.Restrictions{NoUnset: false, NoEmpty: false}

	envMap, err := EnvironMap(env)
	if err != nil {
		return "", err
	}
	envList, err := EnvironList(envMap)
	if err != nil {
		return "", err
	}

	parser := parse.New("tmp", envList, &restr)
	r, err := parser.Parse(input)
	if err != nil {
		return "", errors.Wrapf(err, "envsubst the text: %s", input)
	}
	return r, nil
}

func EnvSubstSliceP(inputs []string, env map[string]string) []string {
	r, err := EnvSubstSlice(inputs, env)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func EnvSubstSlice(inputs []string, env map[string]string) ([]string, error) {
	r := make([]string, 0, len(inputs))
	for _, s := range inputs {
		substed, err := EnvSubst(s, env)
		if err != nil {
			return nil, err
		}
		r = append(r, substed)
	}
	return r, nil
}
