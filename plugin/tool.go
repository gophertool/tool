// plugin/tool.go - 工具相关类型定义和工具选项函数
package plugin

import "encoding/json"

// Tool 表示一个工具的完整定义
// 包含工具的名称、描述和输入参数模式
type Tool struct {
	Name           string          `json:"name"`         // 工具名称
	Description    string          `json:"description"`  // 工具描述
	InputSchema    ToolInputSchema `json:"input_schema"` // 工具输入参数 与 RawInputSchema 二选一
	RawInputSchema json.RawMessage `json:"-"`            // 工具输入参数的原始JSON Schema 与 InputSchema 二选一
}

// ToolInputSchema 表示工具输入参数的JSON Schema结构
type ToolInputSchema struct {
	Type       string         `json:"type"`                 // 参数类型
	Properties map[string]any `json:"properties,omitempty"` // 参数属性
	Required   []string       `json:"required,omitempty"`   // 参数必填
}

// MarshalJSON 自定义ToolInputSchema的JSON序列化方法
func (tis *ToolInputSchema) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"type": tis.Type,
	}
	if len(tis.Properties) > 0 {
		data["properties"] = tis.Properties
	}
	if len(tis.Required) > 0 {
		data["required"] = tis.Required
	}
	return json.Marshal(data)
}

// ToolOption 是配置Tool的函数选项类型
type ToolOption func(*Tool)

// PropertyOption 是配置Tool输入模式中属性的函数选项类型
// 它允许使用函数选项模式灵活配置JSON Schema属性
type PropertyOption func(map[string]any)

// WithName 设置工具名称的选项函数
func WithName(name string) ToolOption {
	return func(t *Tool) {
		t.Name = name
	}
}

// WithDescription 设置工具描述的选项函数
func WithDescription(description string) ToolOption {
	return func(t *Tool) {
		t.Description = description
	}
}

// Description 设置属性描述的选项函数
func Description(desc string) PropertyOption {
	return func(schema map[string]any) {
		schema["description"] = desc
	}
}

// Required 标记属性为必需的
func Required() PropertyOption {
	return func(schema map[string]any) {
		schema["required"] = true
	}
}

// Default 设置属性的默认值
func Default(value any) PropertyOption {
	return func(schema map[string]any) {
		schema["default"] = value
	}
}

// Enum 设置属性的枚举值
func Enum(values ...any) PropertyOption {
	return func(schema map[string]any) {
		schema["enum"] = values
	}
}

// 字符串相关的属性选项

// MinLength 设置字符串的最小长度
func MinLength(min int) PropertyOption {
	return func(schema map[string]any) {
		schema["minLength"] = min
	}
}

// MaxLength 设置字符串的最大长度
func MaxLength(max int) PropertyOption {
	return func(schema map[string]any) {
		schema["maxLength"] = max
	}
}

// Pattern 设置字符串的正则表达式模式
func Pattern(pattern string) PropertyOption {
	return func(schema map[string]any) {
		schema["pattern"] = pattern
	}
}

// Format 设置字符串的格式（如 email、uri 等）
func Format(format string) PropertyOption {
	return func(schema map[string]any) {
		schema["format"] = format
	}
}

// 数字相关的属性选项

// Minimum 设置数字的最小值
func Minimum(min float64) PropertyOption {
	return func(schema map[string]any) {
		schema["minimum"] = min
	}
}

// Maximum 设置数字的最大值
func Maximum(max float64) PropertyOption {
	return func(schema map[string]any) {
		schema["maximum"] = max
	}
}

// ExclusiveMinimum 设置数字的排他最小值
func ExclusiveMinimum(min float64) PropertyOption {
	return func(schema map[string]any) {
		schema["exclusiveMinimum"] = min
	}
}

// ExclusiveMaximum 设置数字的排他最大值
func ExclusiveMaximum(max float64) PropertyOption {
	return func(schema map[string]any) {
		schema["exclusiveMaximum"] = max
	}
}

// MultipleOf 设置数字必须是指定值的倍数
func MultipleOf(value float64) PropertyOption {
	return func(schema map[string]any) {
		schema["multipleOf"] = value
	}
}

// WithString 添加字符串类型属性的选项函数
func WithString(name string, opts ...PropertyOption) ToolOption {
	return withType("string", name, opts...)
}

// WithNumber 添加数字类型属性的选项函数
func WithNumber(name string, opts ...PropertyOption) ToolOption {
	return withType("number", name, opts...)
}

// WithInteger 添加整数类型属性的选项函数
func WithInteger(name string, opts ...PropertyOption) ToolOption {
	return withType("integer", name, opts...)
}

// WithBoolean 添加布尔类型属性的选项函数
func WithBoolean(name string, opts ...PropertyOption) ToolOption {
	return withType("boolean", name, opts...)
}

// WithObject 添加对象类型属性的选项函数
func WithObject(name string, opts ...PropertyOption) ToolOption {
	return withType("object", name, opts...)
}

// WithArray 添加数组类型属性的选项函数
func WithArray(name string, opts ...PropertyOption) ToolOption {
	return withType("array", name, opts...)
}

