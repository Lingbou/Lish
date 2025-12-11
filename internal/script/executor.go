package script

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ExecutionError 执行错误
type ExecutionError struct {
	Message string
	Line    int
}

func (e *ExecutionError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("行 %d: %s", e.Line, e.Message)
	}
	return e.Message
}

// ControlFlow 控制流类型
type ControlFlow int

const (
	FLOW_NORMAL ControlFlow = iota
	FLOW_BREAK
	FLOW_CONTINUE
	FLOW_RETURN
)

// CommandExecutor 命令执行器接口
type CommandExecutor interface {
	ExecuteCommand(ctx context.Context, command string, args []string) error
}

// Executor 脚本执行器
type Executor struct {
	cmdExecutor CommandExecutor
	variables   *VariableManager
	functions   map[string]*FunctionDef
	lastExit    int
	flowType    ControlFlow
	returnVal   string
}

// NewExecutor 创建新的执行器
func NewExecutor(cmdExecutor CommandExecutor) *Executor {
	return &Executor{
		cmdExecutor: cmdExecutor,
		variables:   NewVariableManager(),
		functions:   make(map[string]*FunctionDef),
		lastExit:    0,
		flowType:    FLOW_NORMAL,
	}
}

// Execute 执行脚本
func (e *Executor) Execute(ctx context.Context, program *Program) error {
	for _, stmt := range program.Statements {
		if err := e.executeStatement(ctx, stmt); err != nil {
			return err
		}

		// 检查控制流
		if e.flowType == FLOW_RETURN {
			break
		}
	}

	return nil
}

// ExecuteFile 执行脚本文件
func (e *Executor) ExecuteFile(ctx context.Context, filepath string, args []string) error {
	// 读取文件
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("无法读取脚本文件: %w", err)
	}

	// 设置特殊变量
	scriptArgs := append([]string{filepath}, args...)
	e.variables.SetSpecialVars(scriptArgs, e.lastExit)

	// 解析
	parser := NewParser(string(content))
	program := parser.Parse()

	if len(parser.Errors()) > 0 {
		return fmt.Errorf("解析错误: %s", strings.Join(parser.Errors(), "\n"))
	}

	// 执行
	return e.Execute(ctx, program)
}

// executeStatement 执行语句
func (e *Executor) executeStatement(ctx context.Context, stmt Statement) error {
	switch s := stmt.(type) {
	case *CommandStatement:
		return e.executeCommand(ctx, s)
	case *AssignStatement:
		return e.executeAssign(ctx, s)
	case *IfStatement:
		return e.executeIf(ctx, s)
	case *ForStatement:
		return e.executeFor(ctx, s)
	case *WhileStatement:
		return e.executeWhile(ctx, s)
	case *FunctionDef:
		return e.executeFunctionDef(s)
	case *ReturnStatement:
		return e.executeReturn(ctx, s)
	case *BreakStatement:
		e.flowType = FLOW_BREAK
		return nil
	case *ContinueStatement:
		e.flowType = FLOW_CONTINUE
		return nil
	default:
		return fmt.Errorf("未知语句类型: %T", stmt)
	}
}

// executeCommand 执行命令
func (e *Executor) executeCommand(ctx context.Context, stmt *CommandStatement) error {
	// 展开变量
	command := e.variables.Expand(stmt.Command)
	args := make([]string, len(stmt.Args))
	for i, arg := range stmt.Args {
		args[i] = e.variables.Expand(arg)
	}

	// 检查是否是函数调用
	if fn, ok := e.functions[command]; ok {
		return e.executeFunction(ctx, fn, args)
	}

	// 使用命令执行器执行命令
	err := e.cmdExecutor.ExecuteCommand(ctx, command, args)
	if err != nil {
		e.lastExit = 1
		return err
	}

	e.lastExit = 0
	e.variables.Set("?", "0")
	return nil
}

