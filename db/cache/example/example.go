// cache包使用示例
// 展示如何使用统一的工厂函数创建不同类型的缓存实例
//
// 本文件提供了完整的缓存使用示例，包括Redis、BadgerDB和BuntDB三种实现
// 通过实际代码演示统一接口的优势和使用方法
//
// 示例内容：
// - 基本键值操作：Set/Get/Delete/Exists/Expire
// - 队列操作：Push/Pop/LPush/RPush/LPop/RPop/PopAll/Len
// - 哈希表操作：HSet/HGet/HDel/HGetAll
// - 事务操作：BeginTx/Commit/Rollback
// - 驱动发现：列出所有可用驱动
//
// 演示特点：
// - 统一的API调用方式
// - 相同的业务逻辑代码
// - 灵活的驱动切换
// - 完整的错误处理
//
// 运行方式：
//
//	go run example.go
//	或在测试中调用各个Example函数
//
// 作者: gophertool
package example

import (
	"fmt"
	"log"
	"time"

	"github.com/gophertool/tool/db/cache/config"
	_interface "github.com/gophertool/tool/db/cache/interface"

	// 导入所有实现以确保驱动注册
	_ "github.com/gophertool/tool/db/cache/badgerdb"
	_ "github.com/gophertool/tool/db/cache/buntdb"
	_ "github.com/gophertool/tool/db/cache/redis"
)

