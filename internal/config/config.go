package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config Lish 配置结构
type Config struct {
	Prompt  PromptConfig      `toml:"prompt"`
	History HistoryConfig     `toml:"history"`
	Aliases map[string]string `toml:"aliases"`
	Colors  ColorsConfig      `toml:"colors"`
	Theme   ThemeConfig       `toml:"theme"`
}

// PromptConfig 提示符配置
type PromptConfig struct {
	Format     string `toml:"format"`
	ShowTime   bool   `toml:"show_time"`
	ShowStatus bool   `toml:"show_status"`
}

// HistoryConfig 历史记录配置
type HistoryConfig struct {
	MaxSize int    `toml:"max_size"`
	File    string `toml:"file"`
}

// ColorsConfig 颜色配置
type ColorsConfig struct {
	Enabled    bool   `toml:"enabled"`
	Directory  string `toml:"directory"`
	Executable string `toml:"executable"`
	Prompt     string `toml:"prompt"`
}

// ThemeConfig 主题配置
type ThemeConfig struct {
	Current         string `toml:"current"`
	CustomThemesDir string `toml:"custom_themes_dir"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	histFile := filepath.Join(homeDir, ".lish_history")
	themesDir := filepath.Join(homeDir, ".lish", "themes")

	return &Config{
		Prompt: PromptConfig{
			Format:     "[{user}@{host} {cwd}]$ ",
			ShowTime:   false,
			ShowStatus: false,
		},
		History: HistoryConfig{
			MaxSize: 1000,
			File:    histFile,
		},
		Aliases: make(map[string]string),
		Colors: ColorsConfig{
			Enabled:    true,
			Directory:  "blue",
			Executable: "green",
			Prompt:     "green",
		},
		Theme: ThemeConfig{
			Current:         "dark",
			CustomThemesDir: themesDir,
		},
	}
}

// Load 加载配置文件
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// 查找配置文件
	configPath, err := getConfigPath()
	if err != nil {
		return cfg, nil // 使用默认配置
	}

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	// 读取配置文件
	if _, err := toml.DecodeFile(configPath, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return cfg, nil
}

// Save 保存配置到文件
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// 创建配置目录
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 打开文件
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("创建配置文件失败: %w", err)
	}
	defer file.Close()

	// 编码并写入
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("写入配置失败: %w", err)
	}

	return nil
}

// getConfigPath 获取配置文件路径
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".lishrc"), nil
}

// GetAlias 获取别名
func (c *Config) GetAlias(name string) (string, bool) {
	cmd, exists := c.Aliases[name]
	return cmd, exists
}

// SetAlias 设置别名
func (c *Config) SetAlias(name, command string) {
	if c.Aliases == nil {
		c.Aliases = make(map[string]string)
	}
	c.Aliases[name] = command
}

// RemoveAlias 删除别名
func (c *Config) RemoveAlias(name string) {
	delete(c.Aliases, name)
}
