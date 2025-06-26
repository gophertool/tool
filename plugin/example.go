// plugin/example.go - 使用示例
package plugin

import (
	"encoding/json"
	"fmt"
)

// ExampleTool 创建工具的示例
func ExampleTool() {
	// 创建一个文件处理工具
	tool := NewTool(
		"file_processor",
		"处理文件的工具，支持读取、写入和转换操作",
		WithString("filepath",
			Description("文件路径"),
			Required(),
			Pattern(`^[a-zA-Z0-9_/.-]+$`),
		),
		WithString("operation",
			Description("要执行的操作"),
			Required(),
			Enum("read", "write", "convert"),
		),
		WithString("content",
			Description("写入的内容（仅当operation为write时需要）"),
			MinLength(1),
			MaxLength(10000),
		),
		WithObject("options",
			Description("操作选项"),
			Properties(map[string]any{
				"encoding": map[string]any{
					"type":    "string",
					"default": "utf-8",
					"enum":    []string{"utf-8", "gbk", "ascii"},
				},
				"backup": map[string]any{
					"type":    "boolean",
					"default": false,
				},
			}),
		),
		WithArray("tags",
			Description("文件标签"),
			WithStringItems(
				MinLength(1),
				MaxLength(50),
			),
			MinItems(0),
			MaxItems(10),
		),
		WithArray("permissions",
			Description("权限设置"),
			WithStringEnumItems([]string{"read", "write", "execute"}),
		),
		WithNumber("size_limit",
			Description("文件大小限制（MB）"),
			Minimum(0),
			Maximum(1024),
			Default(100),
		),
		WithInteger("timeout",
			Description("超时时间（秒）"),
			Minimum(1),
			Maximum(300),
			Default(30),
		),
		WithBoolean("debug",
			Description("是否启用调试模式"),
			Default(false),
		),
	)

	// 输出工具的JSON定义
	toolJSON, _ := json.MarshalIndent(tool, "", "  ")
	fmt.Printf("工具定义:\n%s\n", toolJSON)
}

// ExampleCallToolResult 创建工具调用结果的示例
func ExampleCallToolResult() {
	// 成功的工具调用结果
	successResult := NewCallToolResult().
		AddTextContent("文件处理完成", "status").
		AddTextContent("共处理了 1024 个字符", "details").
		SetMeta("processing_time", "150ms").
		SetMeta("version", "1.0")

	// 包含图片的结果
	imageResult := NewCallToolResult().
		AddTextContent("图片处理完成").
		AddImageContent("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
			"image/png", "result_image")

	// 错误的工具调用结果
	errorResult := NewCallToolResult().
		AddTextContent("文件不存在: /path/to/nonexistent/file.txt", "error").
		SetError(true).
		SetMeta("error_code", "FILE_NOT_FOUND")

	// 包含多种内容类型的复杂结果
	complexResult := NewCallToolResult().
		AddTextContent("数据分析完成", "summary").
		AddImageContent("base64_chart_data", "image/png", "chart").
		AddAudioContent("base64_audio_data", "audio/wav", "report_audio").
		AddVideoContent("base64_video_data", "video/mp4", "demo_video").
		AddDocumentContent("base64_pdf_data", "application/pdf", "report.pdf").
		AddStructContent(map[string]any{
			"total_records":   10000,
			"processing_time": "2.5s",
			"success_rate":    98.5,
			"categories":      []string{"success", "warning", "error"},
			"metadata": map[string]any{
				"version":   "1.0",
				"timestamp": "2024-01-15T10:30:00Z",
			},
		}, "analysis_data")

	// 输出结果的JSON
	results := []struct {
		name   string
		result *CallToolResult
	}{
		{"成功结果", successResult},
		{"图片结果", imageResult},
		{"错误结果", errorResult},
		{"复杂结果", complexResult},
	}

	for _, r := range results {
		resultJSON, _ := json.MarshalIndent(r.result, "", "  ")
		fmt.Printf("\n%s:\n%s\n", r.name, resultJSON)
	}
}

