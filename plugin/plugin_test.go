// plugin_test.go
// 插件系统测试文件，提供完整的测试套件
// 测试插件的加载、管理和调用功能
package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestPluginManager 测试插件管理器的基本功能
func TestPluginManager(t *testing.T) {
	// 创建测试插件目录
	testPluginDir := filepath.Join(os.TempDir(), "test_plugins")

	// 确保测试目录存在
	err := os.MkdirAll(testPluginDir, 0755)
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	defer os.RemoveAll(testPluginDir) // 测试结束后清理

	// 创建插件管理器
	manager := NewPluginManager()

	// 测试插件管理器初始化
	if manager == nil {
		t.Fatal("插件管理器创建失败")
	}

	// 测试空目录扫描
	pluginPaths, err := manager.ScanPlugins(testPluginDir)
	if err != nil {
		t.Errorf("扫描空插件目录失败: %v", err)
	}
	if len(pluginPaths) != 0 {
		t.Errorf("空目录应该返回空插件列表，实际返回: %v", pluginPaths)
	}

	// 测试加载所有插件（空目录）
	err = manager.LoadAllPlugins(testPluginDir)
	if err != nil {
		t.Errorf("加载空目录插件失败: %v", err)
	}

	// 测试列出插件（应为空）
	plugins := manager.ListPlugins()
	if len(plugins) != 0 {
		t.Errorf("应该没有加载的插件，实际有: %d", len(plugins))
	}

	// 测试列出工具（应为空）
	tools := manager.ListTools()
	if len(tools) != 0 {
		t.Errorf("应该没有可用的工具，实际有: %d", len(tools))
	}

	// 测试调用不存在的工具
	_, err = manager.CallTool("non_existent_tool", map[string]interface{}{})
	if err == nil {
		t.Error("调用不存在的工具应该返回错误")
	} else {
		t.Logf("调用不存在的工具返回预期错误: %v", err)
	}

	// 测试带上下文的工具调用（不存在的工具）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = manager.CallToolWithContext(ctx, "non_existent_tool", map[string]interface{}{})
	if err == nil {
		t.Error("带上下文调用不存在的工具应该返回错误")
	} else {
		t.Logf("带上下文调用不存在的工具返回预期错误: %v", err)
	}

	// 测试关闭插件系统
	manager.Shutdown()
}

// TestPluginScanAndLoad 测试插件扫描和加载功能
func TestPluginScanAndLoad(t *testing.T) {
	// 创建测试插件目录
	testPluginDir := filepath.Join(os.TempDir(), "test_plugins_scan")

	// 确保测试目录存在
	err := os.MkdirAll(testPluginDir, 0755)
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	defer os.RemoveAll(testPluginDir) // 测试结束后清理

	// 创建一个假的插件文件（不可执行）
	fakePluginPath := filepath.Join(testPluginDir, "fake.tool.plugin")
	fakeFile, err := os.Create(fakePluginPath)
	if err != nil {
		t.Fatalf("创建假插件文件失败: %v", err)
	}
	fakeFile.Close()

	// 创建插件管理器
	manager := NewPluginManager()

	// 测试扫描（应该找不到可执行的插件）
	pluginPaths, err := manager.ScanPlugins(testPluginDir)
	if err != nil {
		t.Errorf("扫描插件目录失败: %v", err)
	}
	if len(pluginPaths) != 0 {
		t.Errorf("应该没有找到可执行的插件，实际找到: %v", pluginPaths)
	}

	// 设置文件为可执行
	err = os.Chmod(fakePluginPath, 0755)
	if err != nil {
		t.Fatalf("设置文件可执行权限失败: %v", err)
	}

	// 再次扫描（应该找到可执行的插件）
	pluginPaths, err = manager.ScanPlugins(testPluginDir)
	if err != nil {
		t.Errorf("扫描插件目录失败: %v", err)
	}
	if len(pluginPaths) != 1 {
		t.Errorf("应该找到1个可执行的插件，实际找到: %d", len(pluginPaths))
	}

	// 测试加载无效插件（会失败，但不应该崩溃）
	_, err = manager.LoadPlugin(fakePluginPath)
	if err == nil {
		t.Error("加载无效插件应该返回错误")
	}
}

