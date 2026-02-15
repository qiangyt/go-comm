package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ============================================================
// TestSecurityChecker 黑名单测试
// ============================================================

func TestSecurityChecker_BlacklistSimpleCommand(t *testing.T) {
	a := require.New(t)

	// 简单黑名单：阻止 rm 命令
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact))

	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}},
	}
	err := checker.Check(cmds)
	a.Error(err)
	a.Contains(err.Error(), "blocked by blacklist")
}

func TestSecurityChecker_BlacklistMultipleCommands(t *testing.T) {
	a := require.New(t)

	// 黑名单：阻止多个命令
	checker := NewSecurityChecker().
		WithBlacklist(
			NewCommandRule("rm", MatchExact),
			NewCommandRule("dd", MatchExact),
			NewCommandRule("mkfs", MatchExact),
		)

	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}},
	}
	err := checker.Check(cmds)
	a.Error(err)

	cmds = []ExtractedCommand{
		{Name: "dd", Args: []string{"if=/dev/zero"}},
	}
	err = checker.Check(cmds)
	a.Error(err)

	cmds = []ExtractedCommand{
		{Name: "ls", Args: []string{"-la"}},
	}
	err = checker.Check(cmds)
	a.NoError(err) // ls 不在黑名单中
}

func TestSecurityChecker_BlacklistWithGlob(t *testing.T) {
	a := require.New(t)

	// 通配符黑名单：阻止所有以 rm 开头的命令
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm*", MatchGlob))

	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}},
	}
	err := checker.Check(cmds)
	a.Error(err)

	cmds = []ExtractedCommand{
		{Name: "rmdir", Args: []string{"/tmp/empty"}},
	}
	err = checker.Check(cmds)
	a.Error(err) // rmdir 也匹配 rm*

	cmds = []ExtractedCommand{
		{Name: "ls", Args: []string{"-la"}},
	}
	err = checker.Check(cmds)
	a.NoError(err)
}

func TestSecurityChecker_BlacklistWithRegex(t *testing.T) {
	a := require.New(t)

	// 正则表达式黑名单：阻止所有以 rm 开头的命令
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule(`^rm.*`, MatchRegex))

	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{}},
	}
	err := checker.Check(cmds)
	a.Error(err)

	cmds = []ExtractedCommand{
		{Name: "rmdir", Args: []string{}},
	}
	err = checker.Check(cmds)
	a.Error(err)

	cmds = []ExtractedCommand{
		{Name: "lsrm", Args: []string{}}, // 不以 rm 开头
	}
	err = checker.Check(cmds)
	a.NoError(err)
}

// ============================================================
// TestSecurityChecker 参数级别黑名单测试
// ============================================================

func TestSecurityChecker_BlacklistWithArgsFilter(t *testing.T) {
	a := require.New(t)

	// 参数级别黑名单：只阻止 rm -rf
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact).
			WithArgsFilter(ArgMatcher{Position: 0, Pattern: "-rf", Mode: MatchExact}))

	// rm -rf / 应该被阻止
	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}},
	}
	err := checker.Check(cmds)
	a.Error(err)

	// rm file.txt 应该允许
	cmds = []ExtractedCommand{
		{Name: "rm", Args: []string{"file.txt"}},
	}
	err = checker.Check(cmds)
	a.NoError(err)

	// rm -r /tmp 应该允许（因为参数是 -r 不是 -rf）
	cmds = []ExtractedCommand{
		{Name: "rm", Args: []string{"-r", "/tmp"}},
	}
	err = checker.Check(cmds)
	a.NoError(err)
}

func TestSecurityChecker_BlacklistWithAnyPositionArg(t *testing.T) {
	a := require.New(t)

	// 任意位置参数黑名单：阻止包含 /dev/null 的命令
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("*", MatchGlob).
			WithArgsFilter(ArgMatcher{Position: -1, Pattern: "/dev/null", Mode: MatchExact}))

	cmds := []ExtractedCommand{
		{Name: "cat", Args: []string{"/dev/null"}},
	}
	err := checker.Check(cmds)
	a.Error(err)

	cmds = []ExtractedCommand{
		{Name: "echo", Args: []string{"hello", ">", "/dev/null"}},
	}
	err = checker.Check(cmds)
	a.Error(err)

	cmds = []ExtractedCommand{
		{Name: "echo", Args: []string{"hello"}},
	}
	err = checker.Check(cmds)
	a.NoError(err)
}

