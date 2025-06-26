// plugin/result.go - 工具调用结果相关类型定义
package plugin

import (
	"fmt"
	"reflect"
)

// validateStructType 验证数据类型是否为结构化数据类型
// 支持的类型：结构体、map、slice、array、指向结构体的指针
func validateStructType(data any) error {
	if data == nil {
		return fmt.Errorf("数据不能为 nil")
	}

	v := reflect.ValueOf(data)
	t := v.Type()

	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return fmt.Errorf("指针不能为 nil")
		}
		// 获取指针指向的类型
		t = t.Elem()
	}

	// 检查是否为支持的结构化数据类型
	switch t.Kind() {
	case reflect.Struct:
		// 结构体类型 - 允许
		return nil
	case reflect.Map:
		// Map类型 - 允许
		return nil
	case reflect.Slice, reflect.Array:
		// 切片和数组类型 - 允许
		return nil
	case reflect.Interface:
		// 接口类型 - 如果底层是结构化数据则允许
		if !v.IsNil() {
			return validateStructType(v.Elem().Interface())
		}
		return fmt.Errorf("接口不能为 nil")
	default:
		return fmt.Errorf("不支持的数据类型: %s，只支持结构体、map、slice、array 或指向结构体的指针", t.Kind())
	}
}

// Result 表示工具调用的基础结果结构
// 所有工具调用结果都应该包含这个基础结构
type Result struct {
	// Meta 元数据属性由协议保留，允许客户端和服务器向其响应附加额外的元数据
	Meta map[string]any `json:"_meta,omitempty"`
}

// ContentType 定义内容类型的枚举
type ContentType string

const (
	// ContentTypeText 文本内容类型
	ContentTypeText ContentType = "text"
	// ContentTypeFile 文件内容类型（包括图片、音频、文档等所有文件）
	ContentTypeFile ContentType = "file"
	// ContentTypeStruct 结构体内容类型
	ContentTypeStruct ContentType = "struct"
)

// Content 定义了工具调用结果中内容的接口
// 所有内容类型都应该实现此接口
type Content interface {
	// GetType 返回内容的类型
	GetType() ContentType
}

// TextContent 表示文本内容
type TextContent struct {
	Type ContentType `json:"type"`           // 内容类型，固定为 "text"
	Text string      `json:"text"`           // 文本内容
	Name string      `json:"name,omitempty"` // 内容名称（可选）
}

// GetType 返回文本内容的类型
func (tc TextContent) GetType() ContentType {
	return ContentTypeText
}

// FileType 定义文件类型的枚举
type FileType string

const (
	// FileTypeImage 图片文件
	FileTypeImage FileType = "image"
	// FileTypeAudio 音频文件
	FileTypeAudio FileType = "audio"
	// FileTypeVideo 视频文件
	FileTypeVideo FileType = "video"
	// FileTypeDocument 文档文件
	FileTypeDocument FileType = "document"
	// FileTypeArchive 压缩文件
	FileTypeArchive FileType = "archive"
	// FileTypeCode 代码文件
	FileTypeCode FileType = "code"
	// FileTypeData 数据文件
	FileTypeData FileType = "data"
	// FileTypeOther 其他文件
	FileTypeOther FileType = "other"
)

// FileContent 表示文件内容（统一的文件类型，包括图片、音频、文档等）
type FileContent struct {
	Type     ContentType `json:"type"`           // 内容类型，固定为 "file"
	FileType FileType    `json:"fileType"`       // 文件类型
	Data     string      `json:"data"`           // 文件数据（Base64编码）
	MimeType string      `json:"mimeType"`       // MIME类型
	Name     string      `json:"name,omitempty"` // 文件名称（可选）
	Size     int64       `json:"size,omitempty"` // 文件大小（字节）（可选）

	// 图片特定属性
	Width  int `json:"width,omitempty"`  // 图片宽度（仅图片文件）
	Height int `json:"height,omitempty"` // 图片高度（仅图片文件）

	// 音频/视频特定属性
	Duration float64 `json:"duration,omitempty"` // 音频/视频时长（秒）（仅音频/视频文件）
	Bitrate  int     `json:"bitrate,omitempty"`  // 比特率（仅音频/视频文件）

	// 文档特定属性
	PageCount int    `json:"pageCount,omitempty"` // 页数（仅文档文件）
	Author    string `json:"author,omitempty"`    // 作者（仅文档文件）

	// 通用文件属性
	Encoding string         `json:"encoding,omitempty"` // 文件编码
	Checksum string         `json:"checksum,omitempty"` // 文件校验和
	URL      string         `json:"url,omitempty"`      // 文件URL（可选）
	Metadata map[string]any `json:"metadata,omitempty"` // 额外元数据
}

