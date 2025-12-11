package script

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

// TokenType 表示 token 类型
type TokenType int

const (
	// 特殊 token
	TOKEN_EOF TokenType = iota
	TOKEN_ILLEGAL
	TOKEN_COMMENT

	// 标识符和字面量
	TOKEN_IDENT  // 标识符
	TOKEN_STRING // 字符串
	TOKEN_NUMBER // 数字

	// 关键字
	TOKEN_IF
	TOKEN_THEN
	TOKEN_ELIF
	TOKEN_ELSE
	TOKEN_FI
	TOKEN_FOR
	TOKEN_IN
	TOKEN_DO
	TOKEN_DONE
	TOKEN_WHILE
	TOKEN_FUNCTION
	TOKEN_RETURN
	TOKEN_BREAK
	TOKEN_CONTINUE
	TOKEN_LOCAL

	// 操作符
	TOKEN_ASSIGN    // =
	TOKEN_EQ        // ==
	TOKEN_NE        // !=
	TOKEN_LT        // <
	TOKEN_GT        // >
	TOKEN_LE        // <=
	TOKEN_GE        // >=
	TOKEN_AND       // &&
	TOKEN_OR        // ||
	TOKEN_NOT       // !
	TOKEN_PIPE      // |
	TOKEN_REDIRECT  // >, >>, <
	TOKEN_SEMICOLON // ;

	// 分隔符
	TOKEN_LPAREN   // (
	TOKEN_RPAREN   // )
	TOKEN_LBRACE   // {
	TOKEN_RBRACE   // }
	TOKEN_LBRACKET // [
	TOKEN_RBRACKET // ]
	TOKEN_NEWLINE  // \n
)

// Token 表示一个词法单元
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Lexer 词法分析器
type Lexer struct {
	reader *bufio.Reader
	line   int
	column int
	ch     rune
}

// NewLexer 创建新的词法分析器
func NewLexer(input string) *Lexer {
	l := &Lexer{
		reader: bufio.NewReader(strings.NewReader(input)),
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar 读取下一个字符
func (l *Lexer) readChar() {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.ch = 0
		}
		return
	}
	l.ch = ch
	l.column++
}

// peekChar 查看下一个字符但不移动位置
func (l *Lexer) peekChar() rune {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		return 0
	}
	l.reader.UnreadRune()
	return ch
}

// NextToken 获取下一个 token
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	tok := Token{
		Line:   l.line,
		Column: l.column,
	}

	switch l.ch {
	case 0:
		tok.Type = TOKEN_EOF
		tok.Literal = ""
	case '#':
		tok.Type = TOKEN_COMMENT
		tok.Literal = l.readComment()
	case '\n':
		tok.Type = TOKEN_NEWLINE
		tok.Literal = "\n"
		l.line++
		l.column = 0
		l.readChar()
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TOKEN_EQ
			tok.Literal = "=="
			l.readChar()
		} else {
			tok.Type = TOKEN_ASSIGN
			tok.Literal = "="
			l.readChar()
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TOKEN_NE
			tok.Literal = "!="
			l.readChar()
		} else {
			tok.Type = TOKEN_NOT
			tok.Literal = "!"
			l.readChar()
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TOKEN_LE
			tok.Literal = "<="
			l.readChar()
		} else {
			tok.Type = TOKEN_LT
			tok.Literal = "<"
			l.readChar()
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TOKEN_GE
			tok.Literal = ">="
			l.readChar()
		} else if l.peekChar() == '>' {
			l.readChar()
			tok.Type = TOKEN_REDIRECT
			tok.Literal = ">>"
			l.readChar()
		} else {
			tok.Type = TOKEN_GT
			tok.Literal = ">"
			l.readChar()
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok.Type = TOKEN_AND
			tok.Literal = "&&"
			l.readChar()
		} else {
			tok.Type = TOKEN_ILLEGAL
			tok.Literal = string(l.ch)
			l.readChar()
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok.Type = TOKEN_OR
			tok.Literal = "||"
			l.readChar()
		} else {
			tok.Type = TOKEN_PIPE
			tok.Literal = "|"
			l.readChar()
		}
	case ';':
		tok.Type = TOKEN_SEMICOLON
		tok.Literal = ";"
		l.readChar()
	case '(':
		tok.Type = TOKEN_LPAREN
		tok.Literal = "("
		l.readChar()
	case ')':
		tok.Type = TOKEN_RPAREN
		tok.Literal = ")"
		l.readChar()
	case '{':
		tok.Type = TOKEN_LBRACE
		tok.Literal = "{"
		l.readChar()
	case '}':
		tok.Type = TOKEN_RBRACE
		tok.Literal = "}"
		l.readChar()
	case '[':
		tok.Type = TOKEN_LBRACKET
		tok.Literal = "["
		l.readChar()
	case ']':
		tok.Type = TOKEN_RBRACKET
		tok.Literal = "]"
		l.readChar()
	case '"', '\'':
		tok.Type = TOKEN_STRING
		tok.Literal = l.readString(l.ch)
	default:
		if isLetter(l.ch) || l.ch == '_' || l.ch == '$' {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupKeyword(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TOKEN_NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok.Type = TOKEN_ILLEGAL
			tok.Literal = string(l.ch)
			l.readChar()
		}
	}

	return tok
}

// skipWhitespace 跳过空白字符（不包括换行）
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// readComment 读取注释
func (l *Lexer) readComment() string {
	var buf bytes.Buffer
	for l.ch != '\n' && l.ch != 0 {
		buf.WriteRune(l.ch)
		l.readChar()
	}
	return buf.String()
}

// readIdentifier 读取标识符
func (l *Lexer) readIdentifier() string {
	var buf bytes.Buffer
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '$' {
		buf.WriteRune(l.ch)
		l.readChar()
	}
	return buf.String()
}

// readNumber 读取数字
func (l *Lexer) readNumber() string {
	var buf bytes.Buffer
	for isDigit(l.ch) {
		buf.WriteRune(l.ch)
		l.readChar()
	}
	return buf.String()
}

// readString 读取字符串
func (l *Lexer) readString(quote rune) string {
	var buf bytes.Buffer
	l.readChar() // 跳过开始引号

	for l.ch != quote && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			// 处理转义字符
			switch l.ch {
			case 'n':
				buf.WriteRune('\n')
			case 't':
				buf.WriteRune('\t')
			case 'r':
				buf.WriteRune('\r')
			case '\\':
				buf.WriteRune('\\')
			case quote:
				buf.WriteRune(quote)
			default:
				buf.WriteRune(l.ch)
			}
		} else {
			buf.WriteRune(l.ch)
		}
		l.readChar()
	}

	if l.ch == quote {
		l.readChar() // 跳过结束引号
	}

	return buf.String()
}

// isLetter 判断是否是字母
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

// isDigit 判断是否是数字
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// lookupKeyword 查找关键字
func lookupKeyword(ident string) TokenType {
	keywords := map[string]TokenType{
		"if":       TOKEN_IF,
		"then":     TOKEN_THEN,
		"elif":     TOKEN_ELIF,
		"else":     TOKEN_ELSE,
		"fi":       TOKEN_FI,
		"for":      TOKEN_FOR,
		"in":       TOKEN_IN,
		"do":       TOKEN_DO,
		"done":     TOKEN_DONE,
		"while":    TOKEN_WHILE,
		"function": TOKEN_FUNCTION,
		"return":   TOKEN_RETURN,
		"break":    TOKEN_BREAK,
		"continue": TOKEN_CONTINUE,
		"local":    TOKEN_LOCAL,
	}

	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}
