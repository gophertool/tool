# GopherTool

一个功能丰富的Go语言工具库，提供插件系统、缓存管理、图像处理和日志记录等核心功能，专为被其他项目引用而设计。

## 项目简介

GopherTool是一个现代化的Go工具集合库，提供了多种高性能、易用的工具模块。该库的设计理念是：

1. **模块化设计** - 每个工具模块可以独立引用和使用
2. **插件化架构** - 基于hashicorp/go-plugin的强大插件系统
3. **统一接口** - 提供一致的API设计和使用体验
4. **高性能** - 专注于提供高质量、高性能的工具实现
5. **易扩展** - 支持自定义插件和驱动扩展

## 项目结构

```
├── .github/              # GitHub工作流和CI/CD配置
│   └── workflows/        # 自动化构建和发布流程
├── db/                   # 数据库相关工具
│   └── cache/            # 统一缓存接口和多驱动实现
│       ├── badgerdb/     # BadgerDB本地缓存实现
│       ├── buntdb/       # BuntDB内存缓存实现
│       ├── redis/        # Redis分布式缓存实现
│       ├── interface/    # 统一缓存接口定义
│       ├── config/       # 缓存配置管理
│       └── example/      # 缓存使用示例
├── image/                # 图像处理工具
│   ├── example/          # 图像处理示例
│   └── image.go          # 图像加载、保存和格式转换
├── log/                  # 高级日志工具
│   ├── color.go          # 彩色输出支持
│   └── log.go            # 多级别日志记录
├── plugin/               # 插件系统核心
│   ├── example/          # 插件开发和使用示例
│   │   ├── Makefile      # 插件构建脚本
│   │   └── plugin/       # 示例插件实现
│   ├── plugin.go         # 插件管理器和核心功能
│   ├── result.go         # 插件调用结果类型
│   └── tool.go           # 工具定义和选项
├── go.mod                # Go模块文件
├── go.sum                # Go模块依赖文件
├── tool.go               # 版本信息和主入口
├── LICENSE               # MIT许可证文件
└── README.md             # 项目说明文档
```

## 核心功能

### 🔌 插件系统

基于hashicorp/go-plugin的强大插件架构，支持动态加载和管理插件：

- **插件管理器** - 统一的插件生命周期管理
- **工具调用** - 支持结构化参数和类型安全的工具调用
- **插件扫描** - 自动发现和加载.tool.plugin文件
- **RPC通信** - 进程间安全通信和错误处理

### 💾 缓存系统

统一的缓存接口，支持多种后端驱动：

- **Redis** - 分布式内存缓存，支持集群和持久化
- **BadgerDB** - 高性能本地LSM树存储
- **BuntDB** - 快速内存数据库，支持持久化
- **统一接口** - 一致的API，轻松切换不同缓存后端
- **事务支持** - 原子性操作和事务管理

### 🖼️ 图像处理

完整的图像处理工具集：

- **多源加载** - 支持文件、URL、Base64、字节数组等多种来源
- **格式支持** - JPEG、PNG等常见格式的读取和保存
- **接口设计** - 灵活的Loader接口，易于扩展

### 📝 日志系统

高级日志记录功能：

- **多级别日志** - DEBUG、INFO、WARN、ERROR等级别
- **彩色输出** - 支持终端彩色显示
- **调用者追踪** - 自动显示日志调用位置
- **灵活配置** - 可自定义输出格式和过滤规则

## 安装使用

### 安装

```bash
go get github.com/gophertool/tool
```

### 使用缓存系统

```go
import (
    "time"
    "github.com/gophertool/tool/db/cache/interface"
    "github.com/gophertool/tool/db/cache/config"
)

func main() {
    // 创建Redis缓存配置
    cfg := config.Cache{
        Driver:   "redis",
        Host:     "localhost",
        Port:     "6379",
        Password: "",
        DB:       0,
    }
    
    // 创建缓存实例
    cache, err := _interface.New(cfg)
    if err != nil {
        panic(err)
    }
    defer cache.Close()
    
    // 基本操作
    err = cache.Set("key", "value", 5*time.Minute)
    if err != nil {
        panic(err)
    }
    
    value, err := cache.Get("key")
    if err != nil {
        panic(err)
    }
    fmt.Println(value) // 输出: value
    
    // 哈希操作
    err = cache.HSet("user:1", "name", "张三", time.Hour)
    if err != nil {
        panic(err)
    }
    
    // 队列操作
    err = cache.LPush("queue", "task1")
    if err != nil {
        panic(err)
    }
}
```

### 使用插件系统

```go
import "github.com/gophertool/tool/plugin"

func main() {
    // 创建插件管理器
    pm := plugin.NewPluginManager()
    defer pm.Shutdown()
    
    // 加载所有插件
    err := pm.LoadAllPlugins("./plugins")
    if err != nil {
        panic(err)
    }
    
    // 列出可用工具
    tools := pm.ListTools()
    for _, tool := range tools {
        fmt.Printf("工具: %s - %s\n", tool.Name, tool.Description)
    }
    
    // 调用工具
    result, err := pm.CallTool("current_time", map[string]interface{}{
        "format":   "2006-01-02 15:04:05",
        "timezone": "Asia/Shanghai",
    })
    if err != nil {
        panic(err)
    }
    
    // 处理结果
    for _, content := range result.Content {
        if textContent, ok := content.(plugin.TextContent); ok {
            fmt.Println(textContent.Text)
        }
    }
}
```