// ExampleAdvancedTool 高级工具配置示例
func ExampleAdvancedTool() {
	// 创建一个数据处理工具，展示更复杂的配置
	tool := NewTool(
		"data_processor",
		"高级数据处理工具，支持多种数据格式和转换操作",
		WithObject("input",
			Description("输入数据配置"),
			Required(),
			Properties(map[string]any{
				"source": map[string]any{
					"type":        "string",
					"description": "数据源类型",
					"enum":        []string{"file", "database", "api", "stream"},
				},
				"format": map[string]any{
					"type":        "string",
					"description": "数据格式",
					"enum":        []string{"json", "csv", "xml", "parquet"},
				},
				"connection": map[string]any{
					"type":        "object",
					"description": "连接配置",
					"properties": map[string]any{
						"host": map[string]any{"type": "string"},
						"port": map[string]any{"type": "integer", "minimum": 1, "maximum": 65535},
						"credentials": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"username": map[string]any{"type": "string"},
								"password": map[string]any{"type": "string"},
							},
						},
					},
				},
			}),
		),
		WithArray("transformations",
			Description("数据转换操作列表"),
			WithObjectItems(
				Properties(map[string]any{
					"type": map[string]any{
						"type":        "string",
						"description": "转换类型",
						"enum":        []string{"filter", "map", "aggregate", "join"},
					},
					"config": map[string]any{
						"type":        "object",
						"description": "转换配置参数",
					},
				}),
			),
			MinItems(1),
			MaxItems(20),
		),
		WithObject("output",
			Description("输出配置"),
			Required(),
			Properties(map[string]any{
				"destination": map[string]any{
					"type":        "string",
					"description": "输出目标",
					"enum":        []string{"file", "database", "console", "memory"},
				},
				"format": map[string]any{
					"type":        "string",
					"description": "输出格式",
					"enum":        []string{"json", "csv", "xml", "parquet"},
				},
			}),
		),
		WithArray("validation_rules",
			Description("数据验证规则"),
			WithStringItems(
				Description("验证规则表达式"),
				Pattern(`^[a-zA-Z0-9_\s\(\)\.\>\<\=\!\&\|]+$`),
			),
			UniqueItems(true),
		),
		WithBoolean("parallel_processing",
			Description("是否启用并行处理"),
			Default(true),
		),
		WithInteger("batch_size",
			Description("批处理大小"),
			Minimum(1),
			Maximum(10000),
			Default(1000),
		),
		WithNumber("memory_limit_gb",
			Description("内存使用限制（GB）"),
			Minimum(0.1),
			Maximum(64.0),
			Default(2.0),
		),
	)

	toolJSON, _ := json.MarshalIndent(tool, "", "  ")
	fmt.Printf("高级工具定义:\n%s\n", toolJSON)
}

