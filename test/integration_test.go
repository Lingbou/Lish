package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Lingbou/Lish/internal/commands"
	"github.com/Lingbou/Lish/internal/parser"
)

// TestBasicCommands 测试基本命令执行
func TestBasicCommands(t *testing.T) {
	// 创建命令注册表
	registry := commands.NewRegistry()

	// 注册基本命令
	registry.Register(commands.NewPwdCommand(os.Stdout))
	registry.Register(commands.NewEchoCommand(os.Stdout))

	tests := []struct {
		name    string
		cmdLine string
		wantErr bool
	}{
		{
			name:    "pwd命令",
			cmdLine: "pwd",
			wantErr: false,
		},
		{
			name:    "echo命令",
			cmdLine: "echo hello",
			wantErr: false,
		},
		{
			name:    "不存在的命令",
			cmdLine: "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parser.Parse(tt.cmdLine)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Parse failed: %v", err)
				}
				return
			}
			if parsed == nil {
				t.Fatal("Parse returned nil")
			}

			cmd, exists := registry.Get(parsed.Command)
			if !exists {
				if !tt.wantErr {
					t.Errorf("Command %q not found", parsed.Command)
				}
				return
			}

			ctx := context.Background()
			execErr := cmd.Execute(ctx, parsed.Args)
			if (execErr != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", execErr, tt.wantErr)
			}
		})
	}
}

// TestFileOperations 测试文件操作命令
func TestFileOperations(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello, World!")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 测试 cat 命令
	t.Run("cat command", func(t *testing.T) {
		cmd := commands.NewCatCommand(os.Stdout)
		ctx := context.Background()
		err := cmd.Execute(ctx, []string{testFile})
		if err != nil {
			t.Errorf("cat command failed: %v", err)
		}
	})

	// 测试 ls 命令
	t.Run("ls command", func(t *testing.T) {
		cmd := commands.NewLsCommand(os.Stdout)
		ctx := context.Background()
		err := cmd.Execute(ctx, []string{tmpDir})
		if err != nil {
			t.Errorf("ls command failed: %v", err)
		}
	})
}

// TestCommandPipeline 测试命令管道
func TestCommandPipeline(t *testing.T) {
	t.Skip("Pipeline test requires shell executor")

	// TODO: 实现管道测试
	// 例如: echo "hello" | grep "hello"
}

// TestCommandRedirection 测试命令重定向
func TestCommandRedirection(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	// 解析带重定向的命令
	parsed, err := parser.Parse("echo hello > " + outputFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if parsed == nil {
		t.Fatal("Parse returned nil")
	}

	if parsed.RedirectOut != outputFile {
		t.Errorf("RedirectOut = %v, want %v", parsed.RedirectOut, outputFile)
	}

	// TODO: 测试实际的重定向执行
}

// BenchmarkCommandExecution 命令执行性能测试
func BenchmarkCommandExecution(b *testing.B) {
	cmd := commands.NewPwdCommand(os.Stdout)
	ctx := context.Background()
	args := []string{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.Execute(ctx, args)
	}
}

// BenchmarkParseAndExecute 解析和执行性能测试
func BenchmarkParseAndExecute(b *testing.B) {
	registry := commands.NewRegistry()
	registry.Register(commands.NewEchoCommand(os.Stdout))

	cmdLine := "echo hello world"
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parsed, _ := parser.Parse(cmdLine)
		if parsed != nil {
			if cmd, exists := registry.Get(parsed.Command); exists {
				cmd.Execute(ctx, parsed.Args)
			}
		}
	}
}
