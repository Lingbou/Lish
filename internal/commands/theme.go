package commands

import (
	"context"
	"fmt"
	"sort"

	"github.com/Lingbou/Lish/internal/theme"
)

// ThemeCommand 主题管理命令
type ThemeCommand struct {
	manager *theme.Manager
}

// NewThemeCommand 创建新的主题命令
func NewThemeCommand(manager *theme.Manager) *ThemeCommand {
	return &ThemeCommand{manager: manager}
}

func (c *ThemeCommand) Name() string {
	return "theme"
}

func (c *ThemeCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		// 显示当前主题
		fmt.Printf("当前主题: %s\n", c.manager.CurrentTheme())
		fmt.Println("\n使用 'theme list' 查看所有可用主题")
		fmt.Println("使用 'theme set <name>' 切换主题")
		return nil
	}

	subcommand := args[0]

	switch subcommand {
	case "list":
		return c.listThemes()
	case "show":
		if len(args) < 2 {
			return fmt.Errorf("用法: theme show <name>")
		}
		return c.showTheme(args[1])
	case "set":
		if len(args) < 2 {
			return fmt.Errorf("用法: theme set <name>")
		}
		return c.setTheme(args[1])
	case "export":
		if len(args) < 2 {
			return fmt.Errorf("用法: theme export <name>")
		}
		return c.exportTheme(args[1])
	case "import":
		if len(args) < 2 {
			return fmt.Errorf("用法: theme import <file>")
		}
		return c.importTheme(args[1])
	default:
		return fmt.Errorf("未知子命令: %s\n使用 'theme' 或 'theme list' 查看帮助", subcommand)
	}
}

func (c *ThemeCommand) listThemes() error {
	themes := c.manager.ListThemes()
	current := c.manager.CurrentTheme()

	// 分离内置主题和自定义主题
	builtin := []string{}
	custom := []string{}

	for _, name := range themes {
		if theme.BuiltinThemes[name] != nil {
			builtin = append(builtin, name)
		} else {
			custom = append(custom, name)
		}
	}

	sort.Strings(builtin)
	sort.Strings(custom)

	fmt.Println("\n内置主题:")
	for _, name := range builtin {
		if name == current {
			fmt.Printf("  * %s (当前)\n", name)
		} else {
			fmt.Printf("    %s\n", name)
		}
	}

	if len(custom) > 0 {
		fmt.Println("\n自定义主题:")
		for _, name := range custom {
			fmt.Printf("    %s\n", name)
		}
	}

	fmt.Println("\n提示: 使用 'theme show <name>' 预览主题")
	fmt.Println("      使用 'theme set <name>' 切换主题")

	return nil
}

func (c *ThemeCommand) showTheme(name string) error {
	// 移除可能的 (custom) 后缀
	name = removeCustomSuffix(name)

	// 保存当前主题
	currentTheme := c.manager.CurrentTheme()

	// 加载要预览的主题
	if err := c.manager.LoadTheme(name); err != nil {
		return err
	}

	scheme := c.manager.CurrentScheme()

	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("主题预览: %s\n", name)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	fmt.Println("基础颜色:")
	fmt.Println("  " + scheme.Primary().Apply("■ Primary (主要)"))
	fmt.Println("  " + scheme.Success().Apply("■ Success (成功)"))
	fmt.Println("  " + scheme.Warning().Apply("■ Warning (警告)"))
	fmt.Println("  " + scheme.Error().Apply("■ Error (错误)"))
	fmt.Println("  " + scheme.Info().Apply("■ Info (信息)"))

	fmt.Println("\n文件类型:")
	fmt.Println("  " + scheme.Directory().Apply("■ Directory/"))
	fmt.Println("  " + scheme.Executable().Apply("■ Executable*"))
	fmt.Println("  " + scheme.Symlink().Apply("■ Symlink@"))
	fmt.Println("  " + scheme.Archive().Apply("■ Archive.zip"))

	fmt.Println("\n提示符颜色:")
	fmt.Print("  ")
	fmt.Print(scheme.PromptUser().Apply("user"))
	fmt.Print("@")
	fmt.Print(scheme.PromptHost().Apply("hostname"))
	fmt.Print(" ")
	fmt.Print(scheme.PromptPath().Apply("~/path"))
	fmt.Print(scheme.PromptGit().Apply(" (main)"))
	fmt.Print("$ \n")

	fmt.Println("\n语法高亮:")
	fmt.Print("  ")
	fmt.Print(scheme.SyntaxCommand().Apply("command"))
	fmt.Print(" ")
	fmt.Print(scheme.SyntaxArgument().Apply("-flag"))
	fmt.Print(" ")
	fmt.Print(scheme.SyntaxString().Apply("\"string\""))
	fmt.Print(" ")
	fmt.Print(scheme.SyntaxVariable().Apply("$var"))
	fmt.Print(" ")
	fmt.Print(scheme.SyntaxOperator().Apply("|"))
	fmt.Print("\n")

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 恢复当前主题
	c.manager.LoadTheme(currentTheme)

	return nil
}

func (c *ThemeCommand) setTheme(name string) error {
	// 移除可能的 (custom) 后缀
	name = removeCustomSuffix(name)

	if err := c.manager.LoadTheme(name); err != nil {
		return err
	}

	fmt.Printf("✓ 主题已切换为: %s\n", name)
	fmt.Println("\n提示: 重启 shell 以应用新的提示符颜色")
	fmt.Println("      或使用 'theme show %s' 查看效果", name)

	return nil
}

func (c *ThemeCommand) exportTheme(name string) error {
	return fmt.Errorf("导出功能开发中...")
}

func (c *ThemeCommand) importTheme(file string) error {
	return fmt.Errorf("导入功能开发中...")
}

func (c *ThemeCommand) Help() string {
	return `theme - 主题管理

用法:
  theme                  显示当前主题
  theme list             列出所有可用主题
  theme show <name>      显示主题预览
  theme set <name>       切换主题
  theme export <name>    导出主题配置
  theme import <file>    导入主题配置

内置主题:
  dark              - 默认暗色主题
  light             - 亮色主题
  solarized-dark    - Solarized Dark
  solarized-light   - Solarized Light
  gruvbox           - Gruvbox Dark
  dracula           - Dracula
  nord              - Nord
  monokai           - Monokai Pro

示例:
  theme list                    # 列出所有主题
  theme show dracula            # 预览 Dracula 主题
  theme set dracula             # 切换到 Dracula 主题
  theme                         # 显示当前主题`
}

func (c *ThemeCommand) ShortHelp() string {
	return "主题管理"
}

// removeCustomSuffix 移除主题名称中的 (custom) 后缀
func removeCustomSuffix(name string) string {
	if len(name) > 9 && name[len(name)-9:] == " (custom)" {
		return name[:len(name)-9]
	}
	return name
}

