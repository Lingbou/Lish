package parser

import (
	"strings"
	"unicode"
)

// ParsedCommand 解析后的命令
type ParsedCommand struct {
	Command     string
	Args        []string
	Raw         string
	RedirectOut string // 输出重定向文件
	RedirectIn  string // 输入重定向文件
	RedirectErr string // 错误重定向文件
	Append      bool   // 是否追加模式
}

// Parse 解析命令行输入
func Parse(input string) (*ParsedCommand, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}

	tokens := tokenize(input)
	if len(tokens) == 0 {
		return nil, nil
	}

	return &ParsedCommand{
		Command: tokens[0],
		Args:    tokens[1:],
		Raw:     input,
	}, nil
}

// tokenize 将输入字符串分词，处理引号和转义字符
func tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	var inQuote rune

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		// 处理转义字符
		if ch == '\\' && i+1 < len(input) {
			next := rune(input[i+1])
			current.WriteRune(next)
			i++
			continue
		}

		// 处理引号
		if ch == '"' || ch == '\'' {
			if inQuote == 0 {
				inQuote = ch
			} else if inQuote == ch {
				inQuote = 0
			} else {
				current.WriteRune(ch)
			}
			continue
		}

		// 在引号内，所有字符都加入当前 token
		if inQuote != 0 {
			current.WriteRune(ch)
			continue
		}

		// 处理空格分隔
		if unicode.IsSpace(ch) {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteRune(ch)
	}

	// 添加最后一个 token
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