### 使用图像处理

```go
import "github.com/gophertool/tool/image"

func main() {
    // 创建图像加载器
    loader := image.NewLoader()
    
    // 从文件加载图像
    img, err := loader.LoadFromFile("test.jpg")
    if err != nil {
        panic(err)
    }
    
    // 保存为PNG格式
    err = image.SaveImage(img, "output.png", "png")
    if err != nil {
        panic(err)
    }
    
    // 从URL加载图像
    img2, err := loader.LoadFromURL("https://example.com/image.jpg")
    if err != nil {
        panic(err)
    }
    
    // 从Base64加载图像
    img3, err := loader.LoadFromBase64("data:image/jpeg;base64,/9j/4AAQ...")
    if err != nil {
        panic(err)
    }
}
```

### 使用日志系统

```go
import "github.com/gophertool/tool/log"

func main() {
    // 基本日志记录
    log.Info("这是一条信息日志")
    log.Warn("这是一条警告日志")
    log.Error("这是一条错误日志")
    log.Debug("这是一条调试日志")
    
    // 使用级别记录
    log.Println(log.INFO, "使用级别的信息日志")
    log.Printf(log.ERROR, "格式化错误日志: %s", "错误信息")
    
    // 设置调用者层级（用于显示正确的调用位置）
    log.SetCallerLevel(3)
}
```

## 详细功能说明

### 插件系统 (plugin/)

**核心组件：**
- `PluginManager` - 插件生命周期管理器
- `Tool` - 工具定义和参数模式
- `CallToolResult` - 工具调用结果封装
- `LoadedPlugin` - 已加载插件的运行时信息

**主要功能：**
- 🔍 **插件扫描** - 递归扫描目录，自动发现.tool.plugin文件
- 🚀 **动态加载** - 运行时加载和卸载插件，支持热更新
- 🛠️ **工具调用** - 类型安全的工具调用，支持结构化参数
- 🔒 **进程隔离** - 基于RPC的进程间通信，确保主程序稳定性
- 📊 **状态管理** - 实时监控插件状态和健康检查

### 缓存系统 (db/cache/)

**统一接口设计：**
```go
type Cache interface {
    // 基本操作
    Get(key string) (string, error)
    Set(key string, value string, ttl time.Duration) error
    Delete(key string) error
    Exists(key string) (bool, error)
    Expire(key string, ttl time.Duration) error
    
    // 哈希操作
    HGet(key, field string) (string, error)
    HSet(key, field, value string, ttl time.Duration) error
    HDel(key, field string) error
    HGetAll(key string) (map[string]string, error)
    
    // 队列操作
    LPush(key string, value string) error
    RPush(key string, value string) error
    LPop(key string) (string, error)
    RPop(key string) (string, error)
    PopAll(key string) ([]string, error)
    Len(key string) (int64, error)
    
    // 事务操作
    BeginTx() (Tx, error)
}
```

**支持的驱动：**
- 🔴 **Redis** - 分布式缓存，支持集群、持久化、发布订阅
- 🟡 **BadgerDB** - 高性能LSM树存储，适合大数据量本地缓存
- 🟢 **BuntDB** - 内存数据库，支持事务和持久化

**特性：**
- 🔄 **驱动切换** - 通过配置轻松切换不同缓存后端
- 🏭 **工厂模式** - 统一的实例创建和管理
- 🔐 **事务支持** - 原子性操作，确保数据一致性
- ⚡ **高性能** - 优化的连接池和批量操作

### 图像处理 (image/)

**Loader接口：**
```go
type Loader interface {
    LoadFromFile(filePath string) (image.Image, error)
    LoadFromURL(url string) (image.Image, error)
    LoadFromBase64(base64Str string) (image.Image, error)
    LoadFromBytes(data []byte) (image.Image, error)
    LoadFromReader(reader io.Reader) (image.Image, error)
}
```

**功能特性：**
- 📁 **多源加载** - 文件、URL、Base64、字节数组、io.Reader
- 🖼️ **格式支持** - JPEG、PNG等主流图像格式
- 💾 **智能保存** - 自动格式检测和转换
- 🔧 **易扩展** - 接口化设计，便于添加新的加载方式
- 🛡️ **错误处理** - 完善的错误处理和类型检查

### 日志系统 (log/)

**日志级别：**
- `DEBUG` - 调试信息，开发阶段使用
- `INFO` - 一般信息，正常运行状态
- `WARN` - 警告信息，需要注意但不影响运行
- `ERROR` - 错误信息，影响功能但不致命
- `DATA` - 数据输出，纯数据记录