// withType 内部函数，用于创建指定类型的属性选项
func withType(propertyType string, name string, opts ...PropertyOption) ToolOption {
	return func(t *Tool) {
		// 确保 Properties 已初始化
		if t.InputSchema.Properties == nil {
			t.InputSchema.Properties = make(map[string]any)
		}

		schema := map[string]any{
			"type": propertyType,
		}

		// 如果是对象类型，初始化 properties
		if propertyType == "object" {
			schema["properties"] = map[string]any{}
		}

		for _, opt := range opts {
			opt(schema)
		}

		// 从属性模式中移除required并添加到InputSchema.required
		if required, ok := schema["required"].(bool); ok && required {
			delete(schema, "required")
			t.InputSchema.Required = append(t.InputSchema.Required, name)
		}

		t.InputSchema.Properties[name] = schema
	}
}

// Properties 定义对象模式的属性
func Properties(props map[string]any) PropertyOption {
	return func(schema map[string]any) {
		schema["properties"] = props
	}
}

// AdditionalProperties 指定对象中是否允许额外属性
// 或为额外属性定义模式
func AdditionalProperties(schema any) PropertyOption {
	return func(schemaMap map[string]any) {
		schemaMap["additionalProperties"] = schema
	}
}

// MinProperties 设置对象的最小属性数量
func MinProperties(min int) PropertyOption {
	return func(schema map[string]any) {
		schema["minProperties"] = min
	}
}

// MaxProperties 设置对象的最大属性数量
func MaxProperties(max int) PropertyOption {
	return func(schema map[string]any) {
		schema["maxProperties"] = max
	}
}

// PropertyNames 为对象中的属性名称定义模式
func PropertyNames(schema map[string]any) PropertyOption {
	return func(schemaMap map[string]any) {
		schemaMap["propertyNames"] = schema
	}
}

// Items 定义数组项目的模式
// 接受任何模式定义以获得最大灵活性
func Items(schema any) PropertyOption {
	return func(schemaMap map[string]any) {
		schemaMap["items"] = schema
	}
}

// MinItems 设置数组的最小项目数量
func MinItems(min int) PropertyOption {
	return func(schema map[string]any) {
		schema["minItems"] = min
	}
}

// MaxItems 设置数组的最大项目数量
func MaxItems(max int) PropertyOption {
	return func(schema map[string]any) {
		schema["maxItems"] = max
	}
}

// UniqueItems 指定数组项目是否必须唯一
func UniqueItems(unique bool) PropertyOption {
	return func(schema map[string]any) {
		schema["uniqueItems"] = unique
	}
}

// WithStringItems 配置数组的项目为字符串类型
// 支持的选项：Description(), Default(), Enum(), MaxLength(), MinLength(), Pattern()
// 注意：Required() 等选项对项目模式无效，将被忽略
func WithStringItems(opts ...PropertyOption) PropertyOption {
	return func(schema map[string]any) {
		itemSchema := map[string]any{
			"type": "string",
		}
		for _, opt := range opts {
			opt(itemSchema)
		}
		schema["items"] = itemSchema
	}
}

// WithStringEnumItems 配置数组的项目为指定枚举的字符串类型
func WithStringEnumItems(values []string) PropertyOption {
	return func(schema map[string]any) {
		schema["items"] = map[string]any{
			"type": "string",
			"enum": values,
		}
	}
}

// WithNumberItems 配置数组的项目为数字类型
// 支持的选项：Description(), Default(), Minimum(), Maximum(), MultipleOf()
func WithNumberItems(opts ...PropertyOption) PropertyOption {
	return func(schema map[string]any) {
		itemSchema := map[string]any{
			"type": "number",
		}
		for _, opt := range opts {
			opt(itemSchema)
		}
		schema["items"] = itemSchema
	}
}

// WithIntegerItems 配置数组的项目为整数类型
func WithIntegerItems(opts ...PropertyOption) PropertyOption {
	return func(schema map[string]any) {
		itemSchema := map[string]any{
			"type": "integer",
		}
		for _, opt := range opts {
			opt(itemSchema)
		}
		schema["items"] = itemSchema
	}
}

// WithBooleanItems 配置数组的项目为布尔类型
// 支持的选项：Description(), Default()
func WithBooleanItems(opts ...PropertyOption) PropertyOption {
	return func(schema map[string]any) {
		itemSchema := map[string]any{
			"type": "boolean",
		}
		for _, opt := range opts {
			opt(itemSchema)
		}
		schema["items"] = itemSchema
	}
}

// WithObjectItems 配置数组的项目为对象类型
func WithObjectItems(opts ...PropertyOption) PropertyOption {
	return func(schema map[string]any) {
		itemSchema := map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		}
		for _, opt := range opts {
			opt(itemSchema)
		}
		schema["items"] = itemSchema
	}
}

// NewTool 创建一个新的工具实例
func NewTool(name, description string, options ...ToolOption) *Tool {
	tool := &Tool{
		Name:        name,
		Description: description,
		InputSchema: ToolInputSchema{
			Type:       "object",
			Properties: make(map[string]any),
			Required:   make([]string, 0),
		},
	}

	for _, option := range options {
		option(tool)
	}

	return tool
}
