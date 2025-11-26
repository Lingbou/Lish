package commands

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Lingbou/Lish/internal/config"
)

type AliasCommand struct {
	stdout *os.File
	config *config.Config
}

func NewAliasCommand(stdout *os.File, cfg *config.Config) *AliasCommand {
	return &AliasCommand{
		stdout: stdout,
		config: cfg,
	}
}

func (c *AliasCommand) Name() string {
	return "alias"
}

func (c *AliasCommand) Execute(ctx context.Context, args []string) error {
	// 如果没有参数，显示所有别名
	if len(args) == 0 {
		return c.listAliases()
	}

	// 设置别名
	for _, arg := range args {
		if !strings.Contains(arg, "=") {
			// 显示特定别名
			if cmd, exists := c.config.GetAlias(arg); exists {
				fmt.Fprintf(c.stdout, "alias %s='%s'\n", arg, cmd)
			} else {
				fmt.Fprintf(c.stdout, "alias: %s: 未定义\n", arg)
			}
			continue
		}

		// 解析 name=command
		parts := strings.SplitN(arg, "=", 2)
		name := strings.TrimSpace(parts[0])
		command := strings.Trim(strings.TrimSpace(parts[1]), "'\"")

		// 设置别名
		c.config.SetAlias(name, command)
		fmt.Fprintf(c.stdout, "设置别名: %s='%s'\n", name, command)
	}

	// 保存配置
	if err := c.config.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	return nil
}

func (c *AliasCommand) listAliases() error {
	if len(c.config.Aliases) == 0 {
		fmt.Fprintln(c.stdout, "没有定义的别名")
		return nil
	}

	// 排序显示
	names := make([]string, 0, len(c.config.Aliases))
	for name := range c.config.Aliases {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		cmd := c.config.Aliases[name]
		fmt.Fprintf(c.stdout, "alias %s='%s'\n", name, cmd)
	}

	return nil
}

func (c *AliasCommand) Help() string {
	return `alias - 命令别名管理

用法:
  alias                  # 显示所有别名
  alias name             # 显示特定别名
  alias name='command'   # 设置别名

描述:
  创建命令别名，简化常用命令的输入。
  别名会自动保存到配置文件 ~/.lishrc 中。

示例:
  alias                  # 列出所有别名
  alias ll='ls -l'       # 设置别名
  alias la='ls -la'      # 另一个别名
  alias ll               # 查看 ll 别名`
}

func (c *AliasCommand) ShortHelp() string {
	return "命令别名管理"
}
