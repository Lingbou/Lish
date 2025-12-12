package commands

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

// AwkCommand awk 命令 - 文本处理（基础版）
type AwkCommand struct{}

func (c *AwkCommand) Name() string {
	return "awk"
}

func (c *AwkCommand) Execute(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("awk", flag.ContinueOnError)
	fieldSep := flags.StringP("field-separator", "F", " ", "字段分隔符")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remaining := flags.Args()

	if len(remaining) == 0 {
		return fmt.Errorf("用法: awk [-F sep] 'pattern { action }' [file...]")
	}

	// 解析awk脚本
	script := remaining[0]
	filenames := remaining[1:]

	program, err := parseAwkProgram(script)
	if err != nil {
		return err
	}

	// 读取输入
	var lines []string

	if len(filenames) == 0 {
		// 从标准输入读取
		lines, err = readLines(os.Stdin)
		if err != nil {
			return fmt.Errorf("读取输入失败: %w", err)
		}
	} else {
		// 从文件读取
		for _, filename := range filenames {
			fileLines, err := readLinesFromFile(filename)
			if err != nil {
				return fmt.Errorf("读取文件 %s 失败: %w", filename, err)
			}
			lines = append(lines, fileLines...)
		}
	}

	// 执行awk程序
	return executeAwkProgram(program, lines, *fieldSep)
}

// awkProgram awk程序结构
type awkProgram struct {
	pattern string   // 匹配模式（可选）
	actions []string // 操作列表
}

// parseAwkProgram 解析awk程序
func parseAwkProgram(script string) (*awkProgram, error) {
	prog := &awkProgram{}

	script = strings.TrimSpace(script)

	// 检查是否有模式
	if strings.Contains(script, "{") {
		// 有 pattern { action } 形式
		parts := strings.SplitN(script, "{", 2)
		prog.pattern = strings.TrimSpace(parts[0])

		if len(parts) > 1 {
			actionPart := strings.TrimSuffix(parts[1], "}")
			prog.actions = parseActions(actionPart)
		}
	} else {
		// 只有 action，没有 pattern
		prog.actions = parseActions(script)
	}

	return prog, nil
}

// parseActions 解析操作
func parseActions(actionStr string) []string {
	// 简单分割，按分号或换行分割
	actions := []string{}
	for _, action := range strings.Split(actionStr, ";") {
		action = strings.TrimSpace(action)
		if action != "" {
			actions = append(actions, action)
		}
	}
	return actions
}

// awkContext awk执行上下文
type awkContext struct {
	fields []string          // 当前行的字段
	NR     int               // 行号
	NF     int               // 字段数
	FS     string            // 字段分隔符
	vars   map[string]string // 变量
}