// executeAssign 执行赋值
func (e *Executor) executeAssign(ctx context.Context, stmt *AssignStatement) error {
	value := e.evaluateExpression(ctx, stmt.Value)
	e.variables.Set(stmt.Name, value)
	return nil
}

// executeIf 执行 if 语句
func (e *Executor) executeIf(ctx context.Context, stmt *IfStatement) error {
	// 评估条件
	if e.evaluateCondition(ctx, stmt.Condition) {
		return e.executeBlock(ctx, stmt.ThenBlock)
	}

	// 检查 elif
	for _, elseif := range stmt.ElseIfList {
		if e.evaluateCondition(ctx, elseif.Condition) {
			return e.executeBlock(ctx, elseif.Block)
		}
	}

	// 执行 else 块
	if stmt.ElseBlock != nil {
		return e.executeBlock(ctx, stmt.ElseBlock)
	}

	return nil
}

// executeFor 执行 for 循环
func (e *Executor) executeFor(ctx context.Context, stmt *ForStatement) error {
	for _, item := range stmt.List {
		// 展开变量
		item = e.variables.Expand(item)

		// 设置循环变量
		e.variables.Set(stmt.Variable, item)

		// 执行循环体
		if err := e.executeBlock(ctx, stmt.Block); err != nil {
			return err
		}

		// 检查控制流
		if e.flowType == FLOW_BREAK {
			e.flowType = FLOW_NORMAL
			break
		}
		if e.flowType == FLOW_CONTINUE {
			e.flowType = FLOW_NORMAL
			continue
		}
		if e.flowType == FLOW_RETURN {
			return nil
		}
	}

	return nil
}

// executeWhile 执行 while 循环
func (e *Executor) executeWhile(ctx context.Context, stmt *WhileStatement) error {
	for e.evaluateCondition(ctx, stmt.Condition) {
		if err := e.executeBlock(ctx, stmt.Block); err != nil {
			return err
		}

		// 检查控制流
		if e.flowType == FLOW_BREAK {
			e.flowType = FLOW_NORMAL
			break
		}
		if e.flowType == FLOW_CONTINUE {
			e.flowType = FLOW_NORMAL
			continue
		}
		if e.flowType == FLOW_RETURN {
			return nil
		}
	}

	return nil
}

// executeFunctionDef 注册函数定义
func (e *Executor) executeFunctionDef(stmt *FunctionDef) error {
	e.functions[stmt.Name] = stmt
	return nil
}

// executeFunction 执行函数调用
func (e *Executor) executeFunction(ctx context.Context, fn *FunctionDef, args []string) error {
	// 进入新作用域
	e.variables.PushScope()
	defer e.variables.PopScope()

	// 设置参数
	fnArgs := append([]string{fn.Name}, args...)
	e.variables.SetSpecialVars(fnArgs, e.lastExit)

	// 执行函数体
	if err := e.executeBlock(ctx, fn.Block); err != nil {
		return err
	}

	// 重置控制流
	if e.flowType == FLOW_RETURN {
		e.flowType = FLOW_NORMAL
	}

	return nil
}

// executeReturn 执行 return 语句
func (e *Executor) executeReturn(ctx context.Context, stmt *ReturnStatement) error {
	if stmt.Value != nil {
		e.returnVal = e.evaluateExpression(ctx, stmt.Value)
	}
	e.flowType = FLOW_RETURN
	return nil
}

// executeBlock 执行语句块
func (e *Executor) executeBlock(ctx context.Context, block []Statement) error {
	for _, stmt := range block {
		if err := e.executeStatement(ctx, stmt); err != nil {
			return err
		}

		// 检查控制流
		if e.flowType != FLOW_NORMAL {
			break
		}
	}

	return nil
}

