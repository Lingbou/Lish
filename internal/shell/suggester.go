package shell

import (
	"strings"
)

// Suggester 命令建议系统
type Suggester struct {
	history []string
}

// NewSuggester 创建新的建议器
func NewSuggester() *Suggester {
	return &Suggester{
		history: make([]string, 0),
	}
}

// AddToHistory 添加命令到历史
func (s *Suggester) AddToHistory(cmd string) {
	if cmd == "" {
		return
	}

	// 避免重复添加
	if len(s.history) > 0 && s.history[len(s.history)-1] == cmd {
		return
	}

	s.history = append(s.history, cmd)

	// 限制历史记录大小
	if len(s.history) > 1000 {
		s.history = s.history[len(s.history)-1000:]
	}
}

// Suggest 根据输入获取建议
func (s *Suggester) Suggest(input string) string {
	if input == "" {
		return ""
	}

	// 从后往前查找匹配的历史命令
	for i := len(s.history) - 1; i >= 0; i-- {
		if strings.HasPrefix(s.history[i], input) && s.history[i] != input {
			// 返回建议的后续部分
			return s.history[i][len(input):]
		}
	}

	return ""
}

// SpellCheck 拼写检查和纠正
func (s *Suggester) SpellCheck(cmd string, availableCommands []string) string {
	// 精确匹配
	for _, available := range availableCommands {
		if cmd == available {
			return ""
		}
	}

	// 查找最相似的命令
	minDist := 999
	bestMatch := ""

	for _, available := range availableCommands {
		dist := levenshteinDistance(cmd, available)
		if dist < minDist && dist <= 2 { // 最多允许 2 个字符差异
			minDist = dist
			bestMatch = available
		}
	}

	return bestMatch
}

// levenshteinDistance 计算编辑距离
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// 创建矩阵
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// 动态规划计算
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // 删除
				matrix[i][j-1]+1,      // 插入
				matrix[i-1][j-1]+cost, // 替换
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
