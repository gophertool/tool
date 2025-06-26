// plugin.go
// 插件系统核心文件，基于 hashicorp/go-plugin 实现
// 提供插件的加载、管理和调用功能
package plugin

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/hashicorp/go-plugin"
)

// init 函数在包被导入时自动执行
// 用于注册所有需要的 gob 类型
func init() {
	// 自动注册所有需要的 gob 类型
	registerGobTypes()
}

// ToolPluginInterface 定义了工具插件需要实现的核心接口
// 这是插件和主程序之间的通信契约
type ToolPluginInterface interface {
	// GetTools 获取插件提供的所有工具定义
	GetTools() ([]Tool, error)

	// CallTool 调用指定的工具并返回结果
	// 这是原始的方法，接受 map[string]interface{} 类型的参数
	CallTool(toolName string, params map[string]interface{}) (*CallToolResult, error)

	// GetPluginInfo 获取插件的基本信息
	GetPluginInfo() (PluginInfo, error)
}

// ToolPluginGenericInterface 定义了支持结构化参数的工具插件接口
// 这是 ToolPluginInterface 的扩展，允许使用结构化参数
type ToolPluginGenericInterface interface {
	// 包含基础接口的所有方法
	ToolPluginInterface

	// CallToolWithStruct 使用结构化参数调用指定的工具
	// params 是一个结构体，必须是可序列化的
	CallToolWithStruct(toolName string, params any) (*CallToolResult, error)
}

// PluginInfo 插件基本信息结构体
type PluginInfo struct {
	Name        string `json:"name"`        // 插件名称
	Version     string `json:"version"`     // 插件版本
	Description string `json:"description"` // 插件描述
	Author      string `json:"author"`      // 插件作者
}

// ToolPluginRPC RPC客户端实现
// 将接口调用转换为跨进程的RPC调用
type ToolPluginRPC struct {
	client *rpc.Client
}

// GetTools 实现 ToolPluginInterface 接口的 GetTools 方法
func (t *ToolPluginRPC) GetTools() ([]Tool, error) {
	var tools []Tool
	err := t.client.Call("Plugin.GetTools", new(interface{}), &tools)
	return tools, err
}

// CallTool 实现 ToolPluginInterface 接口的 CallTool 方法
func (t *ToolPluginRPC) CallTool(toolName string, params map[string]interface{}) (*CallToolResult, error) {
	args := CallToolArgs{
		ToolName: toolName,
		Params:   params,
	}
	var result CallToolResult
	err := t.client.Call("Plugin.CallTool", args, &result)
	return &result, err
}

// CallToolWithStruct 实现 ToolPluginGenericInterface 接口的 CallToolWithStruct 方法
// 使用结构化参数调用工具
func (t *ToolPluginRPC) CallToolWithStruct(toolName string, params interface{}) (*CallToolResult, error) {
	// 在客户端就将结构体转换为 map[string]interface{}
	// 这样可以避免 gob 类型注册问题
	var paramsMap map[string]interface{}

	// 检查是否已经是map类型
	if mapParams, ok := params.(map[string]interface{}); ok {
		paramsMap = mapParams
	} else {
		// 将结构体转换为map
		paramsMap = structToMap(params)
	}

	// 使用标准的 CallTool 方法，传递 map 参数
	return t.CallTool(toolName, paramsMap)
}

// GetPluginInfo 实现 ToolPluginInterface 接口的 GetPluginInfo 方法
func (t *ToolPluginRPC) GetPluginInfo() (PluginInfo, error) {
	var info PluginInfo
	err := t.client.Call("Plugin.GetPluginInfo", new(interface{}), &info)
	return info, err
}

// CallToolArgs 工具调用参数结构体
type CallToolArgs struct {
	ToolName string                 `json:"tool_name"` // 工具名称
	Params   map[string]interface{} `json:"params"`    // 调用参数
}

// StructCallToolArgs 支持结构化参数的工具调用参数结构体
// 注意：由于客户端现在会将结构体转换为map，这个结构体主要用于接口完整性
// 实际使用中参数会是 map[string]interface{} 类型
type StructCallToolArgs struct {
	ToolName string      `json:"tool_name"` // 工具名称
	Params   interface{} `json:"params"`    // 结构化参数（通常是map[string]interface{}）
}

