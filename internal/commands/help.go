package commands

import (
	"context"
	"fmt"
	"os"
	"sort"
)

type HelpCommand struct {
	registry *Registry
	stdout   *os.File
}

func NewHelpCommand(registry *Registry, stdout *os.File) *HelpCommand {
	return &HelpCommand{
		registry: registry,
		stdout:   stdout,
	}
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Execute(ctx context.Context, args []string) error {
	if len(args) > 0 {
		// 显示特定命令的帮助
		cmdName := args[0]
		cmd, exists := c.registry.Get(cmdName)
		if !exists {
			return fmt.Errorf("未知命令: %s", cmdName)
		}
		fmt.Fprintln(c.stdout, cmd.Help())
		return nil
	}
	
	// 显示所有命令列表
	fmt.Fprintln(c.stdout, "Lish - Linux 风格的轻量级 Shell")
	fmt.Fprintln(c.stdout, "\n可用命令:")
	
	commands := c.registry.GetAll()
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	
	for _, name := range names {
		cmd := commands[name]
		fmt.Fprintf(c.stdout, "  %-12s %s\n", name, cmd.ShortHelp())
	}
	
	fmt.Fprintln(c.stdout, "\n输入 'help <命令>' 查看详细帮助信息")
	
	return nil
}

func (c *HelpCommand) Help() string {
	return `help - 显示帮助信息

用法:
  help [命令]

描述:
  显示命令的帮助信息。不带参数时显示所有可用命令列表。

示例:
  help         # 显示所有命令
  help ls      # 显示 ls 命令的详细帮助`
}

func (c *HelpCommand) ShortHelp() string {
	return "显示帮助信息"
}