// executeAwkProgram 执行awk程序
func executeAwkProgram(prog *awkProgram, lines []string, fieldSep string) error {
	ctx := &awkContext{
		FS:   fieldSep,
		vars: make(map[string]string),
	}

	for lineNum, line := range lines {
		ctx.NR = lineNum + 1

		// 分割字段
		if fieldSep == " " {
			ctx.fields = strings.Fields(line)
		} else {
			ctx.fields = strings.Split(line, fieldSep)
		}
		ctx.NF = len(ctx.fields)

		// 检查模式匹配
		if prog.pattern != "" {
			matched, err := matchPattern(prog.pattern, line, ctx)
			if err != nil {
				return err
			}
			if !matched {
				continue
			}
		}

		// 执行操作
		for _, action := range prog.actions {
			if err := executeAction(action, line, ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

// matchPattern 匹配模式
func matchPattern(pattern string, line string, ctx *awkContext) (bool, error) {
	// 正则表达式模式
	if strings.HasPrefix(pattern, "/") && strings.HasSuffix(pattern, "/") {
		regexStr := strings.Trim(pattern, "/")
		regex, err := regexp.Compile(regexStr)
		if err != nil {
			return false, err
		}
		return regex.MatchString(line), nil
	}

	// 条件表达式（简单实现）
	// 支持 $1 == "value" 这样的比较
	if strings.Contains(pattern, "==") {
		parts := strings.Split(pattern, "==")
		if len(parts) == 2 {
			left := evaluateExpression(strings.TrimSpace(parts[0]), ctx)
			right := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
			return left == right, nil
		}
	}

	return true, nil
}

// executeAction 执行操作
func executeAction(action string, line string, ctx *awkContext) error {
	action = strings.TrimSpace(action)

	// print 语句
	if strings.HasPrefix(action, "print") {
		return executePrint(action, line, ctx)
	}

	// printf 语句
	if strings.HasPrefix(action, "printf") {
		return executePrintf(action, ctx)
	}

	// 变量赋值
	if strings.Contains(action, "=") && !strings.Contains(action, "==") {
		return executeAssignment(action, ctx)
	}

	return fmt.Errorf("不支持的操作: %s", action)
}

// executePrint 执行print语句
func executePrint(action string, line string, ctx *awkContext) error {
	// 去掉 "print" 关键字
	action = strings.TrimSpace(strings.TrimPrefix(action, "print"))

	if action == "" {
		// 打印整行
		fmt.Println(line)
		return nil
	}

	// 解析要打印的内容
	parts := strings.Split(action, ",")
	var output []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		value := evaluateExpression(part, ctx)
		output = append(output, value)
	}

	fmt.Println(strings.Join(output, " "))
	return nil
}

// executePrintf 执行printf语句
func executePrintf(action string, ctx *awkContext) error {
	// 简单实现，暂不支持复杂格式化
	action = strings.TrimSpace(strings.TrimPrefix(action, "printf"))
	action = strings.Trim(action, "()")

	// 分割格式字符串和参数
	parts := strings.SplitN(action, ",", 2)
	if len(parts) < 1 {
		return fmt.Errorf("无效的printf语句")
	}

	format := strings.Trim(parts[0], "\"'")

	if len(parts) > 1 {
		args := strings.Split(parts[1], ",")
		var values []interface{}
		for _, arg := range args {
			values = append(values, evaluateExpression(strings.TrimSpace(arg), ctx))
		}
		fmt.Printf(format, values...)
	} else {
		fmt.Print(format)
	}

	return nil
}

// executeAssignment 执行赋值语句
func executeAssignment(action string, ctx *awkContext) error {
	parts := strings.Split(action, "=")
	if len(parts) != 2 {
		return fmt.Errorf("无效的赋值语句: %s", action)
	}

	varName := strings.TrimSpace(parts[0])
	value := evaluateExpression(strings.TrimSpace(parts[1]), ctx)
	ctx.vars[varName] = value

	return nil
}

// evaluateExpression 计算表达式
func evaluateExpression(expr string, ctx *awkContext) string {
	expr = strings.TrimSpace(expr)

	// $0 - 整行
	if expr == "$0" {
		return strings.Join(ctx.fields, ctx.FS)
	}

	// $1, $2, ... - 字段
	if strings.HasPrefix(expr, "$") {
		fieldNumStr := strings.TrimPrefix(expr, "$")
		fieldNum, err := strconv.Atoi(fieldNumStr)
		if err == nil && fieldNum > 0 && fieldNum <= ctx.NF {
			return ctx.fields[fieldNum-1]
		}
		return ""
	}

	// NR - 行号
	if expr == "NR" {
		return strconv.Itoa(ctx.NR)
	}

	// NF - 字段数
	if expr == "NF" {
		return strconv.Itoa(ctx.NF)
	}

	// FS - 字段分隔符
	if expr == "FS" {
		return ctx.FS
	}

	// 变量
	if val, ok := ctx.vars[expr]; ok {
		return val
	}

	// 字符串字面量
	if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") {
		return strings.Trim(expr, "\"")
	}
	if strings.HasPrefix(expr, "'") && strings.HasSuffix(expr, "'") {
		return strings.Trim(expr, "'")
	}

	// 简单算术表达式
	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		if len(parts) == 2 {
			a := evaluateExpression(parts[0], ctx)
			b := evaluateExpression(parts[1], ctx)
			na, _ := strconv.Atoi(a)
			nb, _ := strconv.Atoi(b)
			return strconv.Itoa(na + nb)
		}
	}

	// 默认返回原始表达式
	return expr
}

func (c *AwkCommand) Help() string {
	return `awk - 文本处理工具（基础版）

用法:
  awk [-F sep] 'pattern { action }' [file...]
  command | awk 'pattern { action }'

说明:
  强大的文本处理工具，支持模式匹配和字段处理。

选项:
  -F, --field-separator   字段分隔符（默认为空格）

内置变量:
  $0          当前行的完整内容
  $1, $2,...  第1个、第2个字段
  NR          当前行号
  NF          当前行的字段数
  FS          字段分隔符

支持的操作:
  print           打印整行
  print $1        打印第1个字段
  print $1, $2    打印多个字段
  printf "fmt"    格式化打印

示例:
  awk '{ print $1 }' file.txt              # 打印第1列
  awk '{ print $1, $3 }' file.txt          # 打印第1和第3列
  awk '{ print NR, $0 }' file.txt          # 打印行号和内容
  awk 'NR > 1 { print $1 }' file.txt       # 跳过第一行
  awk -F: '{ print $1 }' /etc/passwd       # 使用:分隔
  awk '/pattern/ { print $0 }' file.txt    # 只打印匹配行
  ls -l | awk '{ print $9 }'               # 提取文件名
  awk '{ sum += $1 } END { print sum }'    # 求和（未完全支持）`
}

func (c *AwkCommand) ShortHelp() string {
	return "文本处理和字段提取"
}
