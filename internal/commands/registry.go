package commands

import (
	"fmt"
	"sort"
	"sync"
)

// Registry 命令注册表
type Registry struct {
	mu       sync.RWMutex
	commands map[string]Command
}

// NewRegistry 创建新的命令注册表
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

// Register 注册命令
func (r *Registry) Register(cmd Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := cmd.Name()
	if _, exists := r.commands[name]; exists {
		return fmt.Errorf("command %s already registered", name)
	}

	r.commands[name] = cmd
	return nil
}

// Get 获取命令
func (r *Registry) Get(name string) (Command, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, exists := r.commands[name]
	return cmd, exists
}

// List 列出所有命令
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// GetAll 获取所有命令
func (r *Registry) GetAll() map[string]Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]Command, len(r.commands))
	for name, cmd := range r.commands {
		result[name] = cmd
	}

	return result
}