// ToolPluginRPCServer RPC服务器端实现
// 接收RPC调用并转发给实际的插件实现
type ToolPluginRPCServer struct {
	Impl ToolPluginInterface // 实际的插件实现
}

// GetTools 处理来自客户端的 GetTools RPC 调用
func (s *ToolPluginRPCServer) GetTools(args interface{}, resp *[]Tool) error {
	tools, err := s.Impl.GetTools()
	*resp = tools
	return err
}

// CallTool 处理来自客户端的 CallTool RPC 调用
func (s *ToolPluginRPCServer) CallTool(args CallToolArgs, resp *CallToolResult) error {
	result, err := s.Impl.CallTool(args.ToolName, args.Params)
	if err != nil {
		// 创建错误结果
		*resp = *NewErrorResult(fmt.Sprintf("调用工具失败: %v", err))
		return nil
	}
	*resp = *result
	return nil
}

// CallToolWithStruct 处理来自客户端的 CallToolWithStruct RPC 调用
// 这个方法处理结构化参数的工具调用
// 注意：由于客户端已经将结构体转换为map，这个方法现在实际上不会被调用
// 保留此方法是为了接口完整性
func (s *ToolPluginRPCServer) CallToolWithStruct(args StructCallToolArgs, resp *CallToolResult) error {
	// 检查是否实现了泛型接口
	if genericImpl, ok := s.Impl.(ToolPluginGenericInterface); ok {
		// 调用结构化参数方法
		result, err := genericImpl.CallToolWithStruct(args.ToolName, args.Params)
		if err != nil {
			*resp = *NewErrorResult(fmt.Sprintf("调用工具失败: %v", err))
			return nil
		}
		*resp = *result
		return nil
	}

	// 如果没有实现泛型接口，尝试将参数转换为map并使用标准接口
	paramsMap, ok := args.Params.(map[string]interface{})
	if !ok {
		// 如果不是map，尝试将结构体转换为map
		paramsMap = structToMap(args.Params)
	}

	result, err := s.Impl.CallTool(args.ToolName, paramsMap)
	if err != nil {
		*resp = *NewErrorResult(fmt.Sprintf("调用工具失败: %v", err))
		return nil
	}
	*resp = *result
	return nil
}

// GetPluginInfo 处理来自客户端的 GetPluginInfo RPC 调用
func (s *ToolPluginRPCServer) GetPluginInfo(args interface{}, resp *PluginInfo) error {
	info, err := s.Impl.GetPluginInfo()
	*resp = info
	return err
}

// ToolPlugin 实现了 hashicorp/go-plugin 的 Plugin 接口
// 这是插件系统的核心，负责客户端和服务器端的创建
type ToolPlugin struct {
	Impl ToolPluginInterface // 插件的实际实现
}

// Server 返回插件的RPC服务器实现
// 这个方法在插件进程中被调用
func (p *ToolPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ToolPluginRPCServer{Impl: p.Impl}, nil
}

// Client 返回插件的RPC客户端实现
// 这个方法在主程序中被调用，用于与插件通信
func (ToolPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ToolPluginRPC{client: c}, nil
}

// HandshakeConfig 定义了主程序和插件之间的握手配置
// 用于确保版本兼容性和安全性
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,             // 协议版本号
	MagicCookieKey:   "TOOL_PLUGIN", // 魔术cookie的键
	MagicCookieValue: "tool_v1.0.0", // 魔术cookie的值
}

// PluginMap 是插件映射表
// 定义了可用的插件类型
var PluginMap = map[string]plugin.Plugin{
	"tool": &ToolPlugin{}, // 注册工具插件
}

// LoadedPlugin 已加载的插件信息
type LoadedPlugin struct {
	Name     string              // 插件名称
	Path     string              // 插件文件路径
	Client   *plugin.Client      // 插件客户端
	Instance ToolPluginInterface // 插件实例
	Info     PluginInfo          // 插件信息
	Tools    []Tool              // 插件提供的工具
}

// PluginManager 插件管理器
// 负责管理所有已加载的插件
type PluginManager struct {
	mu      sync.RWMutex             // 读写锁
	plugins map[string]*LoadedPlugin // 插件映射表，key为插件名称
	toolMap map[string]*LoadedPlugin // 工具到插件的映射表，key为工具名称
}

