package commands

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestPwdCommand(t *testing.T) {
	// 创建 pwd 命令
	cmd := NewPwdCommand(os.Stdout)

	// 执行命令
	ctx := context.Background()
	err := cmd.Execute(ctx, []string{})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// pwd 命令应该成功执行，输出到 stdout
	// 我们无法捕获 stdout，但可以验证命令不报错
}

func TestPwdCommandName(t *testing.T) {
	cmd := NewPwdCommand(os.Stdout)

	if cmd.Name() != "pwd" {
		t.Errorf("Name() = %v, want pwd", cmd.Name())
	}
}

func TestPwdCommandHelp(t *testing.T) {
	cmd := NewPwdCommand(os.Stdout)

	help := cmd.Help()
	if help == "" {
		t.Error("Help() should return non-empty string")
	}

	// 验证帮助信息包含关键字
	if !strings.Contains(help, "pwd") {
		t.Error("Help should contain command name")
	}
}

func BenchmarkPwdCommand(b *testing.B) {
	cmd := NewPwdCommand(os.Stdout)
	ctx := context.Background()
	args := []string{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.Execute(ctx, args)
	}
}
