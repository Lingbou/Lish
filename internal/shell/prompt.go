package shell

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Lingbou/Lish/internal/theme"
)

// PromptFormatter 格式化提示符
type PromptFormatter struct {
	format string
	scheme theme.ColorScheme
}

// NewPromptFormatter 创建新的提示符格式化器
func NewPromptFormatter(format string, scheme theme.ColorScheme) *PromptFormatter {
	if format == "" {
		format = "[{user}@{host} {cwd}]$ "
	}
	return &PromptFormatter{
		format: format,
		scheme: scheme,
	}
}

// Format 生成提示符字符串
func (p *PromptFormatter) Format() string {
	prompt := p.format

	// 替换变量
	prompt = strings.ReplaceAll(prompt, "{user}", p.getUser())
	prompt = strings.ReplaceAll(prompt, "{host}", p.getHost())
	prompt = strings.ReplaceAll(prompt, "{cwd}", p.getCwd())
	prompt = strings.ReplaceAll(prompt, "{git}", p.getGitBranch())

	return prompt
}

// getUser 获取用户名
func (p *PromptFormatter) getUser() string {
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	if user == "" {
		user = "user"
	}
	return p.scheme.PromptUser().Apply(user)
}

// getHost 获取主机名
func (p *PromptFormatter) getHost() string {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "localhost"
	}
	return p.scheme.PromptHost().Apply(hostname)
}

// getCwd 获取当前工作目录
func (p *PromptFormatter) getCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "~"
	}

	// 将 HOME 替换为 ~
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE") // Windows
	}
	if home != "" && strings.HasPrefix(cwd, home) {
		cwd = "~" + strings.TrimPrefix(cwd, home)
	}

	// 只显示最后两级目录
	parts := strings.Split(filepath.ToSlash(cwd), "/")
	if len(parts) > 3 {
		cwd = ".../" + strings.Join(parts[len(parts)-2:], "/")
	}

	return p.scheme.PromptPath().Apply(cwd)
}

// getGitBranch 获取 Git 分支
func (p *PromptFormatter) getGitBranch() string {
	// 尝试读取 .git/HEAD
	gitHead, err := os.ReadFile(".git/HEAD")
	if err != nil {
		return ""
	}

	branch := strings.TrimSpace(string(gitHead))
	if strings.HasPrefix(branch, "ref: refs/heads/") {
		branch = strings.TrimPrefix(branch, "ref: refs/heads/")
		return " " + p.scheme.PromptGit().Apply("("+branch+")")
	}

	return ""
}

// UpdateScheme 更新配色方案
func (p *PromptFormatter) UpdateScheme(scheme theme.ColorScheme) {
	p.scheme = scheme
}

// SetFormat 设置提示符格式
func (p *PromptFormatter) SetFormat(format string) {
	p.format = format
}