// evaluateCondition 评估条件表达式
func (e *Executor) evaluateCondition(ctx context.Context, expr Expression) bool {
	switch exp := expr.(type) {
	case *StringLiteral:
		// 非空字符串为真
		val := e.variables.Expand(exp.Value)
		return val != "" && val != "0" && val != "false"

	case *TestExpr:
		return e.evaluateTest(exp)

	case *BinaryExpr:
		return e.evaluateBinaryExpr(exp)

	case *UnaryExpr:
		return e.evaluateUnaryExpr(exp)

	default:
		return false
	}
}

// evaluateTest 评估测试表达式
func (e *Executor) evaluateTest(expr *TestExpr) bool {
	if len(expr.Args) == 0 {
		return false
	}

	op := expr.Operator
	args := expr.Args

	// 展开变量
	for i := range args {
		args[i] = e.variables.Expand(args[i])
	}

	switch op {
	case "-f": // 文件存在且是普通文件
		if len(args) < 1 {
			return false
		}
		info, err := os.Stat(args[0])
		return err == nil && !info.IsDir()

	case "-d": // 目录存在
		if len(args) < 1 {
			return false
		}
		info, err := os.Stat(args[0])
		return err == nil && info.IsDir()

	case "-e": // 文件或目录存在
		if len(args) < 1 {
			return false
		}
		_, err := os.Stat(args[0])
		return err == nil

	case "-z": // 字符串为空
		if len(args) < 1 {
			return true
		}
		return args[0] == ""

	case "-n": // 字符串不为空
		if len(args) < 1 {
			return false
		}
		return args[0] != ""

	case "-eq", "-ne", "-lt", "-gt", "-le", "-ge":
		if len(args) < 2 {
			return false
		}
		return e.compareNumbers(args[0], op, args[1])

	case "==", "!=":
		if len(args) < 2 {
			return false
		}
		if op == "==" {
			return args[0] == args[1]
		}
		return args[0] != args[1]

	default:
		// 默认：非空为真
		return len(args) > 0 && args[0] != ""
	}
}

// compareNumbers 比较数字
func (e *Executor) compareNumbers(a, op, b string) bool {
	n1, err1 := strconv.Atoi(a)
	n2, err2 := strconv.Atoi(b)

	if err1 != nil || err2 != nil {
		return false
	}

	switch op {
	case "-eq":
		return n1 == n2
	case "-ne":
		return n1 != n2
	case "-lt":
		return n1 < n2
	case "-gt":
		return n1 > n2
	case "-le":
		return n1 <= n2
	case "-ge":
		return n1 >= n2
	default:
		return false
	}
}

// evaluateBinaryExpr 评估二元表达式
func (e *Executor) evaluateBinaryExpr(expr *BinaryExpr) bool {
	left := e.evaluateCondition(context.Background(), expr.Left)

	switch expr.Operator {
	case "&&":
		if !left {
			return false
		}
		return e.evaluateCondition(context.Background(), expr.Right)
	case "||":
		if left {
			return true
		}
		return e.evaluateCondition(context.Background(), expr.Right)
	default:
		return false
	}
}

// evaluateUnaryExpr 评估一元表达式
func (e *Executor) evaluateUnaryExpr(expr *UnaryExpr) bool {
	result := e.evaluateCondition(context.Background(), expr.Operand)

	if expr.Operator == "!" {
		return !result
	}

	return result
}

// evaluateExpression 评估表达式并返回字符串值
func (e *Executor) evaluateExpression(ctx context.Context, expr Expression) string {
	switch exp := expr.(type) {
	case *StringLiteral:
		return e.variables.Expand(exp.Value)
	case *Variable:
		val, _ := e.variables.Get(exp.Name)
		return val
	default:
		return ""
	}
}

// GetVariable 获取变量值
func (e *Executor) GetVariable(name string) (string, bool) {
	return e.variables.Get(name)
}

// SetVariable 设置变量值
func (e *Executor) SetVariable(name, value string) {
	e.variables.Set(name, value)
}

// LastExitCode 获取上一个命令的退出码
func (e *Executor) LastExitCode() int {
	return e.lastExit
}
