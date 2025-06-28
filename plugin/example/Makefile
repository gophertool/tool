# Makefile
# 用于编译和运行 hashicorp/go-plugin 示例的构建脚本
# 支持多层目录插件扫描和编译

# 定义变量
PLUGIN_DIR = plugin
MAIN_BINARY = main
GO_FILES = $(shell find . -name "*.go" -type f)

# 自动发现所有插件 - 查找所有有main函数的.go文件
PLUGIN_SOURCES = $(shell find $(PLUGIN_DIR) -name "*.go" -type f -exec grep -l "func main()" {} \;)
# 生成插件二进制文件路径（去掉.go后缀，添加.tool.plugin后缀）
PLUGIN_BINARIES = $(PLUGIN_SOURCES:.go=.tool.plugin)

# 默认目标
.PHONY: all
all: clean run

# 清理编译产物
.PHONY: clean
clean:
	@echo "正在清理编译产物..."
	@rm -f $(MAIN_BINARY)
	@find $(PLUGIN_DIR) -type f -name "*.tool.plugin" -delete 2>/dev/null || true
	@echo "清理完成"

# 创建插件目录
$(PLUGIN_DIR):
	@echo "创建插件目录..."
	@mkdir -p $(PLUGIN_DIR)

# 编译单个插件的通用规则
$(PLUGIN_DIR)/%.tool.plugin: $(PLUGIN_DIR)/%.go
	@echo "正在编译插件: $@"
	@mkdir -p $(dir $@)
	@go build -o $@ $<
	@echo "插件编译完成: $@"

# 编译所有插件
.PHONY: build-plugins
build-plugins: $(PLUGIN_DIR) $(PLUGIN_BINARIES)
	@echo "所有插件编译完成!"


# 编译主程序
.PHONY: build-main
build-main:
	@echo "正在编译主程序..."
	@go build -o $(MAIN_BINARY)
	@echo "主程序编译完成: $(MAIN_BINARY)"

# 编译所有组件
.PHONY: build
build: build-plugins build-main
	@echo "编译完成!"

# 运行示例
.PHONY: run
run: build
	@echo "正在运行示例程序..."
	@./$(MAIN_BINARY)
	@echo "示例程序运行完成"
	@$(MAKE) clean

# 运行测试
.PHONY: test
test: build
	@echo "正在运行测试..."
	@go test -v ./...
	@echo "测试完成"

# 下载依赖
.PHONY: deps
deps:
	@echo "正在下载依赖..."
	@go mod tidy
	@echo "依赖下载完成"

# 检查代码
.PHONY: check
check:
	@echo "正在检查代码..."
	@go vet ./...
	@go fmt ./...
	@echo "代码检查完成"

# 列出所有发现的插件
.PHONY: list-plugins
list-plugins:
	@echo "发现的插件源文件:"
	@for plugin in $(PLUGIN_SOURCES); do \
		echo "  - $$plugin"; \
	done
	@echo ""
	@echo "将生成的插件二进制文件:"
	@for binary in $(PLUGIN_BINARIES); do \
		echo "  - $$binary"; \
	done

# 清理特定插件
.PHONY: clean-plugin
clean-plugin:
	@if [ -z "$(PLUGIN)" ]; then \
		echo "使用方法: make clean-plugin PLUGIN=<插件名>"; \
		echo "例如: make clean-plugin PLUGIN=greeter"; \
	else \
		echo "正在清理插件: $(PLUGIN)"; \
		find $(PLUGIN_DIR) -name "$(PLUGIN).tool.plugin" -type f -delete 2>/dev/null || true; \
		echo "插件清理完成: $(PLUGIN)"; \
	fi

# 编译特定插件
.PHONY: build-specific
build-specific:
	@if [ -z "$(PLUGIN)" ]; then \
		echo "使用方法: make build-specific PLUGIN=<插件名>"; \
		echo "例如: make build-specific PLUGIN=greeter"; \
		echo "或者: make build-specific PLUGIN=advanced/calculator"; \
	else \
		echo "正在编译指定插件: $(PLUGIN_DIR)/$(PLUGIN).tool.plugin"; \
		if [ -f "$(PLUGIN_DIR)/$(PLUGIN).go" ]; then \
			mkdir -p $(dir $(PLUGIN_DIR)/$(PLUGIN)); \
			go build -o $(PLUGIN_DIR)/$(PLUGIN).tool.plugin $(PLUGIN_DIR)/$(PLUGIN).go; \
			echo "插件编译完成: $(PLUGIN_DIR)/$(PLUGIN).tool.plugin"; \
		else \
			echo "错误: 插件源文件不存在: $(PLUGIN_DIR)/$(PLUGIN).go"; \
		fi \
	fi

# 测试特定插件
.PHONY: test-plugin
test-plugin:
	@if [ -z "$(PLUGIN)" ]; then \
		echo "使用方法: make test-plugin PLUGIN=<插件名>"; \
		echo "例如: make test-plugin PLUGIN=greeter"; \
	else \
		if [ -f "$(PLUGIN_DIR)/$(PLUGIN).tool.plugin" ]; then \
			echo "正在测试插件: $(PLUGIN)"; \
			echo "插件路径: $(PLUGIN_DIR)/$(PLUGIN).tool.plugin"; \
			if [ -f "$(PLUGIN_DIR)/$(PLUGIN).tool.plugin" ]; then \
				echo "插件文件存在且可执行"; \
			else \
				echo "插件文件不存在或不可执行"; \
			fi; \
		else \
			echo "插件二进制文件不存在: $(PLUGIN_DIR)/$(PLUGIN).tool.plugin"; \
			echo "请先运行: make build-specific PLUGIN=$(PLUGIN)"; \
		fi \
	fi

# 显示帮助信息
.PHONY: help
help:
	@echo "可用的命令:"
	@echo "  make deps           - 下载Go模块依赖"
	@echo "  make build          - 编译所有插件和主程序"
	@echo "  make build-plugins  - 编译所有发现的插件"
	@echo "  make build-main     - 只编译主程序"
	@echo "  make run            - 编译并运行示例"
	@echo "  make check          - 检查和格式化代码"
	@echo "  make clean          - 清理所有编译产物"
	@echo "  make list-plugins   - 列出所有发现的插件"
	@echo ""
	@echo "高级命令:"
	@echo "  make build-specific PLUGIN=<插件名> - 编译指定插件"
	@echo "  make clean-plugin PLUGIN=<插件名>   - 清理指定插件"
	@echo "  make test-plugin PLUGIN=<插件名>    - 测试指定插件"
	@echo ""
	@echo "示例:"
	@echo "  make build-specific PLUGIN=timetool/tool"
	@echo "  make clean-plugin PLUGIN=timetool/tool"
	@echo "  make test-plugin PLUGIN=timetool/tool"
	@echo ""
	@echo "  make help           - 显示此帮助信息"