// NewPluginManager 创建新的插件管理器
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*LoadedPlugin),
		toolMap: make(map[string]*LoadedPlugin),
	}
}

// ScanPlugins 扫描指定目录下的所有.tool.plugin文件
// pluginDir: 要扫描的插件目录路径
// 返回找到的插件文件路径列表
func (pm *PluginManager) ScanPlugins(pluginDir string) ([]string, error) {
	var pluginPaths []string

	// 使用filepath.Walk递归遍历目录
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查文件是否以.tool.plugin结尾且是可执行文件
		if strings.HasSuffix(path, ".tool.plugin") && !info.IsDir() {
			// 检查文件是否可执行
			if info.Mode()&0111 != 0 {
				pluginPaths = append(pluginPaths, path)
			}
		}

		return nil
	})

	return pluginPaths, err
}

// LoadPlugin 加载单个插件
// pluginPath: 插件文件路径
// 返回加载的插件信息
func (pm *PluginManager) LoadPlugin(pluginPath string) (*LoadedPlugin, error) {
	// 从路径中提取插件名称（去掉目录和.tool.plugin后缀）
	pluginName := filepath.Base(pluginPath)
	pluginName = strings.TrimSuffix(pluginName, ".tool.plugin")

	log.Printf("正在加载插件: %s (路径: %s)", pluginName, pluginPath)

	// 创建插件客户端配置
	config := &plugin.ClientConfig{
		HandshakeConfig:  HandshakeConfig,                          // 握手配置，确保版本兼容
		Plugins:          PluginMap,                                // 插件映射表
		Cmd:              exec.Command(pluginPath),                 // 插件可执行文件命令
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC}, // 允许的协议
	}

	// 创建插件客户端
	client := plugin.NewClient(config)

	// 连接到插件
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("连接插件 %s 失败: %v", pluginName, err)
	}

	// 获取插件实例
	raw, err := rpcClient.Dispense("tool")
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("获取插件实例 %s 失败: %v", pluginName, err)
	}

	// 将插件实例转换为我们的接口类型
	toolPlugin := raw.(ToolPluginInterface)

	// 获取插件信息
	pluginInfo, err := toolPlugin.GetPluginInfo()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("获取插件信息 %s 失败: %v", pluginName, err)
	}

	// 获取插件提供的工具
	tools, err := toolPlugin.GetTools()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("获取插件工具 %s 失败: %v", pluginName, err)
	}

	// 创建已加载插件信息
	loadedPlugin := &LoadedPlugin{
		Name:     pluginName,
		Path:     pluginPath,
		Client:   client,
		Instance: toolPlugin,
		Info:     pluginInfo,
		Tools:    tools,
	}

	log.Printf("插件 %s 加载成功! 提供 %d 个工具", pluginName, len(tools))
	return loadedPlugin, nil
}

// LoadAllPlugins 加载所有扫描到的插件
// pluginDir: 插件目录路径
func (pm *PluginManager) LoadAllPlugins(pluginDir string) error {
	// 扫描插件文件
	pluginPaths, err := pm.ScanPlugins(pluginDir)
	if err != nil {
		return fmt.Errorf("扫描插件目录失败: %v", err)
	}

	if len(pluginPaths) == 0 {
		log.Printf("在目录 %s 中未找到任何.tool.plugin文件", pluginDir)
		return nil
	}

	log.Printf("发现 %d 个插件文件", len(pluginPaths))

	// 逐个加载插件
	var loadedCount int
	var failedCount int

	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, pluginPath := range pluginPaths {
		loadedPlugin, err := pm.LoadPlugin(pluginPath)
		if err != nil {
			log.Printf("加载插件失败: %v", err)
			failedCount++
			continue
		}

		// 将插件添加到管理器中
		pm.plugins[loadedPlugin.Name] = loadedPlugin

		// 建立工具到插件的映射
		for _, tool := range loadedPlugin.Tools {
			pm.toolMap[tool.Name] = loadedPlugin
		}

		loadedCount++
	}

	log.Printf("插件加载结果: 成功 %d 个, 失败 %d 个", loadedCount, failedCount)

	// 如果有插件路径但没有成功加载任何插件，才返回错误
	if len(pluginPaths) > 0 && loadedCount == 0 {
		return fmt.Errorf("没有成功加载任何插件")
	}

	return nil
}

