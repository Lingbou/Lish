package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Manager 管理主题的加载、切换、保存
type Manager struct {
	currentTheme  string
	currentScheme ColorScheme
	customThemes  map[string]ColorScheme
	themesDir     string
}

// NewManager 创建主题管理器
func NewManager(themesDir string) *Manager {
	return &Manager{
		currentTheme:  "dark",
		currentScheme: BuiltinThemes["dark"],
		customThemes:  make(map[string]ColorScheme),
		themesDir:     themesDir,
	}
}

// LoadTheme 加载主题
func (m *Manager) LoadTheme(name string) error {
	// 1. 先查找内置主题
	if scheme, ok := BuiltinThemes[name]; ok {
		m.currentTheme = name
		m.currentScheme = scheme
		return nil
	}

	// 2. 查找自定义主题
	if scheme, ok := m.customThemes[name]; ok {
		m.currentTheme = name
		m.currentScheme = scheme
		return nil
	}

	// 3. 从文件加载
	themePath := filepath.Join(m.themesDir, name+".json")
	scheme, err := m.loadThemeFromFile(themePath)
	if err != nil {
		return fmt.Errorf("主题 '%s' 不存在", name)
	}

	m.customThemes[name] = scheme
	m.currentTheme = name
	m.currentScheme = scheme
	return nil
}

// ListThemes 列出所有可用主题
func (m *Manager) ListThemes() []string {
	themes := make([]string, 0)

	// 内置主题
	for name := range BuiltinThemes {
		themes = append(themes, name)
	}

	// 自定义主题
	if entries, err := os.ReadDir(m.themesDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
				name := strings.TrimSuffix(entry.Name(), ".json")
				themes = append(themes, name+" (custom)")
			}
		}
	}

	return themes
}

// CurrentTheme 返回当前主题名称
func (m *Manager) CurrentTheme() string {
	return m.currentTheme
}

// CurrentScheme 返回当前配色方案
func (m *Manager) CurrentScheme() ColorScheme {
	return m.currentScheme
}

// SaveCustomTheme 保存自定义主题
func (m *Manager) SaveCustomTheme(name string, scheme ColorScheme) error {
	// 创建主题目录
	if err := os.MkdirAll(m.themesDir, 0755); err != nil {
		return fmt.Errorf("创建主题目录失败: %w", err)
	}

	// 序列化主题
	data, err := json.MarshalIndent(scheme, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化主题失败: %w", err)
	}

	// 写入文件
	themePath := filepath.Join(m.themesDir, name+".json")
	if err := os.WriteFile(themePath, data, 0644); err != nil {
		return fmt.Errorf("保存主题文件失败: %w", err)
	}

	return nil
}

// loadThemeFromFile 从文件加载主题
func (m *Manager) loadThemeFromFile(path string) (ColorScheme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 简单的自定义主题结构
	// 这里使用 DarkTheme 作为基础，实际项目中可以更复杂
	var customScheme DarkTheme
	if err := json.Unmarshal(data, &customScheme); err != nil {
		return nil, fmt.Errorf("解析主题文件失败: %w", err)
	}

	return &customScheme, nil
}

// GetHighlighter 获取当前主题的高亮器
func (m *Manager) GetHighlighter() *Highlighter {
	return NewHighlighter(m.currentScheme)
}

