package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/Lingbou/Lish/internal/script"
	flag "github.com/spf13/pflag"
)

// SourceCommand source 命令 - 在当前环境执行脚本
type SourceCommand struct {
	executor *script.Executor
}

// NewSourceCommand 创建 source 命令
func NewSourceCommand(executor *script.Executor) *SourceCommand {
	return &SourceCommand{
		executor: executor,
	}
}

func (c *SourceCommand) Name() string {
	return "source"
}

func (c *SourceCommand) Execute(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("source", flag.ContinueOnError)
	verbose := flags.BoolP("verbose", "v", false, "详细模式")
	_ = flags.BoolP("debug", "x", false, "调试模式") // 保留用于未来实现

	if err := flags.Parse(args); err != nil {
		return err
	}

	remaining := flags.Args()
	if len(remaining) < 1 {
		return fmt.Errorf("用法: source [-v] [-x] <script>")
	}

	scriptFile := remaining[0]
	scriptArgs := remaining[1:]

	// 检查文件是否存在
	if _, err := os.Stat(scriptFile); os.IsNotExist(err) {
		return fmt.Errorf("脚本文件不存在: %s", scriptFile)
	}

	// 显示调试信息
	if *verbose {
		fmt.Printf("执行脚本: %s\n", scriptFile)
		if len(scriptArgs) > 0 {
			fmt.Printf("参数: %v\n", scriptArgs)
		}
	}

	// 执行脚本
	if err := c.executor.ExecuteFile(ctx, scriptFile, scriptArgs); err != nil {
		return fmt.Errorf("脚本执行失败: %w", err)
	}

	if *verbose {
		fmt.Printf("✓ 脚本执行完成\n")
	}

	// debug 模式在执行器内部处理

	return nil
}

func (c *SourceCommand) Help() string {
	return `source - 在当前环境执行脚本

用法:
  source [-v] [-x] <script> [args...]
  . <script> [args...]

说明:
  在当前 shell 环境中执行 .lish 脚本文件。
  脚本中定义的变量和函数会在当前环境中保留。

选项:
  -v, --verbose    详细模式，显示执行过程
  -x, --debug      调试模式，显示每条命令

示例:
  source script.lish              # 执行脚本
  source script.lish arg1 arg2    # 带参数执行
  . ~/.lishrc.lish                # 使用别名执行
  source -v script.lish           # 详细模式

脚本语法支持:
  - 变量: name=value, $name
  - 条件: if [ condition ]; then ... fi
  - 循环: for item in list; do ... done
  - 循环: while [ condition ]; do ... done
  - 函数: function name() { ... }
  - 控制: break, continue, return`
}

func (c *SourceCommand) ShortHelp() string {
	return "在当前环境执行脚本"
}

// ExecCommand exec 命令 - 在新环境执行脚本
type ExecCommand struct {
	cmdExecutor script.CommandExecutor
}

// NewExecCommand 创建 exec 命令
func NewExecCommand(cmdExecutor script.CommandExecutor) *ExecCommand {
	return &ExecCommand{
		cmdExecutor: cmdExecutor,
	}
}

func (c *ExecCommand) Name() string {
	return "exec"
}

func (c *ExecCommand) Execute(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("exec", flag.ContinueOnError)
	verbose := flags.BoolP("verbose", "v", false, "详细模式")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remaining := flags.Args()
	if len(remaining) < 1 {
		return fmt.Errorf("用法: exec [-v] <script> [args...]")
	}

	scriptFile := remaining[0]
	scriptArgs := remaining[1:]

	// 检查文件是否存在
	if _, err := os.Stat(scriptFile); os.IsNotExist(err) {
		return fmt.Errorf("脚本文件不存在: %s", scriptFile)
	}

	// 显示调试信息
	if *verbose {
		fmt.Printf("在新环境中执行脚本: %s\n", scriptFile)
		if len(scriptArgs) > 0 {
			fmt.Printf("参数: %v\n", scriptArgs)
		}
	}

	// 创建新的执行器（独立环境）
	executor := script.NewExecutor(c.cmdExecutor)

	// 执行脚本
	if err := executor.ExecuteFile(ctx, scriptFile, scriptArgs); err != nil {
		return fmt.Errorf("脚本执行失败: %w", err)
	}

	if *verbose {
		fmt.Printf("✓ 脚本执行完成（退出码: %d）\n", executor.LastExitCode())
	}

	return nil
}

func (c *ExecCommand) Help() string {
	return `exec - 在新环境执行脚本

用法:
  exec [-v] <script> [args...]

说明:
  在独立的新环境中执行 .lish 脚本文件。
  脚本中定义的变量和函数不会影响当前环境。

选项:
  -v, --verbose    详细模式，显示执行过程

示例:
  exec script.lish              # 在新环境执行
  exec script.lish arg1 arg2    # 带参数执行
  exec -v script.lish           # 详细模式

区别:
  source  - 在当前环境执行，变量和函数会保留
  exec    - 在新环境执行，不影响当前环境`
}

func (c *ExecCommand) ShortHelp() string {
	return "在新环境执行脚本"
}
