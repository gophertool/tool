// example_usage.go
// 插件系统使用示例，展示如何在主程序中加载和使用插件
package plugin

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ExampleUsage 演示如何使用插件系统的完整示例
func ExampleUsage() {
	// 创建插件管理器
	manager := NewPluginManager()

	// 设置插件目录（实际使用时需要替换为真实路径）
	pluginDir := "./plugins"

	fmt.Println("=== 插件系统使用示例 ===")

	// 1. 加载所有插件
	fmt.Println("\n1. 加载插件...")
	err := manager.LoadAllPlugins(pluginDir)
	if err != nil {
		log.Printf("加载插件失败: %v", err)
		fmt.Println("注意：这是一个示例，实际使用时需要有真实的插件文件")
		return
	}

	// 2. 列出所有已加载的插件
	fmt.Println("\n2. 已加载的插件列表:")
	plugins := manager.ListPlugins()
	for _, plugin := range plugins {
		fmt.Printf("插件: %s (版本: %s)\n", plugin.Info.Name, plugin.Info.Version)
		fmt.Printf("  描述: %s\n", plugin.Info.Description)
		fmt.Printf("  作者: %s\n", plugin.Info.Author)
		fmt.Printf("  路径: %s\n", plugin.Path)
		fmt.Printf("  提供工具数量: %d\n", len(plugin.Tools))
	}

	// 3. 列出所有可用的工具
	fmt.Println("\n3. 可用工具列表:")
	tools := manager.ListTools()
	for _, tool := range tools {
		fmt.Printf("工具: %s\n", tool.Name)
		fmt.Printf("  描述: %s\n", tool.Description)
		fmt.Printf("  参数数量: %d\n", len(tool.InputSchema.Properties))
	}

	// 4. 调用工具示例
	fmt.Println("\n4. 工具调用示例:")

	// 示例1: 调用计算器工具
	fmt.Println("\n示例1: 计算器工具 - 加法运算")
	params := map[string]interface{}{
		"operation": "add",
		"numbers":   []interface{}{10.0, 20.0, 30.0},
	}

	result, err := manager.CallTool("calculator", params)
	if err != nil {
		fmt.Printf("调用失败: %v\n", err)
	} else {
		fmt.Printf("调用结果:\n")
		printCallToolResult(result)
	}

	// 示例2: 调用文本处理工具
	fmt.Println("\n示例2: 文本处理工具 - 转换大写")
	params = map[string]interface{}{
		"action": "uppercase",
		"text":   "hello world",
	}

	result, err = manager.CallTool("text_processor", params)
	if err != nil {
		fmt.Printf("调用失败: %v\n", err)
	} else {
		fmt.Printf("调用结果:\n")
		printCallToolResult(result)
	}

	// 5. 带上下文的工具调用示例
	fmt.Println("\n5. 带超时的工具调用示例:")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params = map[string]interface{}{
		"format":   "2006-01-02 15:04:05",
		"timezone": "Asia/Shanghai",
	}

	result, err = manager.CallToolWithContext(ctx, "time_tool", params)
	if err != nil {
		fmt.Printf("调用失败: %v\n", err)
	} else {
		fmt.Printf("调用结果:\n")
		printCallToolResult(result)
	}

	// 6. 错误处理示例
	fmt.Println("\n6. 错误处理示例:")
	fmt.Println("调用不存在的工具:")
	result, err = manager.CallTool("non_existent_tool", map[string]interface{}{})
	if err != nil {
		fmt.Printf("调用失败: %v\n", err)
	} else {
		fmt.Printf("调用结果:\n")
		printCallToolResult(result)
	}

	// 7. 关闭所有插件
	fmt.Println("\n7. 关闭插件系统...")
	manager.Shutdown()
	fmt.Println("插件系统已关闭")
}

// printCallToolResult 打印工具调用结果的辅助函数
func printCallToolResult(result *CallToolResult) {
	if result.IsError {
		fmt.Println("  状态: 错误")
	} else {
		fmt.Println("  状态: 成功")
	}

	fmt.Printf("  内容数量: %d\n", len(result.Content))
	for i, content := range result.Content {
		fmt.Printf("  内容[%d]:\n", i)
		fmt.Printf("    类型: %s\n", content.GetType())

		switch c := content.(type) {
		case TextContent:
			fmt.Printf("    文本: %s\n", c.Text)
			if c.Name != "" {
				fmt.Printf("    名称: %s\n", c.Name)
			}
		case FileContent:
			fmt.Printf("    文件类型: %s\n", c.FileType)
			fmt.Printf("    MIME类型: %s\n", c.MimeType)
			fmt.Printf("    数据长度: %d 字节\n", len(c.Data))
			if c.Name != "" {
				fmt.Printf("    文件名: %s\n", c.Name)
			}
		case StructContent:
			fmt.Printf("    结构体数据: %+v\n", c.Data)
			if c.Name != "" {
				fmt.Printf("    名称: %s\n", c.Name)
			}
		}
	}

	if result.Meta != nil && len(result.Meta) > 0 {
		fmt.Printf("  元数据: %+v\n", result.Meta)
	}
}

// ExamplePluginDevelopment 展示如何开发一个简单的插件
func ExamplePluginDevelopment() {
	fmt.Println("=== 插件开发示例 ===")
	fmt.Println()
	fmt.Println("要开发一个插件，需要实现 ToolPluginInterface 接口：")
	fmt.Println()
	fmt.Println("1. 创建插件结构体并实现三个方法：")
	fmt.Println("   - GetPluginInfo() (PluginInfo, error)")
	fmt.Println("   - GetTools() ([]Tool, error)")
	fmt.Println("   - CallTool(toolName string, params map[string]interface{}) (*CallToolResult, error)")
	fmt.Println()
	fmt.Println("2. 在main函数中调用 ServePlugin(yourPlugin) 启动插件服务器")
	fmt.Println()
	fmt.Println("3. 编译为可执行文件，文件名以 .tool.plugin 结尾")
	fmt.Println()
	fmt.Println("4. 将插件文件放入插件目录，主程序会自动扫描和加载")
	fmt.Println()
	fmt.Println("完整的插件开发示例请参考 example_plugin.go 文件")
}

// ExampleBuildingPlugins 展示如何构建和部署插件
func ExampleBuildingPlugins() {
	fmt.Println("=== 插件构建和部署示例 ===")
	fmt.Println()
	fmt.Println("构建插件的步骤：")
	fmt.Println()
	fmt.Println("1. 创建插件目录结构：")
	fmt.Println("   plugins/")
	fmt.Println("   ├── calculator/")
	fmt.Println("   │   └── main.go")
	fmt.Println("   │   └── main.go")
	fmt.Println("   └── build.sh")
	fmt.Println()
	fmt.Println("2. 编写构建脚本 (build.sh)：")
	fmt.Println("   #!/bin/bash")
	fmt.Println("   for dir in */; do")
	fmt.Println("     cd \"$dir\"")
	fmt.Println("     go build -o \"../${dir%/}.tool.plugin\" .")
	fmt.Println("     cd ..")
	fmt.Println("   done")
	fmt.Println()
	fmt.Println("3. 运行构建脚本：")
	fmt.Println("   chmod +x build.sh")
	fmt.Println("   ./build.sh")
	fmt.Println()
	fmt.Println("4. 生成的插件文件：")
	fmt.Println("   calculator.tool.plugin")
	fmt.Println("   text_processor.tool.plugin")
	fmt.Println()
	fmt.Println("5. 在主程序中加载：")
	fmt.Println("   manager := NewPluginManager()")
	fmt.Println("   err := manager.LoadAllPlugins(\"./plugins\")")
}
