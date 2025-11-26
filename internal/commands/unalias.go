package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/Lingbou/Lish/internal/config"
)

type UnaliasCommand struct {
	stdout *os.File
	config *config.Config
}

func NewUnaliasCommand(stdout *os.File, cfg *config.Config) *UnaliasCommand {
	return &UnaliasCommand{
		stdout: stdout,
		config: cfg,
	}
}

func (c *UnaliasCommand) Name() string {
	return "unalias"
}

func (c *UnaliasCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("unalias: 需要指定别名名称")
	}

	for _, name := range args {
		if _, exists := c.config.GetAlias(name); !exists {
			fmt.Fprintf(c.stdout, "unalias: %s: 未定义\n", name)
			continue
		}

		c.config.RemoveAlias(name)
		fmt.Fprintf(c.stdout, "已删除别名: %s\n", name)
	}

	// 保存配置
	if err := c.config.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	return nil
}

func (c *UnaliasCommand) Help() string {
	return `unalias - 删除命令别名

用法:
  unalias 别名...

描述:
  删除一个或多个命令别名。

示例:
  unalias ll             # 删除 ll 别名
  unalias ll la          # 删除多个别名`
}

func (c *UnaliasCommand) ShortHelp() string {
	return "删除命令别名"
}
