package qlang

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v3/qerr"
	"gopkg.in/yaml.v3"
)

func RenderWithTemplateP(w io.Writer, name string, tmpl string, data map[string]any) {
	err := RenderWithTemplate(w, name, tmpl, data)
	if err != nil {
		panic(qerr.NewConfigError(err.Error(), err))
	}
}

func RenderWithTemplate(w io.Writer, name string, tmpl string, data map[string]any) error {
	t, err := template.New(name).Funcs(sprig.FuncMap()).Parse(tmpl)
	if err != nil {
		return errors.Wrapf(err, "parse template %s: %s", name, tmpl)
	}

	err = t.Execute(w, data)
	if err != nil {
		return errors.Wrapf(err, "render template %s: %v", name, data)
	}
	return nil
}

func RenderAsTemplateArrayP(tmplArray []string, data map[string]any) []string {
	r, err := RenderAsTemplateArray(tmplArray, data)
	if err != nil {
		panic(qerr.NewConfigError(err.Error(), err))
	}
	return r
}

func RenderAsTemplateArray(tmplArray []string, data map[string]any) ([]string, error) {
	r := make([]string, 0, len(tmplArray))
	for _, tmpl := range tmplArray {
		txt, err := RenderAsTemplate(tmpl, data)
		if err != nil {
			return nil, err
		}
		r = append(r, txt)
	}
	return r, nil
}

func RenderAsTemplateP(tmpl string, data map[string]any) string {
	r, err := RenderAsTemplate(tmpl, data)
	if err != nil {
		panic(qerr.NewConfigError(err.Error(), err))
	}
	return r
}

func RenderAsTemplate(tmpl string, data map[string]any) (string, error) {
	output := &strings.Builder{}
	if err := RenderWithTemplate(output, "", tmpl, data); err != nil {
		return "", err
	}
	return output.String(), nil
}

func JoinedLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func JoinedLinesAsBytes(lines ...string) []byte {
	return []byte(JoinedLines(lines...))
}

func ToYamlP(hint string, me any) string {
	r, err := ToYaml(hint, me)
	if err != nil {
		panic(qerr.NewConfigError(err.Error(), err))
	}
	return r
}

func ToYaml(hint string, me any) (string, error) {
	r, err := yaml.Marshal(me)
	if err != nil {
		if len(hint) > 0 {
			return "", errors.Wrapf(err, "marshal %s to yaml", hint)
		} else {
			return "", errors.Wrapf(err, "marshal to yaml")
		}
	}
	return string(r), nil
}

func SubstVarsP(useGoTemplate bool, m map[string]any, parentVars map[string]any, keysToSkip ...string) map[string]any {
	r, err := SubstVars(useGoTemplate, m, parentVars, keysToSkip...)
	if err != nil {
		panic(qerr.NewConfigError(err.Error(), err))
	}
	return r
}

func SubstVars(useGoTemplate bool, m map[string]any, parentVars map[string]any, keysToSkip ...string) (map[string]any, error) {
	newVars := map[string]any{}

	// copy parent vars, it could be overwritten by local vars
	if len(parentVars) > 0 {
		for k, v := range parentVars {
			newVars[k] = v
		}
	}

	for k, v := range m {
		if k == "vars" {
			if localVarsMap, isMap := v.(map[string]any); isMap {
				// overwrite by local vars
				for k2, v2 := range localVarsMap {
					newVars[k2] = v2
				}
			}
		}
	}

	mapNoVars := map[string]any{}
	for k, v := range m {
		if k != "vars" {
			skip := false
			if len(keysToSkip) > 0 {
				for _, keyToSkip := range keysToSkip {
					if keyToSkip == k {
						skip = true
					}
				}
			}
			if !skip {
				/*vYaml := ToYaml("", v)
				vYaml = RenderAsTemplate(vYaml, newVars)
				if err := yaml.Unmarshal([]byte(vYaml), &v); err != nil {
					panic(qerr.NewConfigError("parse yaml", err))
				}*/
				mapNoVars[k] = v
			}
		}
	}

	yamlNoVars, err := ToYaml("", mapNoVars)
	if err != nil {
		return nil, err
	}

	if useGoTemplate {
		yamlNoVars, err = RenderAsTemplate(yamlNoVars, newVars)
		if err != nil {
			return nil, err
		}
	} else {
		yamlNoVars = os.Expand(yamlNoVars,
			func(k string) string {
				v := newVars[k]
				if v != nil {
					if vStr, isStr := v.(string); isStr {
						return vStr
					}
					return fmt.Sprintf("%v", v)
				}
				return ""
			})
	}

	r := map[string]any{}
	if err := yaml.Unmarshal([]byte(yamlNoVars), &r); err != nil {
		return nil, errors.Wrapf(err, "parse yaml: %s", yamlNoVars)
	}
	r["vars"] = newVars

	// put back skipped key/values
	if len(keysToSkip) > 0 {
		for _, keyToSkip := range keysToSkip {
			r[keyToSkip] = m[keyToSkip]
		}
	}

	return r, nil
}

func TextLine2Array(line string) []string {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return []string{}
	}

	var r []string
	if strings.ContainsAny(line, ",") {
		r = strings.Split(line, ",")
	} else if strings.ContainsAny(line, "\t") {
		r = strings.Split(line, "\t")
	} else if strings.ContainsAny(line, "\n") {
		r = strings.Split(line, "\n")
	} else if strings.ContainsAny(line, "\r") {
		r = strings.Split(line, "\r")
	} else if strings.ContainsAny(line, ";") {
		r = strings.Split(line, ";")
	} else if strings.ContainsAny(line, "|") {
		r = strings.Split(line, "|")
	} else {
		r = strings.Split(line, " ")
	}

	for i, t := range r {
		r[i] = strings.TrimSpace(t)
	}
	return r
}

func Text2Lines(text string) []string {
	rdr := strings.NewReader(text)
	return ReadLines(rdr)
}

func ReadLines(reader io.Reader) []string {
	r := make([]string, 0, 32)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		r = append(r, line)
	}

	return r
}

// JoinLines 将行数组合并为文本，使用换行符连接
func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

// BlockedCommands 定义被阻止的命令列表
var BlockedCommands = map[string]bool{
	"rm":         true, // 删除文件
	"mkfs":       true, // 格式化
	"dd":         true, // 磁盘操作
	"fdisk":      true, // 分区
	"shred":      true, // 安全删除
	"partprobe":  true, // 通知内核分区变化
	"sfdisk":     true, // 分区操作
	"cfdisk":     true, // 分区工具
	"cryptsetup": true, // 加密设备
	"losetup":    true, // 循环设备
	"init":       true, // 系统初始化
	"shutdown":   true, // 关机
	"reboot":     true, // 重启
	"halt":       true, // 停止
	"poweroff":   true, // 关机
}

// CheckTerminalCommand 检查终端命令是否允许执行
// 使用 panic 替代 error 返回，符合项目 Error 处理规范
func CheckTerminalCommand(command string, args []string) {
	// 检查命令是否在黑名单中
	if BlockedCommands[command] {
		panic(qerr.NewSecurityErrorf("命令 '%s' 被安全策略阻止", command))
	}

	// 检查危险的 chmod -R 操作
	if command == "chmod" {
		for _, arg := range args {
			if arg == "-R" || arg == "--recursive" {
				panic(qerr.NewSecurityErrorf("递归 chmod 被安全策略阻止"))
			}
		}
	}
}
