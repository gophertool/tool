# GopherTool

一个聚合多种工具的Go语言工具库，专为被其他项目引用而设计。

## 项目简介

GopherTool是一个工具集合库，提供了多种常用工具函数和实用程序，可以被其他Go项目引用使用。该库的设计理念是：

1. 提供模块化的工具集，每个工具可以单独引用
2. 同时提供统一的接口，方便整体引入
3. 专注于提供高质量、高性能的工具函数

## 项目结构

```
├── cmd/                  # 示例和测试工具的命令行入口
├── db/                   # 数据库相关工具
│   ├── cache/            # 缓存工具
│   │   ├── badgerdb/     # BadgerDB缓存实现
│   │   ├── buntdb/       # BuntDB缓存实现
│   │   ├── redis/        # Redis缓存实现
│   │   ├── interface/    # 缓存接口定义
│   │   ├── config/       # 缓存配置
│   │   └── example/      # 缓存使用示例
│   └── sql/              # SQL数据库工具
├── docs/                 # 文档
├── image/                # 图像处理相关工具
├── internal/             # 内部使用的代码，不对外暴露
├── pkg/                  # 可被外部引用的公共包
│   ├── convert/          # 类型转换工具
│   ├── crypto/           # 加密解密工具
│   ├── file/             # 文件操作工具
│   ├── http/             # HTTP相关工具
│   ├── json/             # JSON处理工具
│   ├── log/              # 日志工具
│   ├── math/             # 数学计算工具
│   ├── net/              # 网络工具
│   ├── slice/            # 切片操作工具
│   ├── str/              # 字符串处理工具
│   ├── time/             # 时间处理工具
│   └── util/             # 通用工具
├── scripts/              # 构建和维护脚本
├── test/                 # 测试文件
├── go.mod                # Go模块文件
├── go.sum                # Go模块依赖文件
├── LICENSE               # 许可证文件
└── README.md             # 项目说明文档
```

## 使用方法

### 安装

```bash
go get github.com/gophertool/tool
```

### 单独引用某个工具

```go
import "github.com/gophertool/tool/pkg/str"

func main() {
    result := str.Reverse("hello")
    fmt.Println(result) // 输出: olleh
}
```

### 使用统一接口

```go
import "github.com/gophertool/tool"

func main() {
    result := tool.Str.Reverse("hello")
    fmt.Println(result) // 输出: olleh
}
```

### 使用缓存工具

```go
import (
    "github.com/gophertool/tool/db/cache/redis"
    "github.com/gophertool/tool/db/cache/config"
)

func main() {
    // 创建Redis缓存配置
    cfg := &config.RedisConfig{
        Host:     "localhost",
        Port:     6379,
        Password: "",
        DB:       0,
    }
  
    // 初始化Redis缓存
    cache := redis.NewRedisCache(cfg)
  
    // 设置缓存
    err := cache.Set("key", "value", 300) // 缓存5分钟
    if err != nil {
        panic(err)
    }
  
    // 获取缓存
    value, err := cache.Get("key")
    if err != nil {
        panic(err)
    }
    fmt.Println(value) // 输出: value
}
```

## 工具列表

### 数据库相关工具

- **db/cache**: 缓存工具集合
  - **badgerdb**: 基于BadgerDB的本地缓存实现
  - **buntdb**: 基于BuntDB的内存缓存实现
  - **redis**: 基于Redis的分布式缓存实现
  - **interface**: 缓存接口定义，支持多种缓存后端
  - **config**: 缓存配置管理工具
  - **example**: 缓存使用示例和最佳实践

### 通用工具

- **convert**: 提供各种类型转换工具
- **crypto**: 提供加密解密、哈希等功能
- **file**: 文件和目录操作工具
- **http**: HTTP请求和响应处理工具
- **image**: 图像处理相关工具
- **json**: JSON序列化和反序列化工具
- **log**: 日志记录工具
- **math**: 数学计算和统计工具
- **net**: 网络相关工具
- **slice**: 切片操作工具
- **str**: 字符串处理工具
- **time**: 时间处理和格式化工具
- **util**: 通用工具函数

## 贡献指南

欢迎贡献代码或提出建议！请遵循以下步骤：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件
