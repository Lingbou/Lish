package script

// Node 是所有 AST 节点的基础接口
type Node interface {
	String() string
}

// Statement 表示一个语句
type Statement interface {
	Node
	statementNode()
}

// Expression 表示一个表达式
type Expression interface {
	Node
	expressionNode()
}

// ============================================================================
// 程序和语句
// ============================================================================

// Program 表示整个脚本程序
type Program struct {
	Statements []Statement
}

func (p *Program) String() string { return "Program" }

// CommandStatement 表示一个命令语句（如: ls -la）
type CommandStatement struct {
	Command string
	Args    []string
}

func (cs *CommandStatement) statementNode() {}
func (cs *CommandStatement) String() string { return "Command: " + cs.Command }

// AssignStatement 表示变量赋值语句（如: name=value）
type AssignStatement struct {
	Name  string
	Value Expression
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) String() string { return "Assign: " + as.Name }

// IfStatement 表示 if 条件语句
type IfStatement struct {
	Condition  Expression
	ThenBlock  []Statement
	ElseIfList []*ElseIfClause
	ElseBlock  []Statement
}

type ElseIfClause struct {
	Condition Expression
	Block     []Statement
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string { return "If" }

// ForStatement 表示 for 循环语句
type ForStatement struct {
	Variable string      // 循环变量名
	List     []string    // 要遍历的列表
	Block    []Statement // 循环体
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string { return "For: " + fs.Variable }

// WhileStatement 表示 while 循环语句
type WhileStatement struct {
	Condition Expression
	Block     []Statement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string { return "While" }

// FunctionDef 表示函数定义
type FunctionDef struct {
	Name  string
	Block []Statement
}

func (fd *FunctionDef) statementNode() {}
func (fd *FunctionDef) String() string { return "Function: " + fd.Name }

// ReturnStatement 表示 return 语句
type ReturnStatement struct {
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string { return "Return" }

// BreakStatement 表示 break 语句
type BreakStatement struct{}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) String() string { return "Break" }

// ContinueStatement 表示 continue 语句
type ContinueStatement struct{}

func (cs *ContinueStatement) statementNode() {}
func (cs *ContinueStatement) String() string { return "Continue" }

// ============================================================================
// 表达式
// ============================================================================

// StringLiteral 表示字符串字面量
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return sl.Value }

// Variable 表示变量引用
type Variable struct {
	Name string
}

func (v *Variable) expressionNode() {}
func (v *Variable) String() string  { return "$" + v.Name }

// BinaryExpr 表示二元表达式（如: a == b, a > b）
type BinaryExpr struct {
	Left     Expression
	Operator string // ==, !=, <, >, <=, >=, &&, ||
	Right    Expression
}

func (be *BinaryExpr) expressionNode() {}
func (be *BinaryExpr) String() string  { return "BinaryExpr" }

// UnaryExpr 表示一元表达式（如: !condition）
type UnaryExpr struct {
	Operator string // !, -
	Operand  Expression
}

func (ue *UnaryExpr) expressionNode() {}
func (ue *UnaryExpr) String() string  { return "UnaryExpr" }

// TestExpr 表示测试表达式（如: [ -f file ], [ $a -eq $b ]）
type TestExpr struct {
	Operator string   // -f, -d, -e, -z, -n, -eq, -ne, -lt, -gt, -le, -ge
	Args     []string // 参数列表
}

func (te *TestExpr) expressionNode() {}
func (te *TestExpr) String() string  { return "Test" }

// CommandSubstitution 表示命令替换（如: $(command)）
type CommandSubstitution struct {
	Command string
	Args    []string
}

func (cs *CommandSubstitution) expressionNode() {}
func (cs *CommandSubstitution) String() string  { return "$()" }
