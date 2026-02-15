package comm

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"mvdan.cc/sh/v3/interp"
)

// ============================================================
// GoshConfig Tests
// ============================================================

func TestDefaultGoshConfig_happy(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig()
	a.Equal(6*time.Second, config.KillTimeout)
	a.False(config.WhitelistMode)
	a.Empty(config.Blacklist)
	a.Empty(config.Whitelist)
	a.Empty(config.GoHandlers)
}

func TestGoshConfig_WithKillTimeout(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig().WithKillTimeout(30 * time.Second)
	a.Equal(30*time.Second, config.KillTimeout)
}

func TestGoshConfig_WithBlacklist(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig().
		WithBlacklist(NewCommandRule("rm", MatchExact).WithArgsFilter(
			ArgMatcher{Position: -1, Pattern: "-rf", Mode: MatchGlob},
		)).
		WithBlacklistSimple("dd", "mkfs")

	a.Len(config.Blacklist, 3)
	a.Equal("rm", config.Blacklist[0].Pattern)
	a.Len(config.Blacklist[0].ArgsMatchers, 1)
	a.Equal("dd", config.Blacklist[1].Pattern)
	a.Len(config.Blacklist[1].ArgsMatchers, 0)
}

func TestGoshConfig_WithWhitelist(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig().
		WithWhitelistMode(true).
		WithWhitelistSimple("git", "npm", "go")

	a.True(config.WhitelistMode)
	a.Len(config.Whitelist, 3)
}

func TestGoshConfig_WithGoHandler(t *testing.T) {
	a := require.New(t)

	handler := func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		return nil
	}

	config := DefaultGoshConfig().
		WithGoHandler("curl", handler).
		WithGoHandler("wget*", handler)

	a.Len(config.GoHandlers, 2)
}

// ============================================================
// Pattern Matching Tests
// ============================================================

func TestMatchPattern_exact(t *testing.T) {
	a := require.New(t)
	e := NewGoshExecutor(nil)

	a.True(e.matchPattern("rm", "rm"))
	a.False(e.matchPattern("rm", "rmdir"))
}

func TestMatchPattern_wildcard(t *testing.T) {
	a := require.New(t)
	e := NewGoshExecutor(nil)

	a.True(e.matchPattern("rm*", "rm"))
	a.True(e.matchPattern("rm*", "rmdir"))
	a.True(e.matchPattern("rm*", "rmmod"))
	a.False(e.matchPattern("rm*", "ls"))

	a.True(e.matchPattern("curl*", "curl"))
	a.True(e.matchPattern("curl*", "curlconfig"))
}

func TestMatchPattern_question(t *testing.T) {
	a := require.New(t)
	e := NewGoshExecutor(nil)

	// ? 匹配单个字符
	a.False(e.matchPattern("rm?", "rmdir")) // rm? 不匹配 rmdir（5字符，rm?是3字符）
	a.False(e.matchPattern("rm?", "rm"))    // rm? 不匹配 rm（2字符，rm?是3字符）
	a.True(e.matchPattern("rm?", "rmx"))    // rm? 匹配 rmx（3字符）
}

// ============================================================
// Args Matching Tests (使用 SecurityChecker 测试)
// ============================================================

func TestMatchArgs_happy(t *testing.T) {
	a := require.New(t)

	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact).
			WithArgsFilter(ArgMatcher{Position: -1, Pattern: "-rf", Mode: MatchGlob}))

	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/tmp"}},
	}
	err := checker.Check(cmds)
	a.Error(err) // 应该被阻止
}

func TestMatchArgs_noMatch(t *testing.T) {
	a := require.New(t)

	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact).
			WithArgsFilter(ArgMatcher{Position: -1, Pattern: "-rf", Mode: MatchGlob}))

	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-i", "file.txt"}},
	}
	err := checker.Check(cmds)
	a.NoError(err) // 不应该被阻止
}

// ============================================================
// Executor Tests
// ============================================================

func TestGoshExecutor_Run_simpleCommand(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig()
	executor := NewGoshExecutor(config)

	var out bytes.Buffer
	err := executor.Run(context.Background(), "", "echo hello", nil, &out, &out)
	a.NoError(err)
	a.Equal("hello\n", out.String())
}

func TestGoshExecutor_Run_blacklistedCommand(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig().WithBlacklistSimple("rm")
	executor := NewGoshExecutor(config)

	var out bytes.Buffer
	err := executor.Run(context.Background(), "", "rm test.txt", nil, &out, &out)
	a.Error(err)
	a.Contains(err.Error(), "blocked by blacklist")
}

