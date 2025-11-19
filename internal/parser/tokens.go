package parser

// TokenType 定义 token 类型
type TokenType int

const (
	TokenWord           TokenType = iota
	TokenPipe                     // |
	TokenRedirectOut              // >
	TokenRedirectAppend           // >>
	TokenRedirectIn               // <
	TokenRedirectErr              // 2>
	TokenAnd                      // &&
	TokenOr                       // ||
	TokenSemicolon                // ;
	TokenEOF
)

// Token 表示一个词法单元
type Token struct {
	Type  TokenType
	Value string
}