// GetPlugin 获取指定名称的插件
func (pm *PluginManager) GetPlugin(name string) (*LoadedPlugin, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, exists := pm.plugins[name]
	return plugin, exists
}

// GetPluginByTool 根据工具名称获取对应的插件
func (pm *PluginManager) GetPluginByTool(toolName string) (*LoadedPlugin, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, exists := pm.toolMap[toolName]
	return plugin, exists
}

// ListPlugins 列出所有已加载的插件
func (pm *PluginManager) ListPlugins() []*LoadedPlugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugins := make([]*LoadedPlugin, 0, len(pm.plugins))
	for _, plugin := range pm.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// ListTools 列出所有可用的工具
func (pm *PluginManager) ListTools() []Tool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var tools []Tool
	for _, plugin := range pm.plugins {
		tools = append(tools, plugin.Tools...)
	}
	return tools
}

// CallTool 调用指定的工具
func (pm *PluginManager) CallTool(toolName string, params map[string]interface{}) (*CallToolResult, error) {
	// 查找工具对应的插件
	plugin, exists := pm.GetPluginByTool(toolName)
	if !exists {
		return nil, fmt.Errorf("工具 '%s' 不存在", toolName)
	}

	// 调用插件的工具
	return plugin.Instance.CallTool(toolName, params)
}

// structToMap 将任意结构体转换为 map[string]interface{}
// 使用 JSON 序列化和反序列化实现，保留字段标签
func structToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// 首先检查是否为结构体
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 如果不是结构体，返回空 map
	if v.Kind() != reflect.Struct {
		return result
	}

	// 使用 JSON 序列化和反序列化
	// 这样可以保留 JSON 标签信息
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Printf("结构体转换为 JSON 失败: %v", err)
		return result
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Printf("JSON 转换为 map 失败: %v", err)
		return make(map[string]interface{})
	}

	return result
}

// CallToolWithStruct 使用结构化参数调用指定的工具
// 这个方法允许使用强类型参数而不是 map[string]interface{}
func (pm *PluginManager) CallToolWithStruct(toolName string, params interface{}) (*CallToolResult, error) {
	// 获取工具对应的插件
	plugin, found := pm.GetPluginByTool(toolName)
	if !found {
		return nil, fmt.Errorf("找不到工具: %s", toolName)
	}

	// 检查插件是否实现了泛型接口
	if genericPlugin, ok := plugin.Instance.(ToolPluginGenericInterface); ok {
		// 使用泛型接口调用
		return genericPlugin.CallToolWithStruct(toolName, params)
	}

	// 如果插件没有实现泛型接口，尝试将参数转换为map[string]interface{}
	var paramsMap map[string]interface{}

	// 检查是否已经是map类型
	if mapParams, ok := params.(map[string]interface{}); ok {
		paramsMap = mapParams
	} else {
		// 尝试将结构体转换为map
		paramsMap = structToMap(params)
	}

	// 使用标准接口调用
	return plugin.Instance.CallTool(toolName, paramsMap)
}

// CallToolWithContext 带上下文调用指定的工具
func (pm *PluginManager) CallToolWithContext(ctx context.Context, toolName string, params map[string]interface{}) (*CallToolResult, error) {
	// 创建带取消功能的通道
	resultChan := make(chan *CallToolResult, 1)
	errorChan := make(chan error, 1)

	// 在goroutine中执行工具调用
	go func() {
		result, err := pm.CallTool(toolName, params)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()

	// 等待结果或上下文取消
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("工具调用被取消: %w", ctx.Err())
	}
}