// ExampleContentCreation 展示如何创建不同类型的内容
func ExampleContentCreation() {
	// 创建不同类型的内容
	textContent := NewTextContent("这是一个文本内容示例", "example_text")

	// 创建图片文件内容
	imageContent := NewImageContent(
		"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
		"image/png",
		"sample_image",
	).SetImageProperties(100, 100).SetFileProperties(1024, "binary", "sha256:abc123")

	// 创建音频文件内容
	audioContent := NewAudioContent(
		"UklGRiQAAABXQVZFZm10IBAAAAABAAEARKwAAIhYAQACABAAZGF0YQAAAAA=",
		"audio/wav",
		"sample_audio",
	).SetMediaProperties(5.5, 128).SetFileProperties(2048, "binary", "sha256:def456")

	// 创建视频文件内容
	videoContent := NewVideoContent(
		"AAAAIGZ0eXBpc29tAAACAGlzb21pc28yYXZjMW1wNDEAAAAIZnJlZQAAAr5tZGF0",
		"video/mp4",
		"sample_video",
	).SetMediaProperties(30.0, 1000).SetImageProperties(1920, 1080)

	// 创建文档文件内容
	documentContent := NewDocumentContent(
		"JVBERi0xLjQKJdPr6eEKMSAwIG9iago8PAovVHlwZSAvQ2F0YWxvZwovUGFnZXMgMiAwIFIKPj4KZW5kb2JqCg==",
		"application/pdf",
		"sample_document",
	).SetDocumentProperties(10, "John Doe").SetFileProperties(5120, "binary", "sha256:ghi789")

	// 创建结构体内容
	structContent := NewStructContent(
		map[string]any{
			"user": map[string]any{
				"id":    12345,
				"name":  "张三",
				"email": "zhangsan@example.com",
				"roles": []string{"admin", "user"},
			},
			"settings": map[string]any{
				"theme":    "dark",
				"language": "zh-CN",
				"notifications": map[string]any{
					"email": true,
					"sms":   false,
				},
			},
			"metadata": map[string]any{
				"created_at": "2024-01-15T10:30:00Z",
				"version":    "1.0",
			},
		},
		"user_profile",
	).SetStructSchema("UserProfile").SetStructFormat("json")

	// 使用这些内容创建结果
	result := NewCallToolResult()
	result.Content = []Content{
		textContent,
		imageContent,
		audioContent,
		videoContent,
		documentContent,
		structContent,
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("内容创建示例:\n%s\n", resultJSON)
}

// ExampleFileTypes 展示不同文件类型的使用
func ExampleFileTypes() {
	result := NewCallToolResult()

	// 图片文件
	result.AddFileContent(FileTypeImage, "base64_image_data", "image/jpeg", "photo.jpg")

	// 音频文件
	result.AddFileContent(FileTypeAudio, "base64_audio_data", "audio/mp3", "song.mp3")

	// 视频文件
	result.AddFileContent(FileTypeVideo, "base64_video_data", "video/mp4", "video.mp4")

	// 文档文件
	result.AddFileContent(FileTypeDocument, "base64_pdf_data", "application/pdf", "report.pdf")

	// 压缩文件
	result.AddFileContent(FileTypeArchive, "base64_zip_data", "application/zip", "archive.zip")

	// 代码文件
	result.AddFileContent(FileTypeCode, "Y29uc29sZS5sb2coIkhlbGxvLCBXb3JsZCEiKTs=", "text/javascript", "script.js")

	// 数据文件
	result.AddFileContent(FileTypeData, "base64_csv_data", "text/csv", "data.csv")

	// 其他文件
	result.AddFileContent(FileTypeOther, "base64_unknown_data", "application/octet-stream", "unknown.bin")

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("文件类型示例:\n%s\n", resultJSON)
}

// ExampleStructTypes 展示结构体类型的使用
func ExampleStructTypes() {
	result := NewCallToolResult()

	// 用户信息结构体
	result.AddStructContent(map[string]any{
		"id":       123,
		"username": "john_doe",
		"profile": map[string]any{
			"firstName": "John",
			"lastName":  "Doe",
			"age":       30,
			"hobbies":   []string{"reading", "coding", "music"},
		},
	}, "user_info")

	// 系统配置结构体
	result.AddStructContent(map[string]any{
		"database": map[string]any{
			"host":    "localhost",
			"port":    5432,
			"name":    "myapp",
			"ssl":     true,
			"timeout": 30,
		},
		"cache": map[string]any{
			"type":    "redis",
			"servers": []string{"127.0.0.1:6379"},
			"ttl":     3600,
		},
		"logging": map[string]any{
			"level":   "info",
			"format":  "json",
			"outputs": []string{"console", "file"},
		},
	}, "system_config")

	// 分析结果结构体
	result.AddStructContent(map[string]any{
		"summary": map[string]any{
			"total_items":  1000,
			"processed":    950,
			"failed":       50,
			"success_rate": 95.0,
		},
		"details": []map[string]any{
			{
				"category": "validation",
				"count":    30,
				"errors":   []string{"missing field", "invalid format"},
			},
			{
				"category": "processing",
				"count":    20,
				"errors":   []string{"timeout", "connection failed"},
			},
		},
		"recommendations": []string{
			"增加输入验证",
			"优化网络连接",
			"添加重试机制",
		},
	}, "analysis_result")

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("结构体类型示例:\n%s\n", resultJSON)
}

// ExampleStructValidation 展示结构体类型验证的使用
func ExampleStructValidation() {
	fmt.Println("=== 结构体类型验证示例 ===")

	// 定义一些测试用的结构体
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type Config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	result := NewCallToolResult()

	// 1. 正确的用法 - 结构体
	user := User{ID: 1, Name: "张三", Age: 30}
	result.AddStructContent(user, "user_struct")
	fmt.Println("✓ 添加结构体成功")

	// 2. 正确的用法 - 指向结构体的指针
	config := &Config{Host: "localhost", Port: 8080}
	result.AddStructContent(config, "config_pointer")
	fmt.Println("✓ 添加结构体指针成功")

	// 3. 正确的用法 - map
	mapData := map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": []string{"a", "b", "c"},
	}
	result.AddStructContent(mapData, "map_data")
	fmt.Println("✓ 添加map成功")

	// 4. 正确的用法 - slice
	sliceData := []User{
		{ID: 1, Name: "用户1", Age: 25},
		{ID: 2, Name: "用户2", Age: 30},
	}
	result.AddStructContent(sliceData, "slice_data")
	fmt.Println("✓ 添加slice成功")

	// 5. 正确的用法 - array
	arrayData := [3]string{"item1", "item2", "item3"}
	result.AddStructContent(arrayData, "array_data")
	fmt.Println("✓ 添加array成功")

	// 演示错误用法（这些会panic）
	fmt.Println("\n=== 错误用法演示（会panic） ===")

	// 定义一个安全的演示函数
	safeDemo := func(name string, data any) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("✗ %s: %v\n", name, r)
			}
		}()
		result.AddStructContent(data, name)
		fmt.Printf("✓ %s: 成功\n", name)
	}

	// 6. 错误用法 - 基础类型
	safeDemo("string类型", "这是一个字符串")
	safeDemo("int类型", 42)
	safeDemo("bool类型", true)
	safeDemo("float类型", 3.14)

	// 7. 错误用法 - nil
	safeDemo("nil值", nil)

	// 8. 错误用法 - nil指针
	var nilPtr *User
	safeDemo("nil指针", nilPtr)

	// 输出最终结果
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("\n=== 最终结果 ===\n%s\n", resultJSON)
}

