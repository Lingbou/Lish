package history

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultHistoryFile = ".lish_history"
	maxHistorySize     = 1000
)

// Manager 历史记录管理器
type Manager struct {
	historyFile string
	maxSize     int
}

// NewManager 创建历史记录管理器
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	historyFile := filepath.Join(homeDir, defaultHistoryFile)
	
	return &Manager{
		historyFile: historyFile,
		maxSize:     maxHistorySize,
	}, nil
}

// GetHistoryFile 返回历史记录文件路径
func (m *Manager) GetHistoryFile() string {
	return m.historyFile
}

// Search 搜索历史记录中匹配的命令
func (m *Manager) Search(prefix string) (string, error) {
	content, err := os.ReadFile(m.historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	
	lines := strings.Split(string(content), "\n")
	
	// 从后往前搜索，找到最近的匹配
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" && strings.HasPrefix(line, prefix) {
			return line, nil
		}
	}
	
	return "", nil
}