// CallToolWithStructContext 使用结构化参数和上下文调用指定的工具
// 这个方法允许使用结构体参数而不是 map[string]interface{}
// 同时支持上下文取消和超时
func (pm *PluginManager) CallToolWithStructContext(ctx context.Context, toolName string, params interface{}) (*CallToolResult, error) {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// 继续执行
	}

	// 查找工具对应的插件
	plugin, exists := pm.GetPluginByTool(toolName)
	if !exists {
		return nil, fmt.Errorf("工具 '%s' 不存在", toolName)
	}

	// 检查插件是否实现了泛型接口
	if genericPlugin, ok := plugin.Instance.(ToolPluginGenericInterface); ok {
		// 使用结构化参数接口调用
		// 注意：这里我们没有传递上下文，因为插件接口没有定义带上下文的方法
		// 但我们已经在调用前检查了上下文状态
		return genericPlugin.CallToolWithStruct(toolName, params)
	}

	// 如果插件没有实现泛型接口，尝试将参数转换为 map[string]interface{}
	var paramsMap map[string]interface{}

	// 检查是否已经是map类型
	if mapParams, ok := params.(map[string]interface{}); ok {
		paramsMap = mapParams
	} else {
		// 尝试将结构体转换为map
		paramsMap = structToMap(params)
	}

	// 使用标准接口调用
	return plugin.Instance.CallTool(toolName, paramsMap)
}

// Shutdown 关闭所有插件
func (pm *PluginManager) Shutdown() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	log.Println("正在关闭所有插件...")
	for name, plugin := range pm.plugins {
		log.Printf("关闭插件: %s", name)
		plugin.Client.Kill()
	}

	// 清空映射表
	pm.plugins = make(map[string]*LoadedPlugin)
	pm.toolMap = make(map[string]*LoadedPlugin)

	log.Println("所有插件已关闭")
}

// registerGobTypes 注册所有需要的 gob 类型，用于 RPC 通信
// 这个函数会在 ServePlugin 中自动调用，插件开发者不需要手动注册
func registerGobTypes() {
	// 注册基础类型
	gob.Register(map[string]any{}) // 注意：这里使用了空的map而不是nil
	gob.Register([]any{})
	gob.Register(map[string]int{})
	gob.Register(map[string]string{})
	gob.Register(map[string]float64{})
	gob.Register(map[string]bool{})
	gob.Register([]string{})
	gob.Register([]int{})
	gob.Register([]float64{})
	gob.Register([]bool{})

	// 注册内容类型
	gob.RegisterName("github.com/gophertool/tool/plugin.TextContent", TextContent{})
	gob.RegisterName("github.com/gophertool/tool/plugin.FileContent", FileContent{})
	gob.RegisterName("github.com/gophertool/tool/plugin.StructContent", StructContent{})
	gob.RegisterName("github.com/gophertool/tool/plugin.CallToolResult", CallToolResult{})
	gob.RegisterName("github.com/gophertool/tool/plugin.Content", []Content{})

	// 注册工具相关类型
	gob.RegisterName("github.com/gophertool/tool/plugin.Tool", Tool{})
	gob.RegisterName("github.com/gophertool/tool/plugin.ToolInputSchema", ToolInputSchema{})
	gob.RegisterName("github.com/gophertool/tool/plugin.PluginInfo", PluginInfo{})
	gob.RegisterName("github.com/gophertool/tool/plugin.CallToolArgs", CallToolArgs{})
	gob.RegisterName("github.com/gophertool/tool/plugin.StructCallToolArgs", StructCallToolArgs{})
}

// RegisterStructType 注册自定义结构体类型，用于 RPC 通信
// 注意：由于客户端现在会自动将结构体转换为map，通常不再需要调用此函数
// 此函数保留用于特殊情况下需要传递原始结构体的场景
// 例如：RegisterStructType(MyCustomStruct{}) 将注册 MyCustomStruct 类型
func RegisterStructType(structType interface{}) {
	// 获取类型信息
	typeName := fmt.Sprintf("%T", structType)
	// 注册类型
	gob.Register(structType)
	log.Printf("已注册结构体类型: %s", typeName)
}

// ServePlugin 启动插件服务器
// 这个函数应该在插件的main函数中调用
// 它会自动注册所有需要的 gob 类型，插件开发者不需要手动注册
func ServePlugin(impl ToolPluginInterface) {
	// 自动注册所有需要的 gob 类型，用于 RPC 通信
	registerGobTypes()

	// 创建插件映射，将我们的实现注册到插件系统
	pluginMap := map[string]plugin.Plugin{
		"tool": &ToolPlugin{Impl: impl},
	}

	// 启动插件服务器
	// 这会让插件进入监听状态，等待主程序的连接
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig, // 使用共享的握手配置
		Plugins:         pluginMap,       // 注册插件映射
	})
}
