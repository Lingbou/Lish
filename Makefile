# Makefile for Lish

# 变量定义
BINARY_NAME=lish
MAIN_PATH=cmd/lish/main.go
BUILD_DIR=build
INSTALL_PATH=/usr/local/bin

# 编译标志
LDFLAGS=-ldflags="-s -w"

# 默认目标
.PHONY: all
all: build

# 编译当前平台
.PHONY: build
build:
	@echo "正在编译 Lish（当前平台）..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ 编译完成: $(BINARY_NAME)"
	@ls -lh $(BINARY_NAME)

# 编译 Windows 版本（在 Linux/WSL 下使用）
.PHONY: build-windows
build-windows:
	@echo "正在编译 Windows 版本..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o lish.exe $(MAIN_PATH)
	@echo "✓ 编译完成: lish.exe"
	@ls -lh lish.exe

# 构建所有平台的版本
.PHONY: build-all
build-all:
	@echo "正在为所有平台编译..."
	@mkdir -p $(BUILD_DIR)
	@echo "  - Linux (amd64)"
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "  - Linux (arm64)"
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "  - Windows (amd64)"
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "  - macOS (amd64)"
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo "  - macOS (arm64)"
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "✓ 所有平台编译完成"
	@ls -lh $(BUILD_DIR)/

# 运行
.PHONY: run
run: build
	@echo "正在运行 Lish..."
	@./$(BINARY_NAME)

# 安装到系统
.PHONY: install
install: build
	@echo "正在安装 Lish 到 $(INSTALL_PATH)..."
	@sudo cp $(BINARY_NAME) $(INSTALL_PATH)/
	@sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ 安装完成"
	@which $(BINARY_NAME)

# 卸载
.PHONY: uninstall
uninstall:
	@echo "正在卸载 Lish..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ 卸载完成"

# 清理
.PHONY: clean
clean:
	@echo "正在清理..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "✓ 清理完成"

# 测试
.PHONY: test
test:
	@echo "正在运行测试..."
	@go test -v ./...

# 格式化代码
.PHONY: fmt
fmt:
	@echo "正在格式化代码..."
	@go fmt ./...
	@echo "✓ 格式化完成"

# 代码检查
.PHONY: lint
lint:
	@echo "正在检查代码..."
	@go vet ./...
	@echo "✓ 检查完成"

# 整理依赖
.PHONY: tidy
tidy:
	@echo "正在整理依赖..."
	@go mod tidy
	@echo "✓ 依赖整理完成"

# 开发模式（格式化、检查、编译）
.PHONY: dev
dev: fmt lint build
	@echo "✓ 开发构建完成"

# 发布构建（包含所有平台）
.PHONY: release
release: clean tidy fmt lint test build-all
	@echo "✓ 发布构建完成"

# 帮助信息
.PHONY: help
help:
	@echo "Lish Makefile 命令："
	@echo ""
	@echo "  make build         - 编译当前平台的版本"
	@echo "  make build-windows - 编译 Windows 版本（交叉编译）"
	@echo "  make build-all     - 编译所有平台的版本"
	@echo "  make run           - 编译并运行"
	@echo "  make install       - 安装到系统 (需要 sudo)"
	@echo "  make uninstall     - 从系统卸载 (需要 sudo)"
	@echo "  make clean         - 清理编译文件"
	@echo "  make test          - 运行测试"
	@echo "  make fmt           - 格式化代码"
	@echo "  make lint          - 代码检查"
	@echo "  make tidy          - 整理依赖"
	@echo "  make dev           - 开发模式构建"
	@echo "  make release       - 发布构建（所有平台）"
	@echo "  make help          - 显示此帮助信息"
	@echo ""
	@echo "提示："
	@echo "  在 WSL2 中为 Windows 编译请使用: make build-windows"
	@echo ""

