package commands

import (
	"context"
	"os"
	"testing"

	_ "strings" // 保留用于未来测试
)

func TestEchoCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "单个参数",
			args: []string{"hello"},
		},
		{
			name: "多个参数",
			args: []string{"hello", "world"},
		},
		{
			name: "空参数",
			args: []string{},
		},
		{
			name: "带特殊字符",
			args: []string{"hello", "world!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewEchoCommand(os.Stdout)
			ctx := context.Background()

			err := cmd.Execute(ctx, tt.args)
			if err != nil {
				t.Fatalf("Execute failed: %v", err)
			}
		})
	}
}

func TestEchoCommandFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "无换行符 -n",
			args: []string{"-n", "hello"},
		},
		{
			name: "无换行符 --no-newline",
			args: []string{"--no-newline", "hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewEchoCommand(os.Stdout)
			ctx := context.Background()

			err := cmd.Execute(ctx, tt.args)
			if err != nil {
				t.Fatalf("Execute failed: %v", err)
			}
		})
	}
}

func BenchmarkEchoCommand(b *testing.B) {
	cmd := NewEchoCommand(os.Stdout)
	ctx := context.Background()
	args := []string{"hello", "world", "benchmark", "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.Execute(ctx, args)
	}
}
