// content_test.go
// 插件系统内容类型和工具模式测试文件
// 测试内容类型和工具模式的功能
package plugin

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestContentCreation 测试内容创建和类型判断
func TestContentCreation(t *testing.T) {
	// 测试创建文本内容
	textContent := TextContent{
		Text: "测试文本内容",
		Name: "测试文本",
	}

	// 验证类型
	if textContent.GetType() != "text" {
		t.Errorf("文本内容类型错误，期望: text, 实际: %s", textContent.GetType())
	}

	// 验证内容
	if textContent.Text != "测试文本内容" {
		t.Errorf("文本内容错误，期望: 测试文本内容, 实际: %s", textContent.Text)
	}

	// 验证名称
	if textContent.Name != "测试文本" {
		t.Errorf("文本名称错误，期望: 测试文本, 实际: %s", textContent.Name)
	}

	// 测试创建文件内容
	fileContent := FileContent{
		Data:     "测试文件内容",
		MimeType: "text/plain",
		FileType: FileTypeDocument,
		Name:     "test.txt",
		Type:     ContentTypeFile,
	}

	// 验证类型
	if fileContent.GetType() != "file" {
		t.Errorf("文件内容类型错误，期望: file, 实际: %s", fileContent.GetType())
	}

	// 验证内容
	if string(fileContent.Data) != "测试文件内容" {
		t.Errorf("文件内容错误，期望: 测试文件内容, 实际: %s", string(fileContent.Data))
	}

	// 验证MIME类型
	if fileContent.MimeType != "text/plain" {
		t.Errorf("文件MIME类型错误，期望: text/plain, 实际: %s", fileContent.MimeType)
	}

	// 验证文件类型
	if fileContent.FileType != FileTypeDocument {
		t.Errorf("文件类型错误，期望: %s, 实际: %s", FileTypeDocument, fileContent.FileType)
	}

	// 验证文件名
	if fileContent.Name != "test.txt" {
		t.Errorf("文件名错误，期望: test.txt, 实际: %s", fileContent.Name)
	}

	// 测试创建结构体内容
	testStruct := map[string]interface{}{
		"key":  "value",
		"num":  123,
		"bool": true,
	}
	structContent := StructContent{
		Data: testStruct,
		Name: "测试结构体",
	}

	// 验证类型
	if structContent.GetType() != "struct" {
		t.Errorf("结构体内容类型错误，期望: struct, 实际: %s", structContent.GetType())
	}

	// 验证内容
	if structContent.Data.(map[string]interface{})["key"] != "value" {
		t.Errorf("结构体内容错误，期望: value, 实际: %v", structContent.Data.(map[string]interface{})["key"])
	}

	// 验证名称
	if structContent.Name != "测试结构体" {
		t.Errorf("结构体名称错误，期望: 测试结构体, 实际: %s", structContent.Name)
	}
}