// GetType 返回文件内容的类型
func (fc FileContent) GetType() ContentType {
	return ContentTypeFile
}

// SetImageProperties 设置图片属性
func (fc FileContent) SetImageProperties(width, height int) FileContent {
	fc.Width = width
	fc.Height = height
	return fc
}

// SetMediaProperties 设置媒体属性（音频/视频）
func (fc FileContent) SetMediaProperties(duration float64, bitrate int) FileContent {
	fc.Duration = duration
	fc.Bitrate = bitrate
	return fc
}

// SetDocumentProperties 设置文档属性
func (fc FileContent) SetDocumentProperties(pageCount int, author string) FileContent {
	fc.PageCount = pageCount
	fc.Author = author
	return fc
}

// SetFileProperties 设置通用文件属性
func (fc FileContent) SetFileProperties(size int64, encoding, checksum string) FileContent {
	fc.Size = size
	fc.Encoding = encoding
	fc.Checksum = checksum
	return fc
}

// SetFileURL 设置文件URL
func (fc FileContent) SetFileURL(url string) FileContent {
	fc.URL = url
	return fc
}

// SetFileMetadata 设置文件元数据
func (fc FileContent) SetFileMetadata(metadata map[string]any) FileContent {
	fc.Metadata = metadata
	return fc
}

// StructContent 表示结构体内容
type StructContent struct {
	Type   ContentType `json:"type"`             // 内容类型，固定为 "struct"
	Data   any         `json:"data"`             // 结构体数据（可以是任意类型）
	Name   string      `json:"name,omitempty"`   // 结构体名称（可选）
	Schema string      `json:"schema,omitempty"` // 结构体模式定义（可选）
	Format string      `json:"format,omitempty"` // 数据格式（如 json、yaml 等）
}

// GetType 返回结构体内容的类型
func (sc StructContent) GetType() ContentType {
	return ContentTypeStruct
}

// SetStructSchema 设置结构体模式定义
func (sc StructContent) SetStructSchema(schema string) StructContent {
	sc.Schema = schema
	return sc
}

// SetStructFormat 设置结构体数据格式
func (sc StructContent) SetStructFormat(format string) StructContent {
	sc.Format = format
	return sc
}

// CallToolResult 表示服务器对工具调用的响应
//
// 任何来自工具的错误都应该在结果对象内报告，将 `isError` 设置为 true，
// 而不是作为 MCP 协议级错误响应。否则，LLM 将无法看到发生了错误并自我纠正。
//
// 但是，在查找工具时的任何错误、表示服务器不支持工具调用的错误，
// 或任何其他异常情况，都应该作为 MCP 错误响应报告。
type CallToolResult struct {
	Result
	// Content 工具调用返回的内容数组，可以包含文本、文件或结构体内容
	Content []Content `json:"content"`
	// IsError 表示工具调用是否以错误结束
	// 如果未设置，则假定为 false（调用成功）
	IsError bool `json:"isError,omitempty"`
}

// NewCallToolResult 创建一个新的工具调用结果
func NewCallToolResult() *CallToolResult {
	return &CallToolResult{
		Content: make([]Content, 0),
		IsError: false,
	}
}

// AddTextContent 向结果中添加文本内容
func (ctr *CallToolResult) AddTextContent(text string, name ...string) *CallToolResult {
	content := TextContent{
		Type: ContentTypeText,
		Text: text,
	}
	if len(name) > 0 {
		content.Name = name[0]
	}
	ctr.Content = append(ctr.Content, content)
	return ctr
}

// AddFileContent 向结果中添加文件内容
func (ctr *CallToolResult) AddFileContent(fileType FileType, data, mimeType string, name ...string) *CallToolResult {
	content := FileContent{
		Type:     ContentTypeFile,
		FileType: fileType,
		Data:     data,
		MimeType: mimeType,
	}
	if len(name) > 0 {
		content.Name = name[0]
	}
	ctr.Content = append(ctr.Content, content)
	return ctr
}

