// interface包：定义统一的缓存接口和工厂函数
// 提供缓存操作的标准接口定义和驱动管理功能
//
// 本包是整个缓存系统的核心，定义了所有缓存实现必须遵循的接口规范
// 同时提供了统一的工厂函数和驱动注册机制
//
// 主要组件：
// - Cache接口：定义所有缓存操作的标准方法
// - Tx接口：定义事务操作的标准方法
// - 工厂函数：提供统一的缓存实例创建方法
// - 驱动注册：支持动态注册不同的缓存实现
// - 错误定义：统一的错误类型定义
//
// 支持的操作类型：
// - 基本键值操作（Get/Set/Delete/Exists/Expire）
// - 哈希表操作（HGet/HSet/HDel/HGetAll）
// - 队列操作（Push/Pop/LPush/RPush/LPop/RPop/PopAll/Len）
// - 事务操作（BeginTx/Commit/Rollback）
//
// 设计模式：
// - 工厂模式：统一创建不同类型的缓存实例
// - 注册模式：动态注册和发现缓存驱动
// - 接口隔离：分离基本操作和事务操作
//
// 作者: gophertool
package _interface

import (
	"errors"
	"fmt"
	"time"

	"github.com/gophertool/tool/db/cache/config"
)

// Cache 缓存接口
type Cache interface {
	Close()
	// Get 获取指定 key 的值
	Get(key string) (string, error)
	// Set 设置 key-value 并设置过期时间
	Set(key string, value string, ttl time.Duration) error
	// Delete 删除指定 key
	Delete(key string) error
	// Exists 判断 key 是否存在
	Exists(key string) (bool, error)
	// Expire 设置 key 的过期时间
	Expire(key string, ttl time.Duration) error

	// HGet 获取哈希表中指定 field 的值
	HGet(key, field string) (string, error)
	// HSet 设置哈希表中的 field-value，并设置过期时间
	HSet(key, field, value string, ttl time.Duration) error
	// HDel 删除哈希表中的一个或多个 field
	HDel(key, field string) error
	// HGetAll 获取哈希表中所有的 field 和 value
	HGetAll(key string) (map[string]string, error)

	// Push 向队列中推入元素（默认实现）
	Push(key string, value string) error
	// LPush 将元素插入到列表左边
	LPush(key string, value string) error
	// RPush 将元素插入到列表右边
	RPush(key string, value string) error
	// Pop 弹出队列中的元素（默认实现）
	Pop(key string) (string, error)
	// LPop 弹出列表最左边的元素
	LPop(key string) (string, error)
	// RPop 弹出列表最右边的元素
	RPop(key string) (string, error)
	// PopAll 弹出队列中所有元素
	PopAll(key string) ([]string, error)
	// Len 获取队列长度
	Len(key string) (int64, error)

	// BeginTx 开启事务操作
	BeginTx() (Tx, error) // 事务操作
}

// Tx 事务接口
type Tx interface {
	// Set 设置 key-value 并设置过期时间
	Set(key string, value string, ttl time.Duration) error
	// Delete 删除指定 key
	Delete(key string) error
	// Commit 提交事务
	Commit() error
	// Rollback 回滚事务
	Rollback() error
}

// NewStoreFunc 创建缓存实例的函数类型
type NewStoreFunc func(config config.Cache) (Cache, error)

var (
	// ErrKeyNotFound 键值不存在
	ErrKeyNotFound = errors.New("key not found")

	// ErrUnsupportedDriver 不支持的驱动类型
	ErrUnsupportedDriver = errors.New("unsupported cache driver")
)

// 存储不同驱动的构造函数
var storeFactories = make(map[string]NewStoreFunc)

// RegisterDriver 注册缓存驱动
func RegisterDriver(driverName string, newFunc NewStoreFunc) {
	storeFactories[driverName] = newFunc
}

// New 根据配置创建缓存实例的工厂函数
// 参数：
//
//	cfg - 缓存配置，包含驱动类型、连接信息等
//
// 返回值：
//
//	Cache - 缓存接口实例
//	error - 创建过程中的错误
func New(cfg config.Cache) (Cache, error) {
	if cfg.Driver == "" {
		return nil, fmt.Errorf("缓存驱动不能为空")
	}

	newFunc, exists := storeFactories[cfg.Driver]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedDriver, cfg.Driver)
	}

	cache, err := newFunc(cfg)
	if err != nil {
		return nil, fmt.Errorf("创建%s缓存实例失败: %w", cfg.Driver, err)
	}

	return cache, nil
}

// GetRegisteredDrivers 获取已注册的所有驱动名称
func GetRegisteredDrivers() []string {
	drivers := make([]string, 0, len(storeFactories))
	for driver := range storeFactories {
		drivers = append(drivers, driver)
	}
	return drivers
}