// TestCallToolResultFunctions 测试工具调用结果的功能
func TestCallToolResultFunctions(t *testing.T) {
	// 创建成功结果
	result := NewCallToolResult()

	// 验证初始状态
	if result.IsError {
		t.Error("成功结果的IsError应该为false")
	}

	if len(result.Content) != 0 {
		t.Errorf("初始内容应该为空，实际长度: %d", len(result.Content))
	}

	if result.Meta != nil {
		t.Error("初始元数据应该为nil")
	}

	// 添加文本内容
	result.AddTextContent("测试文本")

	// 验证内容添加
	if len(result.Content) != 1 {
		t.Errorf("添加文本后内容长度应该为1，实际: %d", len(result.Content))
	}

	// 验证内容类型
	content, ok := result.Content[0].(TextContent)
	if !ok {
		t.Error("添加的内容类型应该是TextContent")
	} else {
		if content.Text != "测试文本" {
			t.Errorf("文本内容错误，期望: 测试文本, 实际: %s", content.Text)
		}
	}

	// 添加带名称的文本内容
	result.AddTextContent("测试文本2", "文本名称")

	// 验证内容添加
	if len(result.Content) != 2 {
		t.Errorf("添加第二个文本后内容长度应该为2，实际: %d", len(result.Content))
	}

	// 验证内容类型和名称
	content, ok = result.Content[1].(TextContent)
	if !ok {
		t.Error("添加的内容类型应该是TextContent")
	} else {
		if content.Text != "测试文本2" {
			t.Errorf("文本内容错误，期望: 测试文本2, 实际: %s", content.Text)
		}
		if content.Name != "文本名称" {
			t.Errorf("文本名称错误，期望: 文本名称, 实际: %s", content.Name)
		}
	}

	// 添加结构体内容
	testStruct := map[string]interface{}{
		"key": "value",
	}
	result.AddStructContent(testStruct)

	// 验证内容添加
	if len(result.Content) != 3 {
		t.Errorf("添加结构体后内容长度应该为3，实际: %d", len(result.Content))
	}

	// 添加文件内容
	testData := "测试文件内容"
	result.AddFileContent(FileTypeDocument, testData, "text/plain", "test.txt")

	// 验证内容添加
	if len(result.Content) != 4 {
		t.Errorf("添加文件后内容长度应该为4，实际: %d", len(result.Content))
	}

	// 添加元数据
	result.SetMeta("meta_key", "meta_value")

	// 验证元数据添加
	if result.Meta == nil {
		t.Error("添加元数据后Meta不应该为nil")
	} else if result.Meta["meta_key"] != "meta_value" {
		t.Errorf("元数据值错误，期望: meta_value, 实际: %v", result.Meta["meta_key"])
	}

	// 测试错误结果
	errorResult := NewErrorResult("测试错误")

	// 验证错误状态
	if !errorResult.IsError {
		t.Error("错误结果的IsError应该为true")
	}

	// 验证错误内容
	if len(errorResult.Content) != 1 {
		t.Errorf("错误结果应该有一个内容项，实际: %d", len(errorResult.Content))
	} else {
		content, ok := errorResult.Content[0].(TextContent)
		if !ok {
			t.Error("错误内容类型应该是TextContent")
		} else if content.Text != "测试错误" {
			t.Errorf("错误消息不匹配，期望: 测试错误, 实际: %s", content.Text)
		}
	}
}

// TestToolInputSchema 测试工具输入模式的功能
func TestToolInputSchema(t *testing.T) {
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
			"tags": map[string]any{
				"type":        "array",
				"description": "标签",
			},
			"settings": map[string]any{
				"type":        "object",
				"description": "设置",
			},
			"enabled": map[string]any{
				"type":        "boolean",
				"description": "是否启用",
			},
		},
		Required: []string{"name", "enabled"},
	}

	// 测试JSON序列化
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		t.Errorf("工具模式序列化失败: %v", err)
	}

	// 测试JSON反序列化
	var deserializedSchema ToolInputSchema
	err = json.Unmarshal(schemaJSON, &deserializedSchema)
	if err != nil {
		t.Errorf("工具模式反序列化失败: %v", err)
	}

	// 验证反序列化结果
	if deserializedSchema.Type != schema.Type {
		t.Errorf("类型不匹配，期望: %s, 实际: %s", schema.Type, deserializedSchema.Type)
	}

	if len(deserializedSchema.Properties) != len(schema.Properties) {
		t.Errorf("属性数量不匹配，期望: %d, 实际: %d", len(schema.Properties), len(deserializedSchema.Properties))
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
		"name":     "测试名称",
		"age":      30,
		"tags":     []interface{}{"tag1", "tag2"},
		"settings": map[string]interface{}{"key": "value"},
		"enabled":  true,
	}

	err = validateParams(schema, validParams)
	if err != nil {
		t.Errorf("有效参数验证失败: %v", err)
	}

	// 测试缺少必需参数
	invalidParams1 := map[string]interface{}{
		"age":     30,
		"enabled": true,
	}

	err = validateParams(schema, invalidParams1)
	if err == nil {
		t.Error("缺少必需参数name应该返回错误")
	}

	// 测试参数类型错误
	invalidParams2 := map[string]interface{}{
		"name":    "测试名称",
		"age":     "三十", // 应该是数字
		"enabled": true,
	}

	err = validateParams(schema, invalidParams2)
	if err == nil {
		t.Error("参数类型错误应该返回错误")
	}

	// 测试布尔值类型错误
	invalidParams3 := map[string]interface{}{
		"name":    "测试名称",
		"enabled": "true", // 应该是布尔值
	}

	err = validateParams(schema, invalidParams3)
	if err == nil {
		t.Error("布尔值类型错误应该返回错误")
	}

	// 测试数组类型错误
	invalidParams4 := map[string]interface{}{
		"name":    "测试名称",
		"tags":    "tag1,tag2", // 应该是数组
		"enabled": true,
	}

	err = validateParams(schema, invalidParams4)
	if err == nil {
		t.Error("数组类型错误应该返回错误")
	}

	// 测试对象类型错误
	invalidParams5 := map[string]interface{}{
		"name":     "测试名称",
		"settings": "key=value", // 应该是对象
		"enabled":  true,
	}

	err = validateParams(schema, invalidParams5)
	if err == nil {
		t.Error("对象类型错误应该返回错误")
	}
}