// AddImageContent 向结果中添加图片内容（便捷方法）
func (ctr *CallToolResult) AddImageContent(data, mimeType string, name ...string) *CallToolResult {
	return ctr.AddFileContent(FileTypeImage, data, mimeType, name...)
}

// AddAudioContent 向结果中添加音频内容（便捷方法）
func (ctr *CallToolResult) AddAudioContent(data, mimeType string, name ...string) *CallToolResult {
	return ctr.AddFileContent(FileTypeAudio, data, mimeType, name...)
}

// AddVideoContent 向结果中添加视频内容（便捷方法）
func (ctr *CallToolResult) AddVideoContent(data, mimeType string, name ...string) *CallToolResult {
	return ctr.AddFileContent(FileTypeVideo, data, mimeType, name...)
}

// AddDocumentContent 向结果中添加文档内容（便捷方法）
func (ctr *CallToolResult) AddDocumentContent(data, mimeType string, name ...string) *CallToolResult {
	return ctr.AddFileContent(FileTypeDocument, data, mimeType, name...)
}

// AddStructContent 向结果中添加结构体内容
// data 必须是结构化数据类型：自定义结构体、map、slice、array 或指向结构体的指针
// 如果传入不支持的类型，方法会panic并提供错误信息
func (ctr *CallToolResult) AddStructContent(data any, name ...string) *CallToolResult {
	if err := validateStructType(data); err != nil {
		panic(fmt.Sprintf("AddStructContent: %v", err))
	}

	content := StructContent{
		Type: ContentTypeStruct,
		Data: data,
	}
	if len(name) > 0 {
		content.Name = name[0]
	}
	ctr.Content = append(ctr.Content, content)
	return ctr
}

// SetError 设置工具调用为错误状态
func (ctr *CallToolResult) SetError(isError bool) *CallToolResult {
	ctr.IsError = isError
	return ctr
}

// SetMeta 设置元数据
func (ctr *CallToolResult) SetMeta(key string, value any) *CallToolResult {
	if ctr.Meta == nil {
		ctr.Meta = make(map[string]any)
	}
	ctr.Meta[key] = value
	return ctr
}

// NewTextContent 创建新的文本内容
func NewTextContent(text string, name ...string) TextContent {
	content := TextContent{
		Type: ContentTypeText,
		Text: text,
	}
	if len(name) > 0 {
		content.Name = name[0]
	}
	return content
}

// NewFileContent 创建新的文件内容
func NewFileContent(fileType FileType, data, mimeType string, name ...string) FileContent {
	content := FileContent{
		Type:     ContentTypeFile,
		FileType: fileType,
		Data:     data,
		MimeType: mimeType,
	}
	if len(name) > 0 {
		content.Name = name[0]
	}
	return content
}

// NewImageContent 创建新的图片内容（便捷方法）
func NewImageContent(data, mimeType string, name ...string) FileContent {
	return NewFileContent(FileTypeImage, data, mimeType, name...)
}

// NewAudioContent 创建新的音频内容（便捷方法）
func NewAudioContent(data, mimeType string, name ...string) FileContent {
	return NewFileContent(FileTypeAudio, data, mimeType, name...)
}

// NewVideoContent 创建新的视频内容（便捷方法）
func NewVideoContent(data, mimeType string, name ...string) FileContent {
	return NewFileContent(FileTypeVideo, data, mimeType, name...)
}

// NewDocumentContent 创建新的文档内容（便捷方法）
func NewDocumentContent(data, mimeType string, name ...string) FileContent {
	return NewFileContent(FileTypeDocument, data, mimeType, name...)
}

// NewStructContent 创建新的结构体内容
// data 必须是结构化数据类型：自定义结构体、map、slice、array 或指向结构体的指针
// 如果传入不支持的类型，函数会panic并提供错误信息
func NewStructContent(data any, name ...string) StructContent {
	if err := validateStructType(data); err != nil {
		panic(fmt.Sprintf("NewStructContent: %v", err))
	}

	content := StructContent{
		Type: ContentTypeStruct,
		Data: data,
	}
	if len(name) > 0 {
		content.Name = name[0]
	}
	return content
}

// NewErrorResult 创建一个表示错误的工具调用结果
// errorMessage: 错误信息
func NewErrorResult(errorMessage string) *CallToolResult {
	return &CallToolResult{
		Content: []Content{
			TextContent{
				Type: ContentTypeText,
				Text: errorMessage,
				Name: "error",
			},
		},
		IsError: true,
	}
}
