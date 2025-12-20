package theme

import (
	"strings"
	"testing"
)

func TestNewColor(t *testing.T) {
	tests := []struct {
		name       string
		foreground string
		background string
		styles     []string
	}{
		{
			name:       "基本颜色",
			foreground: "#FF0000",
			background: "",
			styles:     []string{},
		},
		{
			name:       "带背景色",
			foreground: "#00FF00",
			background: "#000000",
			styles:     []string{},
		},
		{
			name:       "带样式",
			foreground: "#0000FF",
			background: "",
			styles:     []string{"bold", "underline"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := Color{
				Foreground: tt.foreground,
				Background: tt.background,
				Style:      tt.styles,
			}

			if color.Foreground != tt.foreground {
				t.Errorf("Foreground = %v, want %v", color.Foreground, tt.foreground)
			}
			if color.Background != tt.background {
				t.Errorf("Background = %v, want %v", color.Background, tt.background)
			}
		})
	}
}

func TestColorToANSI(t *testing.T) {
	tests := []struct {
		name  string
		color Color
	}{
		{
			name: "红色",
			color: Color{
				Foreground: "#FF0000",
			},
		},
		{
			name: "绿色",
			color: Color{
				Foreground: "#00FF00",
			},
		},
		{
			name: "蓝色带粗体",
			color: Color{
				Foreground: "#0000FF",
				Style:      []string{"bold"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ansi := tt.color.ToANSI()
			
			// 验证 ANSI 码格式
			if !strings.HasPrefix(ansi, "\033[") {
				t.Errorf("ANSI code should start with \\033[, got %q", ansi)
			}
			if !strings.HasSuffix(ansi, "m") {
				t.Errorf("ANSI code should end with m, got %q", ansi)
			}
		})
	}
}

func TestColorApply(t *testing.T) {
	color := Color{
		Foreground: "#FF0000",
	}

	text := "Hello, World!"
	result := color.Apply(text)

	// 验证结果包含原文本
	if !strings.Contains(result, text) {
		t.Errorf("Apply result should contain original text")
	}

	// 验证结果包含 ANSI 码
	if !strings.Contains(result, "\033[") {
		t.Errorf("Apply result should contain ANSI codes")
	}

	// 验证结果包含重置码
	if !strings.Contains(result, "\033[0m") {
		t.Errorf("Apply result should contain reset code")
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		hex  string
		r, g, b int
	}{
		{"#FF0000", 255, 0, 0},
		{"#00FF00", 0, 255, 0},
		{"#0000FF", 0, 0, 255},
		{"#FFFFFF", 255, 255, 255},
		{"#000000", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			r, g, b := hexToRGB(tt.hex)
			if r != tt.r || g != tt.g || b != tt.b {
				t.Errorf("hexToRGB(%q) = (%d, %d, %d), want (%d, %d, %d)",
					tt.hex, r, g, b, tt.r, tt.g, tt.b)
			}
		})
	}
}

func BenchmarkColorToANSI(b *testing.B) {
	color := Color{
		Foreground: "#FF0000",
		Background: "#000000",
		Style:      []string{"bold", "underline"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		color.ToANSI()
	}
}

func BenchmarkColorApply(b *testing.B) {
	color := Color{
		Foreground: "#FF0000",
	}
	text := "Hello, World! This is a benchmark test."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		color.Apply(text)
	}
}

