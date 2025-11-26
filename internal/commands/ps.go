package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type PsCommand struct {
	stdout *os.File
}

func NewPsCommand(stdout *os.File) *PsCommand {
	return &PsCommand{stdout: stdout}
}

func (c *PsCommand) Name() string {
	return "ps"
}

func (c *PsCommand) Execute(ctx context.Context, args []string) error {
	// Windows 使用 tasklist，Linux 使用 ps
	if runtime.GOOS == "windows" {
		return c.windowsPs()
	}
	return c.unixPs()
}

func (c *PsCommand) windowsPs() error {
	cmd := exec.Command("tasklist")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行 tasklist 失败: %w", err)
	}

	fmt.Fprint(c.stdout, string(output))
	return nil
}

func (c *PsCommand) unixPs() error {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行 ps 失败: %w", err)
	}

	fmt.Fprint(c.stdout, string(output))
	return nil
}

func (c *PsCommand) Help() string {
	var examples string
	if runtime.GOOS == "windows" {
		examples = `示例:
  ps                    # 显示所有进程
  ps | grep chrome      # 查找 Chrome 进程`
	} else {
		examples = `示例:
  ps                    # 显示所有进程
  ps | grep firefox     # 查找 Firefox 进程`
	}

	return fmt.Sprintf(`ps - 显示进程列表

用法:
  ps

描述:
  显示当前运行的进程列表。
  在 Windows 上使用 tasklist，在 Linux 上使用 ps aux。

%s`, examples)
}

func (c *PsCommand) ShortHelp() string {
	return "显示进程列表"
}
