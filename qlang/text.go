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

// SliceLines 从文本中按行号和限制切分行
// 参数:
//   - text: 要切分的文本
//   - line: 起始行号（1-based），nil 表示从第一行开始
//   - limit: 最大行数，nil 表示不限制
//
// 返回: 切分后的文本
//
// 示例:
//
//	SliceLines("a\nb\nc", intPtr(2), nil) // 返回 "b\nc"
//	SliceLines("a\nb\nc", nil, intPtr(2)) // 返回 "a\nb"
func SliceLines(text string, line *int, limit *int) string {
	if line == nil && limit == nil {
		return text
	}

	lines := Text2Lines(text)

	start := 0
	if line != nil {
		start = *line - 1 // line 是 1-based
		if start < 0 {
			start = 0
		}
	}

	end := len(lines)
	if limit != nil {
		if start+*limit < end {
			end = start + *limit
		}
	}

	if start > len(lines) {
		start = len(lines)
	}

	return JoinLines(lines[start:end])
}

// SliceLinesP 是 SliceLines 的 panic 版本，使用整数而非指针
// 参数:
//   - text: 要切分的文本
//   - line: 起始行号（1-based），0 或负数表示从第一行开始
//   - limit: 最大行数，0 或负数表示不限制
//
// 返回: 切分后的文本
func SliceLinesP(text string, line int, limit int) string {
	var linePtr *int
	if line > 0 {
		linePtr = &line
	}

	var limitPtr *int
	if limit > 0 {
		limitPtr = &limit
	}

	return SliceLines(text, linePtr, limitPtr)
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

// Shorten 截断字符串到指定长度，如果超过则在末尾添加省略号
// 用于日志输出等场景，避免输出过长内容
//
// 参数:
//   - s: 要截断的字符串
//   - maxLen: 最大长度（包含省略号的长度）
//
// 返回: 截断后的字符串
//
// 示例:
//
//	Shorten("hello world", 8) // 返回 "hello..."
//	Shorten("hi", 10)         // 返回 "hi" (不超过长度，原样返回)
func Shorten(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return "..."
	}
	return s[:maxLen-3] + "..."
}

// ShortenWithSuffix 截断字符串到指定长度，使用自定义后缀
//
// 参数:
//   - s: 要截断的字符串
//   - maxLen: 最大长度
//   - suffix: 截断后添加的后缀（如 "..." 或 "[truncated]"）
//
// 返回: 截断后的字符串
func ShortenWithSuffix(s string, maxLen int, suffix string) string {
	if len(s) <= maxLen {
		return s
	}
	suffixLen := len(suffix)
	if maxLen <= suffixLen {
		return suffix
	}
	return s[:maxLen-suffixLen] + suffix
}
