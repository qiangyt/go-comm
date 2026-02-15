package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ============================================================
// TestExtractCommands 测试命令提取器的基本功能
// ============================================================

func TestExtractCommands_SimpleCommand(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 简单命令
	cmds, err := extractor.Extract("ls -la")
	a.NoError(err)
	a.Len(cmds, 1)
	a.Equal("ls", cmds[0].Name)
	a.Equal([]string{"-la"}, cmds[0].Args)
	a.Equal(SourceDirect, cmds[0].Source)
}

func TestExtractCommands_CommandWithVariables(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 带变量的命令
	cmds, err := extractor.Extract("ls -la $HOME")
	a.NoError(err)
	a.Len(cmds, 1)
	a.Equal("ls", cmds[0].Name)
	// 变量会被保留为原始形式
	a.Equal([]string{"-la", "$HOME"}, cmds[0].Args)
}

func TestExtractCommands_Pipeline(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 带管道的命令
	cmds, err := extractor.Extract("ls -la | grep test")
	a.NoError(err)
	a.Len(cmds, 2)
	a.Equal("ls", cmds[0].Name)
	a.Equal([]string{"-la"}, cmds[0].Args)
	a.Equal("grep", cmds[1].Name)
	a.Equal([]string{"test"}, cmds[1].Args)
}

func TestExtractCommands_CommandSubstitution(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 带命令替换
	cmds, err := extractor.Extract("echo $(rm -rf /)")
	a.NoError(err)
	a.Len(cmds, 2)

	// 外层命令
	a.Equal("echo", cmds[0].Name)
	a.Equal(SourceDirect, cmds[0].Source)

	// 命令替换中的命令
	a.Equal("rm", cmds[1].Name)
	a.Equal(SourceCmdSubst, cmds[1].Source)
	a.Equal([]string{"-rf", "/"}, cmds[1].Args)
}

func TestExtractCommands_BackquoteSubstitution(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 使用反引号的命令替换
	cmds, err := extractor.Extract("echo `ls -la /tmp`")
	a.NoError(err)
	a.Len(cmds, 2)

	a.Equal("echo", cmds[0].Name)
	a.Equal(SourceDirect, cmds[0].Source)

	a.Equal("ls", cmds[1].Name)
	a.Equal(SourceCmdSubst, cmds[1].Source)
	a.Equal([]string{"-la", "/tmp"}, cmds[1].Args)
}

func TestExtractCommands_ConditionalCommand(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 条件命令 &&
	cmds, err := extractor.Extract("cd /tmp && ls -la")
	a.NoError(err)
	a.Len(cmds, 2)
	a.Equal("cd", cmds[0].Name)
	a.Equal([]string{"/tmp"}, cmds[0].Args)
	a.Equal("ls", cmds[1].Name)
	a.Equal([]string{"-la"}, cmds[1].Args)
}

func TestExtractCommands_ConditionalOrCommand(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 条件命令 ||
	cmds, err := extractor.Extract("cd /tmp || echo failed")
	a.NoError(err)
	a.Len(cmds, 2)
	a.Equal("cd", cmds[0].Name)
	a.Equal("echo", cmds[1].Name)
}

func TestExtractCommands_QuotedArguments(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 带引号的命令
	cmds, err := extractor.Extract(`echo "hello world"`)
	a.NoError(err)
	a.Len(cmds, 1)
	a.Equal("echo", cmds[0].Name)
	// 引号内的内容作为一个参数
	a.Len(cmds[0].Args, 1)
}

func TestExtractCommands_MultipleStatements(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 多行命令（分号分隔）
	cmds, err := extractor.Extract("ls; ls -la")
	a.NoError(err)
	a.Len(cmds, 2)
	a.Equal("ls", cmds[0].Name)
	a.Len(cmds[0].Args, 0) // 空参数
	a.Equal("ls", cmds[1].Name)
	a.Equal([]string{"-la"}, cmds[1].Args)
}

func TestExtractCommands_Subshell(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 子 shell
	cmds, err := extractor.Extract("(cd /tmp && ls)")
	a.NoError(err)
	a.Len(cmds, 2)

	a.Equal("cd", cmds[0].Name)
	a.Equal(SourceSubshell, cmds[0].Source)

	a.Equal("ls", cmds[1].Name)
	a.Equal(SourceSubshell, cmds[1].Source)
}

func TestExtractCommands_NestedCommandSubstitution(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 嵌套命令替换
	cmds, err := extractor.Extract("echo $(echo $(rm -rf /))")
	a.NoError(err)
	a.Len(cmds, 3)

	a.Equal("echo", cmds[0].Name)
	a.Equal(SourceDirect, cmds[0].Source)

	a.Equal("echo", cmds[1].Name)
	a.Equal(SourceCmdSubst, cmds[1].Source)

	a.Equal("rm", cmds[2].Name)
	a.Equal(SourceCmdSubst, cmds[2].Source)
	a.Equal([]string{"-rf", "/"}, cmds[2].Args)
}

func TestExtractCommands_EnvironmentVariables(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 带环境变量的命令
	cmds, err := extractor.Extract("VAR=value ls -la")
	a.NoError(err)
	a.Len(cmds, 1)

	a.Equal("ls", cmds[0].Name)
	a.Equal([]string{"-la"}, cmds[0].Args)
	// 环境变量应该被记录
	a.Len(cmds[0].Envs, 1)
	a.Equal("VAR=value", cmds[0].Envs[0])
}

func TestExtractCommands_EmptyCommand(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 空命令
	cmds, err := extractor.Extract("")
	a.NoError(err)
	a.Len(cmds, 0)
}

func TestExtractCommands_CommentOnly(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 只有注释
	cmds, err := extractor.Extract("# this is a comment")
	a.NoError(err)
	a.Len(cmds, 0)
}

func TestExtractCommands_InvalidSyntax(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 无效语法
	cmds, err := extractor.Extract("echo 'unclosed")
	a.Error(err)
	a.Len(cmds, 0)
}

func TestExtractCommands_ComplexPipeline(t *testing.T) {
	a := require.New(t)

	extractor := NewCommandExtractor()

	// 复杂管道
	cmds, err := extractor.Extract("cat file | grep pattern | wc -l")
	a.NoError(err)
	a.Len(cmds, 3)
	a.Equal("cat", cmds[0].Name)
	a.Equal("grep", cmds[1].Name)
	a.Equal("wc", cmds[2].Name)
}
