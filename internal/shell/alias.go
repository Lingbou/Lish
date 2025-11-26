package shell

import (
	"strings"
)

// expandAlias 展开命令别名
func (s *Shell) expandAlias(line string) string {
	// 解析命令（只取第一个词）
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return line
	}

	cmdName := parts[0]

	// 检查是否有别名
	if aliasCmd, exists := s.config.GetAlias(cmdName); exists {
		// 替换命令名
		if len(parts) > 1 {
			return aliasCmd + " " + strings.Join(parts[1:], " ")
		}
		return aliasCmd
	}

	return line
}
