package script

import (
	"fmt"
	"strings"
)

// Scope 表示变量作用域
type Scope struct {
	vars   map[string]string
	parent *Scope
}

// NewScope 创建新的作用域
func NewScope(parent *Scope) *Scope {
	return &Scope{
		vars:   make(map[string]string),
		parent: parent,
	}
}

// Set 设置变量
func (s *Scope) Set(name, value string) {
	s.vars[name] = value
}

// Get 获取变量值
func (s *Scope) Get(name string) (string, bool) {
	// 先在当前作用域查找
	if val, ok := s.vars[name]; ok {
		return val, true
	}
	// 如果没找到，在父作用域查找
	if s.parent != nil {
		return s.parent.Get(name)
	}
	return "", false
}

// Delete 删除变量
func (s *Scope) Delete(name string) {
	delete(s.vars, name)
}

// All 获取所有变量（包括父作用域）
func (s *Scope) All() map[string]string {
	result := make(map[string]string)

	// 先复制父作用域的变量
	if s.parent != nil {
		for k, v := range s.parent.All() {
			result[k] = v
		}
	}

	// 然后覆盖当前作用域的变量
	for k, v := range s.vars {
		result[k] = v
	}

	return result
}

// VariableManager 管理变量和作用域
type VariableManager struct {
	currentScope *Scope
}

// NewVariableManager 创建新的变量管理器
func NewVariableManager() *VariableManager {
	return &VariableManager{
		currentScope: NewScope(nil),
	}
}

// PushScope 进入新作用域
func (vm *VariableManager) PushScope() {
	vm.currentScope = NewScope(vm.currentScope)
}

// PopScope 退出当前作用域
func (vm *VariableManager) PopScope() {
	if vm.currentScope.parent != nil {
		vm.currentScope = vm.currentScope.parent
	}
}

// Set 设置变量
func (vm *VariableManager) Set(name, value string) {
	vm.currentScope.Set(name, value)
}

// Get 获取变量值
func (vm *VariableManager) Get(name string) (string, bool) {
	return vm.currentScope.Get(name)
}

// Delete 删除变量
func (vm *VariableManager) Delete(name string) {
	vm.currentScope.Delete(name)
}

// All 获取所有变量
func (vm *VariableManager) All() map[string]string {
	return vm.currentScope.All()
}

// Expand 展开变量引用
func (vm *VariableManager) Expand(str string) string {
	result := str

	// 处理 ${var} 形式
	for {
		start := strings.Index(result, "${")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		end += start

		varName := result[start+2 : end]
		value, _ := vm.Get(varName)
		result = result[:start] + value + result[end+1:]
	}

	// 处理 $var 形式
	words := strings.Fields(result)
	for i, word := range words {
		if strings.HasPrefix(word, "$") && !strings.Contains(word, "{") {
			varName := strings.TrimPrefix(word, "$")
			if value, ok := vm.Get(varName); ok {
				words[i] = value
			}
		}
	}

	return strings.Join(words, " ")
}

// SetSpecialVars 设置特殊变量
func (vm *VariableManager) SetSpecialVars(args []string, exitCode int) {
	// $0 - 脚本名称
	if len(args) > 0 {
		vm.Set("0", args[0])
	}

	// $1, $2, ... - 位置参数
	for i, arg := range args[1:] {
		vm.Set(fmt.Sprintf("%d", i+1), arg)
	}

	// $# - 参数个数
	if len(args) > 0 {
		vm.Set("#", fmt.Sprintf("%d", len(args)-1))
	} else {
		vm.Set("#", "0")
	}

	// $@ - 所有参数
	if len(args) > 1 {
		vm.Set("@", strings.Join(args[1:], " "))
	} else {
		vm.Set("@", "")
	}

	// $? - 上一个命令的退出码
	vm.Set("?", fmt.Sprintf("%d", exitCode))
}
