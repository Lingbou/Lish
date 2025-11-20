package commands

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
)

type EnvCommand struct {
	stdout *os.File
}

func NewEnvCommand(stdout *os.File) *EnvCommand {
	return &EnvCommand{stdout: stdout}
}

func (c *EnvCommand) Name() string {
	return "env"
}

func (c *EnvCommand) Execute(ctx context.Context, args []string) error {
	// 如果没有参数，显示所有环境变量
	if len(args) == 0 {
		return c.listEnv()
	}

	// 设置或显示特定环境变量
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			// 设置环境变量 KEY=VALUE
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
				fmt.Fprintf(c.stdout, "设置环境变量: %s=%s\n", parts[0], parts[1])
			}
		} else {
			// 显示特定环境变量
			value := os.Getenv(arg)
			if value != "" {
				fmt.Fprintf(c.stdout, "%s=%s\n", arg, value)
			} else {
				fmt.Fprintf(c.stdout, "%s: 未设置\n", arg)
			}
		}
	}

	return nil
}

func (c *EnvCommand) listEnv() error {
	env := os.Environ()
	sort.Strings(env)

	for _, e := range env {
		fmt.Fprintln(c.stdout, e)
	}

	return nil
}

func (c *EnvCommand) Help() string {
	return `env - 环境变量管理

用法:
  env                    # 显示所有环境变量
  env VAR                # 显示特定环境变量
  env VAR=VALUE          # 设置环境变量

描述:
  查看和设置环境变量。

注意:
  使用 env 设置的环境变量只在当前 Shell 会话中有效。

示例:
  env                    # 显示所有环境变量
  env PATH               # 显示 PATH 变量
  env MY_VAR=hello       # 设置环境变量
  env MY_VAR             # 查看刚设置的变量`
}

func (c *EnvCommand) ShortHelp() string {
	return "环境变量管理"
}