func TestSecurityChecker_BlacklistWithRequiredArg(t *testing.T) {
	a := require.New(t)

	// 必需参数黑名单：阻止同时包含 -r 和 -f 的 rm
	// 注意：-rf 是一个参数，需要分别检查 -r 和 -f
	// 这里使用正则表达式来匹配
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact).
			WithArgsFilter(
				ArgMatcher{Position: -1, Pattern: `-r.*f|.*-f.*r`, Mode: MatchRegex, Required: true},
			))

	// rm -rf / 应该被阻止（同时有 -r 和 -f）
	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}},
	}
	err := checker.Check(cmds)
	a.Error(err)

	// rm -r / 应该允许（只有 -r）
	cmds = []ExtractedCommand{
		{Name: "rm", Args: []string{"-r", "/"}},
	}
	err = checker.Check(cmds)
	a.NoError(err)
}

// ============================================================
// TestSecurityChecker 白名单测试
// ============================================================

func TestSecurityChecker_WhitelistMode(t *testing.T) {
	a := require.New(t)

	// 白名单模式：只允许 ls, cat, echo
	checker := NewSecurityChecker().
		WithWhitelistMode(true).
		WithWhitelist(
			NewCommandRule("ls", MatchExact),
			NewCommandRule("cat", MatchExact),
			NewCommandRule("echo", MatchExact),
		)

	// ls 应该允许
	cmds := []ExtractedCommand{
		{Name: "ls", Args: []string{"-la"}},
	}
	err := checker.Check(cmds)
	a.NoError(err)

	// rm 应该被阻止
	cmds = []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}},
	}
	err = checker.Check(cmds)
	a.Error(err)
	a.Contains(err.Error(), "not in whitelist")
}

func TestSecurityChecker_WhitelistWithGlob(t *testing.T) {
	a := require.New(t)

	// 白名单模式：允许所有 git* 命令
	checker := NewSecurityChecker().
		WithWhitelistMode(true).
		WithWhitelist(NewCommandRule("git*", MatchGlob))

	cmds := []ExtractedCommand{
		{Name: "git", Args: []string{"status"}},
	}
	err := checker.Check(cmds)
	a.NoError(err)

	cmds = []ExtractedCommand{
		{Name: "gitconfig", Args: []string{}},
	}
	err = checker.Check(cmds)
	a.NoError(err)

	cmds = []ExtractedCommand{
		{Name: "npm", Args: []string{"install"}},
	}
	err = checker.Check(cmds)
	a.Error(err)
}

func TestSecurityChecker_WhitelistOffByDefault(t *testing.T) {
	a := require.New(t)

	// 默认不启用白名单模式
	checker := NewSecurityChecker()

	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}},
	}
	err := checker.Check(cmds)
	a.NoError(err) // 没有黑名单规则，命令允许执行
}

// ============================================================
// TestSecurityChecker 组合测试
// ============================================================

func TestSecurityChecker_BlacklistAndWhitelist(t *testing.T) {
	a := require.New(t)

	// 白名单模式 + 黑名单的组合
	// 允许所有 git* 命令，但阻止 git push
	checker := NewSecurityChecker().
		WithWhitelistMode(true).
		WithWhitelist(NewCommandRule("git*", MatchGlob)).
		WithBlacklist(NewCommandRule("git", MatchExact).
			WithArgsFilter(ArgMatcher{Position: 0, Pattern: "push", Mode: MatchExact}))

	// git status 允许
	cmds := []ExtractedCommand{
		{Name: "git", Args: []string{"status"}},
	}
	err := checker.Check(cmds)
	a.NoError(err)

	// git push 被黑名单阻止
	cmds = []ExtractedCommand{
		{Name: "git", Args: []string{"push"}},
	}
	err = checker.Check(cmds)
	a.Error(err)
	a.Contains(err.Error(), "blocked by blacklist")

	// npm 不在白名单中
	cmds = []ExtractedCommand{
		{Name: "npm", Args: []string{"install"}},
	}
	err = checker.Check(cmds)
	a.Error(err)
	a.Contains(err.Error(), "not in whitelist")
}

