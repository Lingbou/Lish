package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/Lingbou/Lish/internal/commands"
	"github.com/Lingbou/Lish/internal/completer"
	"github.com/Lingbou/Lish/internal/config"
	"github.com/Lingbou/Lish/internal/history"
	"github.com/Lingbou/Lish/internal/parser"
	"github.com/Lingbou/Lish/internal/script"
	"github.com/Lingbou/Lish/internal/theme"
	"github.com/chzyer/readline"
)

// Shell Lish Shell 结构
type Shell struct {
	registry        *commands.Registry
	history         *history.Manager
	config          *config.Config
	suggester       *Suggester
	themeManager    *theme.Manager
	promptFormatter *PromptFormatter
	scriptExecutor  *script.Executor
	rl              *readline.Instance
	stdout          *os.File
	stderr          *os.File
}

// NewShell 创建新的 Shell 实例
func NewShell() (*Shell, error) {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 创建命令注册表
	registry := commands.NewRegistry()

	// 创建历史管理器
	histMgr, err := history.NewManager()
	if err != nil {
		return nil, fmt.Errorf("创建历史管理器失败: %w", err)
	}

	// 创建主题管理器
	themeManager := theme.NewManager(cfg.Theme.CustomThemesDir)
	if err := themeManager.LoadTheme(cfg.Theme.Current); err != nil {
		// 如果加载失败，使用默认主题
		themeManager.LoadTheme("dark")
	}

	// 创建提示符格式化器
	promptFormatter := NewPromptFormatter(cfg.Prompt.Format, themeManager.CurrentScheme())

	shell := &Shell{
		registry:        registry,
		history:         histMgr,
		config:          cfg,
		suggester:       NewSuggester(),
		themeManager:    themeManager,
		promptFormatter: promptFormatter,
		stdout:          os.Stdout,
		stderr:          os.Stderr,
	}

	// 创建脚本执行器（使用 shell 作为命令执行器）
	shell.scriptExecutor = script.NewExecutor(shell)

	return shell, nil
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
		// 文件浏览
		commands.NewPwdCommand(s.stdout),
		commands.NewCdCommand(),
		commands.NewLsCommand(s.stdout),
		commands.NewFindCommand(s.stdout),
		commands.NewTreeCommand(s.stdout), // v0.3.0 新增

		// 文件操作
		commands.NewCatCommand(s.stdout),
		commands.NewMkdirCommand(),
		commands.NewRmCommand(),
		commands.NewTouchCommand(),
		commands.NewCpCommand(s.stdout),
		commands.NewMvCommand(s.stdout),
		commands.NewDiffCommand(s.stdout), // v0.3.0 新增

		// 文本处理
		commands.NewGrepCommand(s.stdout, os.Stdin),
		commands.NewHeadCommand(s.stdout),
		commands.NewTailCommand(s.stdout),
		commands.NewWcCommand(s.stdout),

		// 系统命令
		commands.NewEchoCommand(s.stdout),
		commands.NewClearCommand(s.stdout),
		commands.NewEnvCommand(s.stdout),
		commands.NewWhichCommand(s.stdout, s.registry),
		commands.NewHistoryCommand(s.stdout),
		commands.NewPsCommand(s.stdout),   // v0.3.0 新增
		commands.NewKillCommand(s.stdout), // v0.3.0 新增
		commands.NewDuCommand(s.stdout),   // v0.3.0 新增
		commands.NewDateCommand(s.stdout), // v0.3.0 新增

		// 配置和别名
		commands.NewAliasCommand(s.stdout, s.config),
		commands.NewUnaliasCommand(s.stdout, s.config),

		// 网络命令
		commands.NewCurlCommand(s.stdout), // v0.4.0 新增
		commands.NewPingCommand(s.stdout), // v0.4.0 新增

		// 压缩命令
		commands.NewZipCommand(s.stdout),   // v0.4.0 新增
		commands.NewUnzipCommand(s.stdout), // v0.4.0 新增

		// 主题命令
		commands.NewThemeCommand(s.themeManager), // v0.5.1 新增

		// 脚本命令
		commands.NewSourceCommand(s.scriptExecutor), // v0.5.2 新增
		commands.NewExecCommand(s),                  // v0.5.2 新增

		// 高级文本命令
		&commands.SortCommand{}, // v0.5.3 新增
		&commands.UniqCommand{}, // v0.5.3 新增
		&commands.SedCommand{},  // v0.5.3 新增
		&commands.AwkCommand{},  // v0.5.3 新增

		// 系统工具命令
		&commands.ChmodCommand{}, // v0.5.4 新增
		&commands.ChownCommand{}, // v0.5.4 新增
		&commands.LnCommand{},    // v0.5.4 新增
		&commands.DfCommand{},    // v0.5.4 新增

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

		// 添加到历史（用于智能建议）
		s.suggester.AddToHistory(line)

		// 展开别名
		line = s.expandAlias(line)

		// 解析命令
		parsed, err := parser.Parse(line)
		if err != nil {
			fmt.Fprintf(s.stderr, "❌ 解析错误: %v\n", err)
			continue
		}

		if parsed == nil || parsed.Command == "" {
			continue
		}

		// 记录开始时间
		startTime := time.Now()

		// 执行命令
		execErr := s.executeCommand(ctx, parsed)

		// 计算执行时间
		duration := time.Since(startTime)

		// 显示错误（带拼写建议）
		if execErr != nil {
			fmt.Fprintf(s.stderr, "❌ 错误: %v\n", execErr)

			// 如果是未知命令，提供拼写建议
			if strings.Contains(execErr.Error(), "未知命令") {
				if suggestion := s.suggester.SpellCheck(parsed.Command, s.registry.List()); suggestion != "" {
					fmt.Fprintf(s.stderr, "💡 你是否想输入: %s\n", suggestion)
				}
			}
		}

		// 显示执行时间（如果超过 100ms）
		if duration > 100*time.Millisecond {
			fmt.Fprintf(s.stderr, "⏱️  执行时间: %s\n", formatDuration(duration))
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

// ExecuteCommand 实现 script.CommandExecutor 接口
func (s *Shell) ExecuteCommand(ctx context.Context, command string, args []string) error {
	cmd, exists := s.registry.Get(command)
	if !exists {
		return fmt.Errorf("未知命令: %s", command)
	}
	return cmd.Execute(ctx, args)
}

// getPrompt 生成提示符
func (s *Shell) getPrompt() string {
	// 使用提示符格式化器
	if s.promptFormatter != nil {
		return s.promptFormatter.Format()
	}

	// 备用方案：简单提示符
	username := "user"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

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
输入 'help' 查看可用命令，输入 'exit' 退出。`
	fmt.Fprintln(s.stdout, banner)
}
