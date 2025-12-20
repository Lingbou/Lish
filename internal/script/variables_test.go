package script

import (
	"testing"
)

func TestVariableManager(t *testing.T) {
	vm := NewVariableManager()

	// 测试设置和获取变量
	vm.Set("name", "test")
	value, exists := vm.Get("name")
	if !exists {
		t.Error("Variable should exist")
	}
	if value != "test" {
		t.Errorf("Value = %v, want test", value)
	}

	// 测试不存在的变量
	_, exists = vm.Get("nonexistent")
	if exists {
		t.Error("Nonexistent variable should not exist")
	}
}

func TestVariableScope(t *testing.T) {
	vm := NewVariableManager()

	// 全局作用域设置变量
	vm.Set("global", "value1")

	// 进入新作用域
	vm.PushScope()
	vm.Set("local", "value2")

	// 在新作用域中应该能访问全局变量
	value, exists := vm.Get("global")
	if !exists || value != "value1" {
		t.Error("Should access global variable from local scope")
	}

	// 在新作用域中应该能访问局部变量
	value, exists = vm.Get("local")
	if !exists || value != "value2" {
		t.Error("Should access local variable")
	}

	// 退出作用域
	vm.PopScope()

	// 局部变量应该不可访问
	_, exists = vm.Get("local")
	if exists {
		t.Error("Local variable should not be accessible after PopScope")
	}

	// 全局变量应该仍然可访问
	value, exists = vm.Get("global")
	if !exists || value != "value1" {
		t.Error("Global variable should still be accessible")
	}
}

func TestVariableExpand(t *testing.T) {
	vm := NewVariableManager()
	vm.Set("name", "world")
	vm.Set("greeting", "hello")

	tests := []struct {
		input    string
		expected string
	}{
		{"$name", "world"},
		{"${name}", "world"},
		{"hello $name", "hello world"},
		{"$greeting $name", "hello world"},
		{"no variables", "no variables"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := vm.Expand(tt.input)
			if result != tt.expected {
				t.Errorf("Expand(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSpecialVariables(t *testing.T) {
	vm := NewVariableManager()
	args := []string{"script.sh", "arg1", "arg2", "arg3"}
	exitCode := 0

	vm.SetSpecialVars(args, exitCode)

	// 测试 $0
	value, _ := vm.Get("0")
	if value != "script.sh" {
		t.Errorf("$0 = %v, want script.sh", value)
	}

	// 测试 $1, $2, $3
	value, _ = vm.Get("1")
	if value != "arg1" {
		t.Errorf("$1 = %v, want arg1", value)
	}

	// 测试 $#
	value, _ = vm.Get("#")
	if value != "3" {
		t.Errorf("$# = %v, want 3", value)
	}

	// 测试 $@
	value, _ = vm.Get("@")
	if value != "arg1 arg2 arg3" {
		t.Errorf("$@ = %v, want 'arg1 arg2 arg3'", value)
	}

	// 测试 $?
	value, _ = vm.Get("?")
	if value != "0" {
		t.Errorf("$? = %v, want 0", value)
	}
}

func BenchmarkVariableSet(b *testing.B) {
	vm := NewVariableManager()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.Set("test", "value")
	}
}

func BenchmarkVariableGet(b *testing.B) {
	vm := NewVariableManager()
	vm.Set("test", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.Get("test")
	}
}

func BenchmarkVariableExpand(b *testing.B) {
	vm := NewVariableManager()
	vm.Set("name", "world")
	input := "hello $name, how are you?"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.Expand(input)
	}
}

