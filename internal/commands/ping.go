package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/spf13/pflag"
)

type PingCommand struct {
	stdout *os.File
}

func NewPingCommand(stdout *os.File) *PingCommand {
	return &PingCommand{stdout: stdout}
}

func (c *PingCommand) Name() string {
	return "ping"
}

func (c *PingCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("ping", pflag.ContinueOnError)
	count := flags.IntP("count", "c", 4, "发送的数据包数量")

	if err := flags.Parse(args); err != nil {
		return err
	}

	hosts := flags.Args()
	if len(hosts) == 0 {
		return fmt.Errorf("ping: 需要指定主机名或 IP 地址")
	}

	host := hosts[0]

	// 构建 ping 命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", strconv.Itoa(*count), host)
	} else {
		cmd = exec.Command("ping", "-c", strconv.Itoa(*count), host)
	}

	// 设置输出
	cmd.Stdout = c.stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ping 失败: %w", err)
	}

	return nil
}

func (c *PingCommand) Help() string {
	return `ping - 测试网络连接

用法:
  ping [选项] 主机名或IP

选项:
  -c, --count    发送的数据包数量（默认 4）

描述:
  向目标主机发送 ICMP 请求，测试网络连通性。

示例:
  ping google.com           # 默认发送 4 个包
  ping -c 10 8.8.8.8       # 发送 10 个包
  ping baidu.com           # 测试国内网络`
}

func (c *PingCommand) ShortHelp() string {
	return "测试网络连接"
}