func TestGoshExecutor_Run_blacklistWithWildcard(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig().WithBlacklist(NewCommandRule("rm*", MatchGlob))
	executor := NewGoshExecutor(config)

	var out bytes.Buffer

	// rm 应该被阻止
	err := executor.Run(context.Background(), "", "rm test.txt", nil, &out, &out)
	a.Error(err)

	// rmdir 也应该被阻止（因为 rm* 匹配）
	err = executor.Run(context.Background(), "", "rmdir testdir", nil, &out, &out)
	a.Error(err)
}

func TestGoshExecutor_Run_blacklistWithArgs(t *testing.T) {
	a := require.New(t)

	// 只阻止 rm -rf，允许普通 rm
	config := DefaultGoshConfig().WithBlacklist(
		NewCommandRule("rm", MatchExact).WithArgsFilter(
			ArgMatcher{Position: -1, Pattern: "-r*", Mode: MatchGlob},
		),
	)
	executor := NewGoshExecutor(config)

	var out bytes.Buffer

	// rm -rf 应该被阻止
	err := executor.Run(context.Background(), "", "rm -rf test", nil, &out, &out)
	a.Error(err)
	a.Contains(err.Error(), "blocked by blacklist")

	// rm -r 也应该被阻止（因为 -r 匹配 -r*）
	err = executor.Run(context.Background(), "", "rm -r testdir", nil, &out, &out)
	a.Error(err)
	a.Contains(err.Error(), "blocked by blacklist")

	// 普通 rm 不应该被黑名单阻止（但因为文件不存在会失败，所以我们用 --help 测试）
	// 使用 rm --help 来测试命令不被阻止（--help 不匹配 -r*）
	err = executor.Run(context.Background(), "", "rm --help", nil, &out, &out)
	// 命令应该执行（不被黑名单阻止），但可能返回非零退出码
	// 关键是检查错误信息不包含 "blocked by blacklist"
	if err != nil {
		a.NotContains(err.Error(), "blocked by blacklist")
	}
}

func TestGoshExecutor_Run_whitelistMode(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig().
		WithWhitelistMode(true).
		WithWhitelistSimple("echo", "ls")
	executor := NewGoshExecutor(config)

	var out bytes.Buffer

	// 白名单中的命令可以执行
	err := executor.Run(context.Background(), "", "echo hello", nil, &out, &out)
	a.NoError(err)

	// 不在白名单中的命令被拒绝
	err = executor.Run(context.Background(), "", "rm test.txt", nil, &out, &out)
	a.Error(err)
	a.Contains(err.Error(), "not in whitelist")
}

func TestGoshExecutor_Run_goHandler(t *testing.T) {
	a := require.New(t)

	called := false
	handler := func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		called = true
		fmt.Fprintln(hc.Stdout, "handled by go")
		return nil
	}

	config := DefaultGoshConfig().WithGoHandler("curl", handler)
	executor := NewGoshExecutor(config)

	var out bytes.Buffer
	err := executor.Run(context.Background(), "", "curl http://example.com", nil, &out, &out)
	a.NoError(err)
	a.True(called)
	a.Contains(out.String(), "handled by go")
}

func TestGoshExecutor_Run_goHandlerWithWildcard(t *testing.T) {
	a := require.New(t)

	called := false
	handler := func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		called = true
		return nil
	}

	config := DefaultGoshConfig().WithGoHandler("wget*", handler)
	executor := NewGoshExecutor(config)

	var out bytes.Buffer

	// wget 应该被处理器拦截
	err := executor.Run(context.Background(), "", "wget http://example.com", nil, &out, &out)
	a.NoError(err)
	a.True(called)

	// wget2 也应该被拦截（因为 wget* 匹配）
	called = false
	err = executor.Run(context.Background(), "", "wget2 http://example.com", nil, &out, &out)
	a.NoError(err)
	a.True(called)
}

func TestGoshExecutor_Run_withVars(t *testing.T) {
	a := require.New(t)

	vars := map[string]string{"NAME": "world"}
	config := DefaultGoshConfig()
	executor := NewGoshExecutor(config)

	var out bytes.Buffer
	err := executor.RunWithVars(context.Background(), vars, "", "echo hello ${NAME}", nil, &out, &out)
	a.NoError(err)
	a.Equal("hello world\n", out.String())
}

func TestGoshExecutor_Run_pipeline(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig()
	executor := NewGoshExecutor(config)

	var out bytes.Buffer
	err := executor.Run(context.Background(), "", "echo hello | cat", nil, &out, &out)
	a.NoError(err)
	a.Equal("hello\n", out.String())
}

func TestGoshExecutor_Run_andOperator(t *testing.T) {
	a := require.New(t)

	config := DefaultGoshConfig()
	executor := NewGoshExecutor(config)

	var out bytes.Buffer
	err := executor.Run(context.Background(), "", "echo first && echo second", nil, &out, &out)
	a.NoError(err)
	a.Contains(out.String(), "first")
	a.Contains(out.String(), "second")
}
