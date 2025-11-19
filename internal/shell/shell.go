package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Lingbou/Lish/internal/commands"
	"github.com/Lingbou/Lish/internal/completer"
	"github.com/Lingbou/Lish/internal/history"
	"github.com/Lingbou/Lish/internal/parser"
	"github.com/chzyer/readline"
)

// Shell Lish Shell 结构
type Shell struct {
	registry *commands.Registry
	history  *history.Manager
	rl       *readline.Instance
	stdout   *os.File
	stderr   *os.File
}

// NewShell 创建新的 Shell 实例
func NewShell() (*Shell, error) {
	// 创建命令注册表
	registry := commands.NewRegistry()

	// 创建历史管理器
	histMgr, err := history.NewManager()
	if err != nil {
		return nil, fmt.Errorf("创建历史管理器失败: %w", err)
	}

	return &Shell{
		registry: registry,
		history:  histMgr,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
	}, nil
}

// Init 初始化 Shell（注册命令、设置 readline）
func (s *Shell) Init() error {
	// 注册所有命令
	if err := s.registerCommands(); err != nil {
		return fmt.Errorf("注册命令失败: %w", err)
	}

	// 创建补全器
	cmdNames := s.registry.List()
	comp := completer.NewCompleter(cmdNames)

	// 配置 readline
	cfg := &readline.Config{
		Prompt:          s.getPrompt(),
		HistoryFile:     s.history.GetHistoryFile(),
		AutoComplete:    comp,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: nil,
	}

	// 创建 readline 实例
	rl, err := readline.NewEx(cfg)
	if err != nil {
		return fmt.Errorf("初始化 readline 失败: %w", err)
	}

	s.rl = rl

	return nil
}

// registerCommands 注册所有内置命令
func (s *Shell) registerCommands() error {
	cmds := []commands.Command{
		commands.NewPwdCommand(s.stdout),
		commands.NewCdCommand(),
		commands.NewLsCommand(s.stdout),
		commands.NewCatCommand(s.stdout),
		commands.NewMkdirCommand(),
		commands.NewRmCommand(),
		commands.NewTouchCommand(),
		commands.NewEchoCommand(s.stdout),
		commands.NewClearCommand(s.stdout),
		commands.NewExitCommand(),
		commands.NewHelpCommand(s.registry, s.stdout),
	}

	for _, cmd := range cmds {
		if err := s.registry.Register(cmd); err != nil {
			return err
		}
	}

	return nil
}

// Run 运行 Shell 主循环
func (s *Shell) Run() error {
	defer s.rl.Close()

	// 显示欢迎信息
	s.printWelcome()

	ctx := context.Background()

	for {
		// 更新提示符
		s.rl.SetPrompt(s.getPrompt())

		// 读取输入
		line, err := s.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Ctrl+C
				continue
			} else if err == io.EOF {
				// Ctrl+D 或 EOF
				fmt.Fprintln(s.stdout, "\n再见!")
				break
			}
			return fmt.Errorf("读取输入失败: %w", err)
		}

		// 解析命令
		parsed, err := parser.Parse(line)
		if err != nil {
			fmt.Fprintf(s.stderr, "解析错误: %v\n", err)
			continue
		}

		if parsed == nil || parsed.Command == "" {
			continue
		}

		// 执行命令
		if err := s.executeCommand(ctx, parsed); err != nil {
			fmt.Fprintf(s.stderr, "错误: %v\n", err)
		}
	}

	return nil
}

// executeCommand 执行解析后的命令
func (s *Shell) executeCommand(ctx context.Context, parsed *parser.ParsedCommand) error {
	cmd, exists := s.registry.Get(parsed.Command)
	if !exists {
		return fmt.Errorf("未知命令: %s。输入 'help' 查看可用命令", parsed.Command)
	}

	return cmd.Execute(ctx, parsed.Args)
}

// getPrompt 生成提示符
func (s *Shell) getPrompt() string {
	// 获取用户名
	username := "user"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	// 获取主机名
	hostname := "localhost"
	if h, err := os.Hostname(); err == nil {
		hostname = h
	}

	// 获取当前目录
	cwd := "~"
	if dir, err := os.Getwd(); err == nil {
		// 尝试用 ~ 替换家目录
		if home, err := os.UserHomeDir(); err == nil {
			if strings.HasPrefix(dir, home) {
				cwd = "~" + strings.TrimPrefix(dir, home)
			} else {
				cwd = dir
			}
		} else {
			cwd = dir
		}
	}

	// 格式: [username@hostname cwd]$
	const (
		colorGreen = "\033[32m"
		colorBlue  = "\033[34m"
		colorReset = "\033[0m"
	)

	prompt := fmt.Sprintf("%s[%s@%s %s%s%s]%s$ ",
		colorGreen,
		username,
		hostname,
		colorBlue,
		filepath.Base(cwd),
		colorGreen,
		colorReset,
	)

	return prompt
}

// printWelcome 打印欢迎信息
func (s *Shell) printWelcome() {
	const banner = `
╦  ╦╔═╗╦ ╦
║  ║╚═╗╠═╣
╩═╝╩╚═╝╩ ╩  Linux-style Shell

欢迎使用 Lish！轻量级 Linux 风格终端。
输入 'help' 查看可用命令，输入 'exit' 退出。
`
	fmt.Fprintln(s.stdout, banner)
}
