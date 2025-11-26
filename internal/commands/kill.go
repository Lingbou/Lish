package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

type KillCommand struct {
	stdout *os.File
}

func NewKillCommand(stdout *os.File) *KillCommand {
	return &KillCommand{stdout: stdout}
}

func (c *KillCommand) Name() string {
	return "kill"
}

func (c *KillCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("kill: 需要指定进程 ID")
	}
	
	for _, arg := range args {
		pid, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "kill: 无效的进程 ID: %s\n", arg)
			continue
		}
		
		if err := c.killProcess(pid); err != nil {
			fmt.Fprintf(os.Stderr, "kill: %v\n", err)
		} else {
			fmt.Fprintf(c.stdout, "已终止进程 %d\n", pid)
		}
	}
	
	return nil
}

func (c *KillCommand) killProcess(pid int) error {
	if runtime.GOOS == "windows" {
		// Windows 使用 taskkill
		cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
		return cmd.Run()
	}
	
	// Unix/Linux
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("查找进程失败: %w", err)
	}
	
	return process.Kill()
}

func (c *KillCommand) Help() string {
	return `kill - 终止进程

用法:
  kill PID...

描述:
  终止指定的进程。需要进程 ID (PID)。
  可以使用 ps 命令查看进程列表。

示例:
  ps | grep chrome      # 查找进程 ID
  kill 1234             # 终止进程
  kill 1234 5678        # 终止多个进程`
}

func (c *KillCommand) ShortHelp() string {
	return "终止进程"
}

