package parser

import (
	"strings"
	"unicode"
)

// Lexer 词法分析器
type Lexer struct {
	input  string
	pos    int
	ch     rune
	tokens []Token
}

// NewLexer 创建词法分析器
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make([]Token, 0),
	}
	l.readChar()
	return l
}

// Tokenize 将输入转换为 token 序列
func (l *Lexer) Tokenize() []Token {
	for l.ch != 0 {
		l.skipWhitespace()

		if l.ch == 0 {
			break
		}

		// 检查操作符
		switch l.ch {
		case '|':
			l.tokens = append(l.tokens, Token{Type: TokenPipe, Value: "|"})
			l.readChar()
		case '>':
			l.readChar()
			if l.ch == '>' {
				l.tokens = append(l.tokens, Token{Type: TokenRedirectAppend, Value: ">>"})
				l.readChar()
			} else {
				l.tokens = append(l.tokens, Token{Type: TokenRedirectOut, Value: ">"})
			}
		case '<':
			l.tokens = append(l.tokens, Token{Type: TokenRedirectIn, Value: "<"})
			l.readChar()
		case '&':
			l.readChar()
			if l.ch == '&' {
				l.tokens = append(l.tokens, Token{Type: TokenAnd, Value: "&&"})
				l.readChar()
			} else {
				// 单个 & 当作普通字符
				l.tokens = append(l.tokens, Token{Type: TokenWord, Value: "&"})
			}
		case ';':
			l.tokens = append(l.tokens, Token{Type: TokenSemicolon, Value: ";"})
			l.readChar()
		case '"', '\'':
			// 处理引号字符串
			quote := l.ch
			l.readChar()
			word := l.readQuotedString(quote)
			l.tokens = append(l.tokens, Token{Type: TokenWord, Value: word})
		default:
			// 普通单词
			word := l.readWord()
			if word != "" {
				// 特殊处理 2>
				if word == "2" && l.ch == '>' {
					l.readChar()
					l.tokens = append(l.tokens, Token{Type: TokenRedirectErr, Value: "2>"})
				} else {
					l.tokens = append(l.tokens, Token{Type: TokenWord, Value: word})
				}
			}
		}
	}

	l.tokens = append(l.tokens, Token{Type: TokenEOF, Value: ""})
	return l.tokens
}

func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = rune(l.input[l.pos])
		l.pos++
	}
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readWord() string {
	var result strings.Builder

	for l.ch != 0 && !unicode.IsSpace(l.ch) && !l.isOperator(l.ch) {
		if l.ch == '\\' {
			l.readChar()
			if l.ch != 0 {
				result.WriteRune(l.ch)
				l.readChar()
			}
		} else {
			result.WriteRune(l.ch)
			l.readChar()
		}
	}

	return result.String()
}

func (l *Lexer) readQuotedString(quote rune) string {
	var result strings.Builder

	for l.ch != 0 && l.ch != quote {
		if l.ch == '\\' {
			l.readChar()
			if l.ch != 0 {
				result.WriteRune(l.ch)
				l.readChar()
			}
		} else {
			result.WriteRune(l.ch)
			l.readChar()
		}
	}

	if l.ch == quote {
		l.readChar()
	}

	return result.String()
}

func (l *Lexer) isOperator(ch rune) bool {
	return ch == '|' || ch == '>' || ch == '<' || ch == '&' || ch == ';'
}
