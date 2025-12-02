package theme

import (
	"fmt"
	"strconv"
	"strings"
)

// ANSI 颜色码常量
const (
	ResetCode      = "\033[0m"
	BoldCode       = "\033[1m"
	DimCode        = "\033[2m"
	ItalicCode     = "\033[3m"
	UnderlineCode  = "\033[4m"
)

// Color 表示一个颜色配置
type Color struct {
	Foreground string   // 前景色 (hex 格式如 "#FF0000" 或 ANSI 码)
	Background string   // 背景色
	Style      []string // bold, italic, underline, dim
}

// ColorScheme 定义主题的颜色方案接口
type ColorScheme interface {
	// 基础颜色
	Primary() Color      // 主要颜色（命令名）
	Secondary() Color    // 次要颜色（参数）
	Success() Color      // 成功（绿色）
	Warning() Color      // 警告（黄色）
	Error() Color        // 错误（红色）
	Info() Color         // 信息（蓝色）

	// 文件类型颜色
	Directory() Color    // 目录
	Executable() Color   // 可执行文件
	Symlink() Color      // 符号链接
	Archive() Color      // 压缩包

	// 提示符颜色
	PromptUser() Color   // 用户名
	PromptHost() Color   // 主机名
	PromptPath() Color   // 路径
	PromptGit() Color    // Git 分支

	// 语法高亮
	SyntaxCommand() Color    // 命令
	SyntaxArgument() Color   // 参数
	SyntaxString() Color     // 字符串
	SyntaxVariable() Color   // 变量
	SyntaxOperator() Color   // 操作符 (|, >, <)
}

// Theme 表示一个完整主题
type Theme struct {
	Name        string
	Description string
	Author      string
	Scheme      ColorScheme
}

// ToANSI 将颜色转换为 ANSI 转义序列
func (c Color) ToANSI() string {
	var codes []string

	// 处理样式
	for _, style := range c.Style {
		switch strings.ToLower(style) {
		case "bold":
			codes = append(codes, "1")
		case "dim":
			codes = append(codes, "2")
		case "italic":
			codes = append(codes, "3")
		case "underline":
			codes = append(codes, "4")
		}
	}

	// 处理前景色
	if c.Foreground != "" {
		if strings.HasPrefix(c.Foreground, "#") {
			// Hex 颜色转 RGB
			r, g, b := hexToRGB(c.Foreground)
			codes = append(codes, fmt.Sprintf("38;2;%d;%d;%d", r, g, b))
		} else {
			// 直接使用 ANSI 码
			codes = append(codes, c.Foreground)
		}
	}

	// 处理背景色
	if c.Background != "" {
		if strings.HasPrefix(c.Background, "#") {
			r, g, b := hexToRGB(c.Background)
			codes = append(codes, fmt.Sprintf("48;2;%d;%d;%d", r, g, b))
		} else {
			codes = append(codes, c.Background)
		}
	}

	if len(codes) == 0 {
		return ""
	}

	return "\033[" + strings.Join(codes, ";") + "m"
}

// Apply 应用颜色到文本
func (c Color) Apply(text string) string {
	ansi := c.ToANSI()
	if ansi == "" {
		return text
	}
	return ansi + text + ResetCode
}

// hexToRGB 将 hex 颜色转换为 RGB 值
func hexToRGB(hex string) (r, g, b int) {
	hex = strings.TrimPrefix(hex, "#")
	
	if len(hex) == 3 {
		// 短格式 #RGB -> #RRGGBB
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	
	if len(hex) != 6 {
		return 0, 0, 0
	}

	r64, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g64, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b64, _ := strconv.ParseInt(hex[4:6], 16, 64)

	return int(r64), int(g64), int(b64)
}

// NewColor 创建一个新的颜色
func NewColor(fg string, styles ...string) Color {
	return Color{
		Foreground: fg,
		Style:      styles,
	}
}

// NewColorWithBg 创建带背景色的颜色
func NewColorWithBg(fg, bg string, styles ...string) Color {
	return Color{
		Foreground: fg,
		Background: bg,
		Style:      styles,
	}
}

