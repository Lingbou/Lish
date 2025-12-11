package script

import (
	"fmt"
	"strings"
)

// Parser 语法分析器
type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

// NewParser 创建新的语法分析器
func NewParser(input string) *Parser {
	p := &Parser{
		lexer:  NewLexer(input),
		errors: []string{},
	}

	// 读取两个 token，初始化 curToken 和 peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken 移动到下一个 token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()

	// 跳过注释
	for p.peekToken.Type == TOKEN_COMMENT {
		p.peekToken = p.lexer.NextToken()
	}
}

// Errors 返回解析错误
func (p *Parser) Errors() []string {
	return p.errors
}

// addError 添加错误信息
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("行 %d:%d - %s",
		p.curToken.Line, p.curToken.Column, msg))
}

// expectToken 期望特定类型的 token
func (p *Parser) expectToken(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("期望 %v，得到 %v", t, p.peekToken.Type))
	return false
}

// skipNewlines 跳过换行符
func (p *Parser) skipNewlines() {
	for p.curToken.Type == TOKEN_NEWLINE || p.peekToken.Type == TOKEN_NEWLINE {
		p.nextToken()
	}
}

// Parse 解析脚本并返回 AST
func (p *Parser) Parse() *Program {
	program := &Program{
		Statements: []Statement{},
	}

	for p.curToken.Type != TOKEN_EOF {
		p.skipNewlines()

		if p.curToken.Type == TOKEN_EOF {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

// parseStatement 解析语句
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case TOKEN_IF:
		return p.parseIfStatement()
	case TOKEN_FOR:
		return p.parseForStatement()
	case TOKEN_WHILE:
		return p.parseWhileStatement()
	case TOKEN_FUNCTION:
		return p.parseFunctionDef()
	case TOKEN_RETURN:
		return p.parseReturnStatement()
	case TOKEN_BREAK:
		return &BreakStatement{}
	case TOKEN_CONTINUE:
		return &ContinueStatement{}
	case TOKEN_IDENT:
		// 检查是否是赋值语句
		if p.peekToken.Type == TOKEN_ASSIGN {
			return p.parseAssignStatement()
		}
		// 否则是命令语句
		return p.parseCommandStatement()
	default:
		// 默认作为命令语句处理
		return p.parseCommandStatement()
	}
}

// parseIfStatement 解析 if 语句
func (p *Parser) parseIfStatement() Statement {
	stmt := &IfStatement{
		ElseIfList: []*ElseIfClause{},
	}

	p.nextToken() // 跳过 'if'

	// 解析条件
	stmt.Condition = p.parseCondition()

	// 期望 'then'
	if !p.expectToken(TOKEN_THEN) {
		return nil
	}
	p.skipNewlines()

	// 解析 then 块
	stmt.ThenBlock = p.parseBlock(TOKEN_ELIF, TOKEN_ELSE, TOKEN_FI)

	// 处理 elif
	for p.curToken.Type == TOKEN_ELIF {
		p.nextToken()
		elseif := &ElseIfClause{}
		elseif.Condition = p.parseCondition()

		if !p.expectToken(TOKEN_THEN) {
			return nil
		}
		p.skipNewlines()

		elseif.Block = p.parseBlock(TOKEN_ELIF, TOKEN_ELSE, TOKEN_FI)
		stmt.ElseIfList = append(stmt.ElseIfList, elseif)
	}

	// 处理 else
	if p.curToken.Type == TOKEN_ELSE {
		p.nextToken()
		p.skipNewlines()
		stmt.ElseBlock = p.parseBlock(TOKEN_FI)
	}

	// 期望 'fi'
	if p.curToken.Type != TOKEN_FI {
		p.addError("if 语句缺少 'fi'")
	}

	return stmt
}

// parseForStatement 解析 for 语句
func (p *Parser) parseForStatement() Statement {
	stmt := &ForStatement{}

	p.nextToken() // 跳过 'for'

	// 读取变量名
	if p.curToken.Type != TOKEN_IDENT {
		p.addError("for 语句需要变量名")
		return nil
	}
	stmt.Variable = p.curToken.Literal

	// 期望 'in'
	if !p.expectToken(TOKEN_IN) {
		return nil
	}
	p.nextToken()

	// 读取列表
	for p.curToken.Type != TOKEN_DO && p.curToken.Type != TOKEN_NEWLINE && p.curToken.Type != TOKEN_SEMICOLON {
		stmt.List = append(stmt.List, p.curToken.Literal)
		p.nextToken()
	}

	p.skipNewlines()

	// 期望 'do'
	if p.curToken.Type == TOKEN_SEMICOLON {
		p.nextToken()
		p.skipNewlines()
	}
	if p.curToken.Type != TOKEN_DO {
		p.addError("for 语句需要 'do'")
		return nil
	}
	p.nextToken()
	p.skipNewlines()

	// 解析循环体
	stmt.Block = p.parseBlock(TOKEN_DONE)

	// 期望 'done'
	if p.curToken.Type != TOKEN_DONE {
		p.addError("for 语句缺少 'done'")
	}

	return stmt
}

// parseWhileStatement 解析 while 语句
func (p *Parser) parseWhileStatement() Statement {
	stmt := &WhileStatement{}

	p.nextToken() // 跳过 'while'

	// 解析条件
	stmt.Condition = p.parseCondition()

	// 期望 'do'
	if !p.expectToken(TOKEN_DO) {
		return nil
	}
	p.skipNewlines()

	// 解析循环体
	stmt.Block = p.parseBlock(TOKEN_DONE)

	// 期望 'done'
	if p.curToken.Type != TOKEN_DONE {
		p.addError("while 语句缺少 'done'")
	}

	return stmt
}

// parseFunctionDef 解析函数定义
func (p *Parser) parseFunctionDef() Statement {
	p.nextToken() // 跳过 'function'

	if p.curToken.Type != TOKEN_IDENT {
		p.addError("函数定义需要函数名")
		return nil
	}

	fn := &FunctionDef{
		Name: p.curToken.Literal,
	}

	p.nextToken()

	// 期望 '()'
	if p.curToken.Type == TOKEN_LPAREN {
		if !p.expectToken(TOKEN_RPAREN) {
			return nil
		}
		p.nextToken()
	}

	p.skipNewlines()

	// 期望 '{'
	if p.curToken.Type != TOKEN_LBRACE {
		p.addError("函数定义需要 '{'")
		return nil
	}
	p.nextToken()
	p.skipNewlines()

	// 解析函数体
	fn.Block = p.parseBlock(TOKEN_RBRACE)

	// 期望 '}'
	if p.curToken.Type != TOKEN_RBRACE {
		p.addError("函数定义缺少 '}'")
	}

	return fn
}

// parseReturnStatement 解析 return 语句
func (p *Parser) parseReturnStatement() Statement {
	stmt := &ReturnStatement{}

	p.nextToken() // 跳过 'return'

	// 如果后面有值，解析它
	if p.curToken.Type != TOKEN_NEWLINE && p.curToken.Type != TOKEN_SEMICOLON && p.curToken.Type != TOKEN_EOF {
		stmt.Value = &StringLiteral{Value: p.curToken.Literal}
	}

	return stmt
}

// parseAssignStatement 解析赋值语句
func (p *Parser) parseAssignStatement() Statement {
	stmt := &AssignStatement{
		Name: p.curToken.Literal,
	}

	p.nextToken() // 跳过变量名
	p.nextToken() // 跳过 '='

	// 解析值
	if p.curToken.Type == TOKEN_STRING {
		stmt.Value = &StringLiteral{Value: p.curToken.Literal}
	} else {
		stmt.Value = &StringLiteral{Value: p.curToken.Literal}
	}

	return stmt
}

// parseCommandStatement 解析命令语句
func (p *Parser) parseCommandStatement() Statement {
	stmt := &CommandStatement{
		Command: p.curToken.Literal,
		Args:    []string{},
	}

	p.nextToken()

	// 读取参数
	for p.curToken.Type != TOKEN_NEWLINE &&
		p.curToken.Type != TOKEN_SEMICOLON &&
		p.curToken.Type != TOKEN_EOF &&
		p.curToken.Type != TOKEN_PIPE &&
		!isBlockEnd(p.curToken.Type) {
		stmt.Args = append(stmt.Args, p.curToken.Literal)
		p.nextToken()
	}

	return stmt
}

// parseCondition 解析条件表达式
func (p *Parser) parseCondition() Expression {
	// 简单实现：支持 [ condition ] 形式
	if p.curToken.Type == TOKEN_LBRACKET {
		p.nextToken()

		// 收集测试表达式的所有部分
		args := []string{}
		for p.curToken.Type != TOKEN_RBRACKET && p.curToken.Type != TOKEN_EOF {
			args = append(args, p.curToken.Literal)
			p.nextToken()
		}

		if p.curToken.Type != TOKEN_RBRACKET {
			return &StringLiteral{Value: "false"}
		}

		// 解析测试表达式
		if len(args) >= 2 {
			return &TestExpr{
				Operator: args[0],
				Args:     args[1:],
			}
		}

		return &StringLiteral{Value: strings.Join(args, " ")}
	}

	// 否则作为字符串条件
	cond := p.curToken.Literal
	p.nextToken()
	return &StringLiteral{Value: cond}
}

// parseBlock 解析语句块
func (p *Parser) parseBlock(endTokens ...TokenType) []Statement {
	statements := []Statement{}

	for !contains(endTokens, p.curToken.Type) && p.curToken.Type != TOKEN_EOF {
		p.skipNewlines()

		if contains(endTokens, p.curToken.Type) {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}

		p.nextToken()
	}

	return statements
}

// isBlockEnd 判断是否是块结束 token
func isBlockEnd(t TokenType) bool {
	return t == TOKEN_FI || t == TOKEN_DONE || t == TOKEN_RBRACE ||
		t == TOKEN_ELIF || t == TOKEN_ELSE
}

// contains 判断切片是否包含元素
func contains(slice []TokenType, item TokenType) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
