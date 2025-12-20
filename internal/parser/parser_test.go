package parser

import (
	"testing"
)

// TestParse 测试基本命令解析
func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "简单命令",
			input:    "ls",
			wantCmd:  "ls",
			wantArgs: []string{},
		},
		{
			name:     "带参数的命令",
			input:    "ls -la /tmp",
			wantCmd:  "ls",
			wantArgs: []string{"-la", "/tmp"},
		},
		{
			name:     "带引号的参数",
			input:    `echo "hello world"`,
			wantCmd:  "echo",
			wantArgs: []string{"hello world"},
		},
		{
			name:     "单引号参数",
			input:    `echo 'hello world'`,
			wantCmd:  "echo",
			wantArgs: []string{"hello world"},
		},
		{
			name:     "多个参数",
			input:    "cp file1.txt file2.txt /tmp",
			wantCmd:  "cp",
			wantArgs: []string{"file1.txt", "file2.txt", "/tmp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if parsed == nil {
				t.Fatal("Parse returned nil")
			}
			if parsed.Command != tt.wantCmd {
				t.Errorf("Command = %v, want %v", parsed.Command, tt.wantCmd)
			}
			if len(parsed.Args) != len(tt.wantArgs) {
				t.Errorf("Args length = %v, want %v", len(parsed.Args), len(tt.wantArgs))
				return
			}
			for i, arg := range parsed.Args {
				if arg != tt.wantArgs[i] {
					t.Errorf("Args[%d] = %v, want %v", i, arg, tt.wantArgs[i])
				}
			}
		})
	}
}

// TestParseEmpty 测试空命令
func TestParseEmpty(t *testing.T) {
	tests := []string{
		"",
		"   ",
		"\t",
		"\n",
	}

	for _, input := range tests {
		parsed, err := Parse(input)
		if err != nil {
			// 空命令可能返回错误，这是正常的
			continue
		}
		if parsed != nil && parsed.Command != "" {
			t.Errorf("Parse(%q) should return empty command, got %v", input, parsed.Command)
		}
	}
}

// TestParseComments 测试注释处理
func TestParseComments(t *testing.T) {
	tests := []struct {
		input   string
		wantCmd string
	}{
		{"ls # this is a comment", "ls"},
		{"echo hello # comment", "echo"},
	}

	for _, tt := range tests {
		parsed, err := Parse(tt.input)
		if err != nil && tt.wantCmd != "" {
			t.Errorf("Parse(%q) failed: %v", tt.input, err)
			continue
		}
		if parsed == nil && tt.wantCmd != "" {
			t.Errorf("Parse(%q) returned nil", tt.input)
			continue
		}
		if parsed != nil && parsed.Command != tt.wantCmd {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, parsed.Command, tt.wantCmd)
		}
	}
}

// TestParseRedirection 测试重定向解析
func TestParseRedirection(t *testing.T) {
	t.Skip("Redirection parsing is handled by pipeline parser")

	// 注意：基本的 Parse 函数不处理重定向
	// 重定向由 ParsePipeline 处理
	// 这里保留测试结构，但跳过执行
}

// BenchmarkParse 性能基准测试
func BenchmarkParse(b *testing.B) {
	input := "ls -la /tmp/test/directory"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse(input)
	}
}

// BenchmarkParseComplex 复杂命令解析性能测试
func BenchmarkParseComplex(b *testing.B) {
	input := `echo "hello world" | grep "hello" > output.txt`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse(input)
	}
}
