package theme

import (
	"strings"
)

// Highlighter 提供语法高亮功能
type Highlighter struct {
	scheme ColorScheme
}

// NewHighlighter 创建新的高亮器
func NewHighlighter(scheme ColorScheme) *Highlighter {
	return &Highlighter{scheme: scheme}
}

// Highlight 高亮命令行文本
func (h *Highlighter) Highlight(input string) string {
	// 简单实现：按空格分割，第一个词是命令
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return input
	}

	var result strings.Builder

	// 高亮命令
	result.WriteString(h.scheme.SyntaxCommand().Apply(parts[0]))

	// 高亮参数
	for i := 1; i < len(parts); i++ {
		result.WriteString(" ")
		arg := parts[i]

		// 检测不同类型
		if strings.HasPrefix(arg, "-") {
			// 选项
			result.WriteString(h.scheme.SyntaxArgument().Apply(arg))
		} else if strings.HasPrefix(arg, "$") {
			// 变量
			result.WriteString(h.scheme.SyntaxVariable().Apply(arg))
		} else if strings.ContainsAny(arg, `"'`) {
			// 字符串
			result.WriteString(h.scheme.SyntaxString().Apply(arg))
		} else if strings.ContainsAny(arg, "|><&;") {
			// 操作符
			result.WriteString(h.scheme.SyntaxOperator().Apply(arg))
		} else {
			// 普通参数
			result.WriteString(h.scheme.SyntaxArgument().Apply(arg))
		}
	}

	return result.String()
}

// HighlightFile 高亮文件名（用于 ls 等命令）
func (h *Highlighter) HighlightFile(name string, isDir, isExec, isLink bool) string {
	if isDir {
		return h.scheme.Directory().Apply(name)
	}
	if isExec {
		return h.scheme.Executable().Apply(name)
	}
	if isLink {
		return h.scheme.Symlink().Apply(name)
	}

	// 根据扩展名判断
	lowerName := strings.ToLower(name)
	if strings.HasSuffix(lowerName, ".zip") ||
		strings.HasSuffix(lowerName, ".tar") ||
		strings.HasSuffix(lowerName, ".gz") ||
		strings.HasSuffix(lowerName, ".7z") ||
		strings.HasSuffix(lowerName, ".rar") {
		return h.scheme.Archive().Apply(name)
	}

	return name
}

// UpdateScheme 更新配色方案
func (h *Highlighter) UpdateScheme(scheme ColorScheme) {
	h.scheme = scheme
}