func TestSecurityChecker_MultipleCommands(t *testing.T) {
	a := require.New(t)

	// 检查多个命令
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact))

	cmds := []ExtractedCommand{
		{Name: "ls", Args: []string{"-la"}},
		{Name: "rm", Args: []string{"-rf", "/"}}, // 这个应该被阻止
		{Name: "cat", Args: []string{"file.txt"}},
	}
	err := checker.Check(cmds)
	a.Error(err)
	a.Contains(err.Error(), "rm")
}

func TestSecurityChecker_EmptyCommands(t *testing.T) {
	a := require.New(t)

	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact))

	// 空命令列表应该通过
	err := checker.Check([]ExtractedCommand{})
	a.NoError(err)
}

func TestSecurityChecker_EmptyRules(t *testing.T) {
	a := require.New(t)

	// 没有规则的检查器
	checker := NewSecurityChecker()

	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}},
	}
	err := checker.Check(cmds)
	a.NoError(err) // 没有规则，所有命令都允许
}

// ============================================================
// TestSecurityChecker 命令来源过滤测试
// ============================================================

func TestSecurityChecker_SourceFilter(t *testing.T) {
	a := require.New(t)

	// 只阻止命令替换中的 rm
	checker := NewSecurityChecker().
		WithBlacklist(NewCommandRule("rm", MatchExact).
			WithSourceFilter(SourceCmdSubst))

	// 直接调用的 rm 允许
	cmds := []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}, Source: SourceDirect},
	}
	err := checker.Check(cmds)
	a.NoError(err)

	// 命令替换中的 rm 被阻止
	cmds = []ExtractedCommand{
		{Name: "rm", Args: []string{"-rf", "/"}, Source: SourceCmdSubst},
	}
	err = checker.Check(cmds)
	a.Error(err)
}

// ============================================================
// TestCommandRule 测试命令规则的构建
// ============================================================

func TestNewCommandRule(t *testing.T) {
	a := require.New(t)

	rule := NewCommandRule("rm", MatchExact)
	a.Equal("rm", rule.Pattern)
	a.Equal(MatchExact, rule.MatchMode)
}

func TestCommandRule_WithArgsFilter(t *testing.T) {
	a := require.New(t)

	rule := NewCommandRule("rm", MatchExact).
		WithArgsFilter(
			ArgMatcher{Position: 0, Pattern: "-rf", Mode: MatchExact},
		)

	a.Len(rule.ArgsMatchers, 1)
	a.Equal(0, rule.ArgsMatchers[0].Position)
	a.Equal("-rf", rule.ArgsMatchers[0].Pattern)
}

func TestCommandRule_WithSourceFilter(t *testing.T) {
	a := require.New(t)

	rule := NewCommandRule("rm", MatchExact).
		WithSourceFilter(SourceCmdSubst, SourceSubshell)

	a.Len(rule.SourceFilter, 2)
	a.Contains(rule.SourceFilter, SourceCmdSubst)
	a.Contains(rule.SourceFilter, SourceSubshell)
}

// ============================================================
// TestSecurityChecker Builder 模式测试
// ============================================================

func TestSecurityChecker_Builder(t *testing.T) {
	a := require.New(t)

	// 测试 Builder 模式返回自身
	checker := NewSecurityChecker()
	a.IsType(&SecurityCheckerT{}, checker)

	// 链式调用
	checker2 := checker.
		WithBlacklist(NewCommandRule("rm", MatchExact)).
		WithWhitelistMode(true).
		WithWhitelist(NewCommandRule("ls", MatchExact))

	a.Equal(checker, checker2) // 应该返回相同的实例
}