// ExampleTextContentWithNames 展示带名称的文本内容使用
func ExampleTextContentWithNames() {
	fmt.Println("=== 带名称的文本内容示例 ===")

	result := NewCallToolResult()

	// 添加带名称的文本内容
	result.AddTextContent("操作开始执行", "start_message").
		AddTextContent("正在验证用户权限...", "auth_step").
		AddTextContent("权限验证通过", "auth_result").
		AddTextContent("正在处理数据...", "processing_step").
		AddTextContent("数据处理完成，共处理 1000 条记录", "processing_result").
		AddTextContent("操作成功完成", "completion_message")

	// 也可以添加没有名称的文本内容
	result.AddTextContent("这是一个没有名称的文本内容")

	// 使用NewTextContent创建带名称的文本内容
	namedText := NewTextContent("这是通过NewTextContent创建的带名称文本", "custom_text")
	unnamedText := NewTextContent("这是通过NewTextContent创建的无名称文本")

	// 手动添加到结果中
	result.Content = append(result.Content, namedText, unnamedText)

	// 设置一些元数据
	result.SetMeta("total_steps", 6).
		SetMeta("execution_time", "2.5s").
		SetMeta("status", "success")

	// 输出结果
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("带名称的文本内容结果:\n%s\n", resultJSON)

	// 展示如何访问具体的内容项
	fmt.Println("\n=== 内容项详情 ===")
	for i, content := range result.Content {
		if textContent, ok := content.(TextContent); ok {
			if textContent.Name != "" {
				fmt.Printf("文本内容 %d: 名称='%s', 内容='%s'\n", i+1, textContent.Name, textContent.Text)
			} else {
				fmt.Printf("文本内容 %d: 无名称, 内容='%s'\n", i+1, textContent.Text)
			}
		}
	}
}