// TestCallToolResultCreation 测试工具调用结果的创建和使用
func TestCallToolResultCreation(t *testing.T) {
	// 测试创建成功结果
	successResult := NewCallToolResult()
	if successResult.IsError {
		t.Error("成功结果的IsError应该为false")
	}
	if len(successResult.Content) != 0 {
		t.Error("新创建的成功结果应该没有内容")
	}

	// 测试创建错误结果
	errorMsg := "测试错误消息"
	errorResult := NewErrorResult(errorMsg)
	if !errorResult.IsError {
		t.Error("错误结果的IsError应该为true")
	}
	if len(errorResult.Content) != 1 {
		t.Error("错误结果应该有一个内容项")
	}

	// 检查错误内容
	if len(errorResult.Content) > 0 {
		content, ok := errorResult.Content[0].(TextContent)
		if !ok {
			t.Error("错误结果的内容应该是TextContent类型")
		} else if content.Text != errorMsg {
			t.Errorf("错误消息不匹配，期望: %s, 实际: %s", errorMsg, content.Text)
		}
	}

	// 测试添加文本内容
	successResult.AddTextContent("测试文本")
	if len(successResult.Content) != 1 {
		t.Error("添加文本后应该有一个内容项")
	}

	// 测试添加结构体内容
	testStruct := map[string]interface{}{
		"key": "value",
	}
	successResult.AddStructContent(testStruct)
	if len(successResult.Content) != 2 {
		t.Error("添加结构体后应该有两个内容项")
	}

	// 测试添加文件内容
	testData := "测试文件内容"
	successResult.AddFileContent(FileTypeDocument, testData, "text/plain", "test.txt")
	if len(successResult.Content) != 3 {
		t.Error("添加文件后应该有三个内容项")
	}

	// 测试添加元数据
	successResult.SetMeta("meta_key", "meta_value")
	if successResult.Meta == nil || successResult.Meta["meta_key"] != "meta_value" {
		t.Error("添加元数据失败")
	}
}

// TestContentTypes 测试内容类型的功能
func TestContentTypes(t *testing.T) {
	// 测试文本内容
	textContent := TextContent{
		Text: "测试文本",
		Name: "测试名称",
	}
	if textContent.GetType() != "text" {
		t.Errorf("文本内容类型应该是text，实际是: %s", textContent.GetType())
	}

	// 测试文件内容
	fileContent := FileContent{
		Data:     "测试文件内容",
		MimeType: "text/plain",
		FileType: FileTypeDocument,
		Name:     "test.txt",
		Type:     ContentTypeFile,
	}
	if fileContent.GetType() != "file" {
		t.Errorf("文件内容类型应该是file，实际是: %s", fileContent.GetType())
	}

	// 测试结构体内容
	structContent := StructContent{
		Data: map[string]interface{}{
			"key": "value",
		},
		Name: "测试结构体",
	}
	if structContent.GetType() != "struct" {
		t.Errorf("结构体内容类型应该是struct，实际是: %s", structContent.GetType())
	}
}

// TestToolSchemaValidation 测试工具模式的功能
func TestToolSchemaValidation(t *testing.T) {
	// 创建一个测试工具模式
	schema := ToolInputSchema{
		Type: "object",
		Properties: map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "名称",
			},
			"age": map[string]any{
				"type":        "number",
				"description": "年龄",
			},
		},
		Required: []string{"name"},
	}

	// 添加一个简单的验证函数
	validateParams := func(schema ToolInputSchema, params map[string]interface{}) error {
		// 检查必需参数
		for _, required := range schema.Required {
			if _, exists := params[required]; !exists {
				return fmt.Errorf("缺少必需参数: %s", required)
			}
		}

		// 检查参数类型
		for name, value := range params {
			propDef, exists := schema.Properties[name]
			if !exists {
				continue // 忽略未定义的参数
			}

			propMap, ok := propDef.(map[string]any)
			if !ok {
				continue
			}

			expectedType, ok := propMap["type"].(string)
			if !ok {
				continue
			}

			switch expectedType {
			case "string":
				if _, ok := value.(string); !ok {
					return fmt.Errorf("参数 %s 应该是字符串类型", name)
				}
			case "number":
				if _, ok := value.(float64); !ok {
					if _, ok := value.(int); !ok {
						return fmt.Errorf("参数 %s 应该是数字类型", name)
					}
				}
			case "boolean":
				if _, ok := value.(bool); !ok {
					return fmt.Errorf("参数 %s 应该是布尔类型", name)
				}
			case "array":
				if _, ok := value.([]interface{}); !ok {
					return fmt.Errorf("参数 %s 应该是数组类型", name)
				}
			case "object":
				if _, ok := value.(map[string]interface{}); !ok {
					return fmt.Errorf("参数 %s 应该是对象类型", name)
				}
			}
		}

		return nil
	}

	// 测试有效参数验证
	validParams := map[string]interface{}{
		"name": "测试名称",
		"age":  30,
	}
	err := validateParams(schema, validParams)
	if err != nil {
		t.Errorf("有效参数验证失败: %v", err)
	}

	// 测试缺少必需参数
	invalidParams := map[string]interface{}{
		"age": 30,
	}
	err = validateParams(schema, invalidParams)
	if err == nil {
		t.Error("缺少必需参数应该返回错误")
	}

	// 测试参数类型错误
	invalidTypeParams := map[string]interface{}{
		"name": "测试名称",
		"age":  "三十", // 应该是数字
	}
	err = validateParams(schema, invalidTypeParams)
	if err == nil {
		t.Error("参数类型错误应该返回错误")
	}
}