// ExampleRedisUsage Redis缓存使用示例
func ExampleRedisUsage() {
	// Redis配置
	cfg := config.Cache{
		Driver:   config.CacheDriverRedis,
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	// 使用统一工厂函数创建Redis缓存实例
	cache, err := _interface.New(cfg)
	if err != nil {
		log.Printf("创建Redis缓存失败: %v", err)
		return
	}
	defer cache.Close()

	// 基本键值操作示例
	demonstrateBasicOperations(cache, "Redis")

	// 队列操作示例
	demonstrateQueueOperations(cache, "Redis")

	// 哈希表操作示例
	demonstrateHashOperations(cache, "Redis")
}

// ExampleBadgerUsage BadgerDB缓存使用示例
func ExampleBadgerUsage() {
	// BadgerDB配置
	cfg := config.Cache{
		Driver: config.CacheDriverBadger,
		Path:   "./badger_data",
	}

	// 使用统一工厂函数创建BadgerDB缓存实例
	cache, err := _interface.New(cfg)
	if err != nil {
		log.Printf("创建BadgerDB缓存失败: %v", err)
		return
	}
	defer cache.Close()

	// 基本键值操作示例
	demonstrateBasicOperations(cache, "BadgerDB")

	// 队列操作示例
	demonstrateQueueOperations(cache, "BadgerDB")

	// 哈希表操作示例
	demonstrateHashOperations(cache, "BadgerDB")
}

// ExampleBuntUsage BuntDB缓存使用示例
func ExampleBuntUsage() {
	// BuntDB配置
	cfg := config.Cache{
		Driver: config.CacheDriverBuntdb,
		Path:   "./bunt_data.db",
	}

	// 使用统一工厂函数创建BuntDB缓存实例
	cache, err := _interface.New(cfg)
	if err != nil {
		log.Printf("创建BuntDB缓存失败: %v", err)
		return
	}
	defer cache.Close()

	// 基本键值操作示例
	demonstrateBasicOperations(cache, "BuntDB")

	// 队列操作示例
	demonstrateQueueOperations(cache, "BuntDB")

	// 哈希表操作示例
	demonstrateHashOperations(cache, "BuntDB")
}

// demonstrateBasicOperations 演示基本的键值操作
func demonstrateBasicOperations(cache _interface.Cache, driverName string) {
	fmt.Printf("\n=== %s 基本操作示例 ===\n", driverName)

	// 设置键值
	err := cache.Set("user:1001", "张三", 5*time.Minute)
	if err != nil {
		fmt.Printf("设置键值失败: %v\n", err)
		return
	}
	fmt.Println("✓ 设置键值: user:1001 = 张三")

	// 获取键值
	value, err := cache.Get("user:1001")
	if err != nil {
		fmt.Printf("获取键值失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 获取键值: user:1001 = %s\n", value)

	// 检查键是否存在
	exists, err := cache.Exists("user:1001")
	if err != nil {
		fmt.Printf("检查键存在失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 键是否存在: user:1001 = %t\n", exists)

	// 设置过期时间
	err = cache.Expire("user:1001", 10*time.Second)
	if err != nil {
		fmt.Printf("设置过期时间失败: %v\n", err)
		return
	}
	fmt.Println("✓ 设置过期时间: user:1001 = 10秒")

	// 删除键
	err = cache.Delete("user:1001")
	if err != nil {
		fmt.Printf("删除键失败: %v\n", err)
		return
	}
	fmt.Println("✓ 删除键: user:1001")
}

// demonstrateQueueOperations 演示队列操作
func demonstrateQueueOperations(cache _interface.Cache, driverName string) {
	fmt.Printf("\n=== %s 队列操作示例 ===\n", driverName)

	queueKey := "task_queue"

	// 向队列推入元素
	tasks := []string{"任务1", "任务2", "任务3", "任务4"}
	for _, task := range tasks {
		err := cache.RPush(queueKey, task)
		if err != nil {
			fmt.Printf("推入队列失败: %v\n", err)
			return
		}
		fmt.Printf("✓ 推入队列: %s\n", task)
	}

	// 获取队列长度
	length, err := cache.Len(queueKey)
	if err != nil {
		fmt.Printf("获取队列长度失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 队列长度: %d\n", length)

	// 从队列弹出元素
	task, err := cache.LPop(queueKey)
	if err != nil {
		fmt.Printf("弹出队列失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 弹出队列: %s\n", task)

	// 弹出所有剩余元素
	remainingTasks, err := cache.PopAll(queueKey)
	if err != nil {
		fmt.Printf("弹出所有元素失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 弹出所有剩余元素: %v\n", remainingTasks)
}

// demonstrateHashOperations 演示哈希表操作
func demonstrateHashOperations(cache _interface.Cache, driverName string) {
	fmt.Printf("\n=== %s 哈希表操作示例 ===\n", driverName)

	hashKey := "user:profile:1001"

	// 设置哈希字段
	fields := map[string]string{
		"name":  "李四",
		"age":   "28",
		"city":  "北京",
		"email": "lisi@example.com",
	}

	for field, value := range fields {
		err := cache.HSet(hashKey, field, value, 10*time.Minute)
		if err != nil {
			fmt.Printf("设置哈希字段失败: %v\n", err)
			return
		}
		fmt.Printf("✓ 设置哈希字段: %s = %s\n", field, value)
	}

	// 获取单个哈希字段
	name, err := cache.HGet(hashKey, "name")
	if err != nil {
		fmt.Printf("获取哈希字段失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 获取哈希字段: name = %s\n", name)

	// 获取所有哈希字段
	allFields, err := cache.HGetAll(hashKey)
	if err != nil {
		fmt.Printf("获取所有哈希字段失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 获取所有哈希字段: %v\n", allFields)

	// 删除哈希字段
	err = cache.HDel(hashKey, "email")
	if err != nil {
		fmt.Printf("删除哈希字段失败: %v\n", err)
		return
	}
	fmt.Println("✓ 删除哈希字段: email")
}

// ExampleTransactionUsage 演示事务操作
func ExampleTransactionUsage() {
	fmt.Println("\n=== 事务操作示例 ===")

	// 使用BadgerDB进行事务示例
	cfg := config.Cache{
		Driver: config.CacheDriverBadger,
		Path:   "./badger_tx_data",
	}

	cache, err := _interface.New(cfg)
	if err != nil {
		log.Printf("创建缓存失败: %v", err)
		return
	}
	defer cache.Close()

	// 开始事务
	tx, err := cache.BeginTx()
	if err != nil {
		fmt.Printf("开始事务失败: %v\n", err)
		return
	}

	// 在事务中执行操作
	err = tx.Set("tx_key1", "事务值1", 0)
	if err != nil {
		fmt.Printf("事务设置失败: %v\n", err)
		tx.Rollback()
		return
	}

	err = tx.Set("tx_key2", "事务值2", 0)
	if err != nil {
		fmt.Printf("事务设置失败: %v\n", err)
		tx.Rollback()
		return
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		fmt.Printf("提交事务失败: %v\n", err)
		return
	}

	fmt.Println("✓ 事务操作成功完成")
}

// ListAvailableDrivers 列出所有可用的缓存驱动
func ListAvailableDrivers() {
	fmt.Println("\n=== 可用的缓存驱动 ===")
	drivers := _interface.GetRegisteredDrivers()
	for i, driver := range drivers {
		fmt.Printf("%d. %s\n", i+1, driver)
	}
}
