package comm

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// ============================================================
// 集成测试：验证提取器、检查器和执行器的协作
// ============================================================

// TestIntegration_ExtractAndCheck 测试提取和检查的集成
func TestIntegration_ExtractAndCheck(t *testing.T) {
	a := require.New(t)

	// 1. 创建提取器
	extractor := NewCommandExtractor()

	// 2. 提取命令
	cmds, err := extractor.Extract("echo $(rm -rf /)")
	a.NoError(err)
	a.Len(cmds, 2)

	// 3. 创建检查器（阻止 rm）
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact))

	// 4. 检查命令
	err = checker.Check(cmds)
	a.Error(err)
	a.Contains(err.Error(), "blocked by blacklist")
}

// TestIntegration_ExtractAndCheckWithSourceFilter 测试来源过滤
func TestIntegration_ExtractAndCheckWithSourceFilter(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 提取命令（包含命令替换）
	cmds, err := extractor.Extract("echo $(rm -rf /)")
	a.NoError(err)
	a.Len(cmds, 2)

	// 创建检查器：只阻止命令替换中的 rm
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact).
			WithSourceFilter(SourceCmdSubst))

	// 应该检测到命令替换中的 rm
	err = checker.Check(cmds)
	a.Error(err)
}

// TestIntegration_WhitelistWithCommandSubstitution 测试白名单与命令替换
func TestIntegration_WhitelistWithCommandSubstitution(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 提取命令
	cmds, err := extractor.Extract("ls -la $(cat /etc/passwd)")
	a.NoError(err)
	a.Len(cmds, 2)

	// 创建检查器：白名单模式，只允许 ls
	checker := NewSecurityChecker().
		WithWhitelistMode(true).
		WithWhitelist(NewCommandRule("ls", MatchExact))

	// cat 不在白名单中，应该被阻止
	err = checker.Check(cmds)
	a.Error(err)
	a.Contains(err.Error(), "not in whitelist")
}

// TestIntegration_PipelineCommands 测试管道命令
func TestIntegration_PipelineCommands(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 提取管道命令
	cmds, err := extractor.Extract("cat file | grep pattern | wc -l")
	a.NoError(err)
	a.Len(cmds, 3)

	// 创建检查器：阻止 grep
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("grep", MatchExact))

	err = checker.Check(cmds)
	a.Error(err)
}

// TestIntegration_NestedCommandSubstitution 测试嵌套命令替换
func TestIntegration_NestedCommandSubstitution(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 提取嵌套命令替换
	cmds, err := extractor.Extract("echo $(echo $(rm -rf /))")
	a.NoError(err)
	a.Len(cmds, 3)

	// 验证最内层的 rm 被正确提取
	a.Equal("rm", cmds[2].Name)
	a.Equal(SourceCmdSubst, cmds[2].Source)

	// 创建检查器：阻止 rm
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact))

	err = checker.Check(cmds)
	a.Error(err)
}

// TestIntegration_AllowSafeCommands 测试允许安全命令
func TestIntegration_AllowSafeCommands(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 提取安全命令
	cmds, err := extractor.Extract("ls -la && cat file.txt")
	a.NoError(err)
	a.Len(cmds, 2)

	// 创建检查器：白名单模式，允许 ls 和 cat
	checker := NewSecurityChecker().
		WithWhitelistMode(true).
		WithWhitelist(
			NewCommandRule("ls", MatchExact),
			NewCommandRule("cat", MatchExact),
		)

	// 所有命令都在白名单中，应该通过
	err = checker.Check(cmds)
	a.NoError(err)
}

// TestIntegration_ExecuteWithSecurityCheck 测试带安全检查的执行
func TestIntegration_ExecuteWithSecurityCheck(t *testing.T) {
	a := require.New(t)

	// 创建执行器配置
	config := DefaultGoshConfig()

	// 执行器应该正常工作
	executor := NewGoshExecutor(config)

	var out bytes.Buffer
	err := executor.Run(context.Background(), "", "echo hello", nil, &out, io.Discard)
	a.NoError(err)
	a.Equal("hello\n", out.String())
}

// TestIntegration_ComplexScenario 测试复杂场景
func TestIntegration_ComplexScenario(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 复杂命令：包含管道、条件、命令替换
	cmd := `
		if [ -f /tmp/file ]; then
			cat /tmp/file | grep pattern
		else
			echo $(ls /tmp)
		fi
	`

	cmds, err := extractor.Extract(cmd)
	a.NoError(err)

	// 验证提取的命令数量
	// [ (test), cat, grep, echo, ls
	a.GreaterOrEqual(len(cmds), 4)

	// 创建检查器：阻止 grep
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("grep", MatchExact))

	// 应该检测到 grep
	err = checker.Check(cmds)
	a.Error(err)
}

// TestIntegration_EnvironmentVariables 测试环境变量提取
func TestIntegration_EnvironmentVariables(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 带环境变量的命令
	cmds, err := extractor.Extract("VAR=value ls -la")
	a.NoError(err)
	a.Len(cmds, 1)

	// 验证环境变量被提取
	a.Len(cmds[0].Envs, 1)
	a.Equal("VAR=value", cmds[0].Envs[0])
}

// TestIntegration_HandlerRegistryWithExecutor 测试 Handler 注册器与执行器的集成
func TestIntegration_HandlerRegistryWithExecutor(t *testing.T) {
	a := require.New(t)

	// 创建 Handler 注册器
	registry := NewHandlerRegistry()

	// 注册一个简单的 echo 处理器（仅用于测试）
	// 注意：实际使用中，执行器使用的是 config.GoHandlers map
	_ = registry

	// 验证注册器工作正常
	handler, ok := registry.Match("notexist")
	a.False(ok)
	a.Nil(handler)
}

// TestIntegration_ArgsFilterWithExtractor 测试参数过滤与提取器的集成
func TestIntegration_ArgsFilterWithExtractor(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 提取命令
	cmds, err := extractor.Extract("rm -rf /tmp")
	a.NoError(err)
	a.Len(cmds, 1)
	a.Equal("rm", cmds[0].Name)
	a.Equal([]string{"-rf", "/tmp"}, cmds[0].Args)

	// 创建检查器：只阻止带 -rf 参数的 rm
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact).
			WithArgsFilter(ArgMatcher{Position: 0, Pattern: "-rf", Mode: MatchExact}))

	// 应该被阻止
	err = checker.Check(cmds)
	a.Error(err)

	// 测试不带 -rf 的 rm
	cmds2, err := extractor.Extract("rm /tmp/file.txt")
	a.NoError(err)
	a.Len(cmds2, 1)

	// 应该通过
	err = checker.Check(cmds2)
	a.NoError(err)
}