// TestTool 测试工具结构体的功能
func TestTool(t *testing.T) {
	// 定义validateParams函数
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
	// 创建一个测试工具
	tool := Tool{
		Name:        "test_tool",
		Description: "测试工具",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"param1": map[string]any{
					"type":        "string",
					"description": "参数1",
				},
				"param2": map[string]any{
					"type":        "number",
					"description": "参数2",
				},
			},
			Required: []string{"param1"},
		},
	}

	// 测试工具属性
	if tool.Name != "test_tool" {
		t.Errorf("工具名称错误，期望: test_tool, 实际: %s", tool.Name)
	}

	if tool.Description != "测试工具" {
		t.Errorf("工具描述错误，期望: 测试工具, 实际: %s", tool.Description)
	}

	if len(tool.InputSchema.Properties) != 2 {
		t.Errorf("工具参数数量错误，期望: 2, 实际: %d", len(tool.InputSchema.Properties))
	}

	// 测试JSON序列化
	toolJSON, err := json.Marshal(tool)
	if err != nil {
		t.Errorf("工具序列化失败: %v", err)
	}

	// 测试JSON反序列化
	var deserializedTool Tool
	err = json.Unmarshal(toolJSON, &deserializedTool)
	if err != nil {
		t.Errorf("工具反序列化失败: %v", err)
	}

	// 验证反序列化结果
	if deserializedTool.Name != tool.Name {
		t.Errorf("工具名称不匹配，期望: %s, 实际: %s", tool.Name, deserializedTool.Name)
	}

	if deserializedTool.Description != tool.Description {
		t.Errorf("工具描述不匹配，期望: %s, 实际: %s", tool.Description, deserializedTool.Description)
	}

	if len(deserializedTool.InputSchema.Properties) != len(tool.InputSchema.Properties) {
		t.Errorf("工具参数数量不匹配，期望: %d, 实际: %d", len(tool.InputSchema.Properties), len(deserializedTool.InputSchema.Properties))
	}

	// 测试参数验证
	validParams := map[string]interface{}{
		"param1": "value1",
		"param2": 42,
	}

	// 使用前面定义的validateParams函数
	err = validateParams(tool.InputSchema, validParams)
	if err != nil {
		t.Errorf("有效参数验证失败: %v", err)
	}

	// 测试缺少必需参数
	invalidParams := map[string]interface{}{
		"param2": 42,
	}

	err = validateParams(tool.InputSchema, invalidParams)
	if err == nil {
		t.Error("缺少必需参数应该返回错误")
	}
}
