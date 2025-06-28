package main

import (
	"fmt"
	"log"

	"github.com/gophertool/tool/plugin"
)

func main() {
	// 创建插件管理器
	pluginManager := plugin.NewPluginManager()
	defer pluginManager.Shutdown()

	// 获取插件目录
	path := "./plugin"

	// 加载插件
	err := pluginManager.LoadAllPlugins(path)
	if err != nil {
		log.Fatalf("加载插件失败: %v", err)
	}

	// 打印已加载的插件
	fmt.Println("已加载的插件:")
	for _, p := range pluginManager.ListPlugins() {
		fmt.Printf("- %s (v%s): %s\n", p.Info.Name, p.Info.Version, p.Info.Description)
	}

	// 获取所有可用的工具
	tools := pluginManager.ListTools()
	fmt.Println("\n可用的工具:")
	for _, tool := range tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// 调用时间工具
	fmt.Println("\n调用当前时间工具:")
	result, err := pluginManager.CallTool("current_time", map[string]interface{}{
		"format":   "2006-01-02 15:04:05",
		"timezone": "Asia/Shanghai",
	})
	if err != nil {
		log.Fatalf("调用工具失败: %v", err)
	}

	// 打印结果
	for _, content := range result.Content {
		if textContent, ok := content.(plugin.TextContent); ok {
			fmt.Printf("当前时间: %s\n", textContent.Text)
		}
	}

	// 调用时间转换工具
	fmt.Println("\n调用时间转换工具:")
	result, err = pluginManager.CallTool("time_convert", map[string]interface{}{
		"time":          "2023-01-01 12:00:00",
		"source_format": "2006-01-02 15:04:05",
		"target_format": "Jan 02, 2006 03:04 PM",
	})
	if err != nil {
		log.Fatalf("调用工具失败: %v", err)
	}

	// 打印结果
	for _, content := range result.Content {
		if textContent, ok := content.(plugin.TextContent); ok {
			fmt.Printf("转换后的时间: %s\n", textContent.Text)
		}
	}

	// 调用时间计算工具
	fmt.Println("\n调用时间计算工具:")
	result, err = pluginManager.CallTool("time_calc", map[string]interface{}{
		"time":    "2023-01-01 12:00:00",
		"format":  "2006-01-02 15:04:05",
		"days":    7,
		"hours":   12,
		"minutes": 30,
	})
	if err != nil {
		log.Fatalf("调用工具失败: %v", err)
	}

	// 打印结果
	for _, content := range result.Content {
		if textContent, ok := content.(plugin.TextContent); ok {
			fmt.Printf("计算后的时间: %s\n", textContent.Text)
		}
	}
}