**高级功能：**
- 🎨 **彩色输出** - 终端彩色显示，提升可读性
- 📍 **调用者追踪** - 自动显示日志调用的文件和行号
- 🔧 **灵活配置** - 可自定义输出格式、过滤规则
- 🎯 **精确定位** - 智能跳过框架代码，显示真实调用位置
- 📝 **多种输出** - 支持标准输出、错误输出等多种目标

## 技术特性

### 🏗️ 架构设计
- **模块化** - 松耦合设计，各模块可独立使用
- **接口驱动** - 统一接口规范，易于扩展和测试
- **工厂模式** - 统一的实例创建和配置管理
- **插件化** - 支持动态扩展和第三方插件

### ⚡ 性能优化
- **连接池** - 数据库和缓存连接复用
- **批量操作** - 支持批量读写，提升吞吐量
- **内存管理** - 优化的内存使用和垃圾回收
- **并发安全** - 线程安全的设计和实现

### 🛡️ 可靠性
- **错误处理** - 完善的错误处理和恢复机制
- **类型安全** - 强类型检查，减少运行时错误
- **测试覆盖** - 完整的单元测试和集成测试
- **文档完整** - 详细的代码注释和使用文档

### 🔧 易用性
- **简单配置** - 最小化配置，开箱即用
- **统一API** - 一致的接口设计和命名规范
- **丰富示例** - 完整的使用示例和最佳实践
- **版本管理** - 语义化版本控制，向后兼容

## 开发指南

### 插件开发

创建自定义插件的步骤：

1. **实现插件接口**
```go
type ToolPlugin interface {
    GetInfo() PluginInfo
    GetTools() []Tool
    CallTool(toolName string, args map[string]any) (*CallToolResult, error)
}
```

2. **编译插件**
```bash
go build -o my-plugin.tool.plugin main.go
```

3. **部署插件**
将编译好的.tool.plugin文件放入插件目录即可自动加载。

### 缓存驱动扩展

添加新的缓存驱动：

1. **实现Cache接口**
```go
type MyCache struct{}

func (c *MyCache) Get(key string) (string, error) {
    // 实现获取逻辑
}
// ... 实现其他接口方法
```

2. **注册驱动**
```go
func init() {
    _interface.RegisterDriver("mycache", func(cfg config.Cache) (_interface.Cache, error) {
        return NewMyCache(cfg), nil
    })
}
```

### 测试

运行测试套件：
```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./db/cache/...
go test ./plugin/...
go test ./image/...
go test ./log/...

# 运行测试并显示覆盖率
go test -cover ./...
```

## 版本信息

当前版本：`v0.0.8-20250724`

### 版本历史
- **v0.0.8** - 完善插件系统，优化缓存接口
- **v0.0.7** - 添加图像处理模块
- **v0.0.6** - 重构日志系统，添加彩色输出
- **v0.0.5** - 实现多驱动缓存系统
- **v0.0.4** - 初始插件系统实现

## 依赖管理

主要依赖：
```go
require (
    github.com/dgraph-io/badger v1.6.2      // BadgerDB存储引擎
    github.com/go-redis/redis v6.15.9       // Redis客户端
    github.com/hashicorp/go-plugin v1.6.3   // 插件系统框架
    github.com/tidwall/buntdb v1.3.2        // BuntDB内存数据库
)
```

## 贡献指南

我们欢迎各种形式的贡献！

### 🐛 报告问题
- 使用GitHub Issues报告bug
- 提供详细的复现步骤和环境信息
- 包含相关的日志和错误信息

### 💡 功能建议
- 在Issues中描述新功能需求
- 说明使用场景和预期效果
- 讨论实现方案和API设计

### 🔧 代码贡献

1. **Fork项目**
```bash
git clone https://github.com/gophertool/tool.git
cd tool
```

2. **创建功能分支**
```bash
git checkout -b feature/amazing-feature
```

3. **开发和测试**
```bash
# 编写代码
# 添加测试
go test ./...

# 检查代码格式
go fmt ./...
go vet ./...
```

4. **提交更改**
```bash
git add .
git commit -m "feat: add amazing feature"
```

5. **推送和PR**
```bash
git push origin feature/amazing-feature
# 在GitHub上创建Pull Request
```

### 📝 代码规范
- 遵循Go官方代码风格
- 为所有公共函数添加注释
- 编写单元测试，保持测试覆盖率
- 使用语义化提交信息

### 🏷️ 提交信息格式
```
type(scope): description

[optional body]

[optional footer]
```

类型：
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

## 社区

- 💬 **讨论**: GitHub Discussions
- 🐛 **问题**: GitHub Issues
- 📖 **文档**: [项目Wiki](https://github.com/gophertool/tool/wiki)

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

```
MIT License

Copyright (c) 2024 GopherTool

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

<div align="center">

**⭐ 如果这个项目对您有帮助，请给我们一个Star！**

[🏠 主页](https://github.com/gophertool/tool) • [📖 文档](https://github.com/gophertool/tool/wiki) • [🐛 问题](https://github.com/gophertool/tool/issues) • [💬 讨论](https://github.com/gophertool/tool/discussions)

</div>
