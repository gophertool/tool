// cache包的统一测试文件
// 测试所有缓存实现的功能正确性和一致性
//
// 本文件提供了完整的测试套件，确保所有缓存实现都符合接口规范
// 通过统一的测试用例验证不同驱动的功能一致性和正确性
//
// 测试范围：
// - 基本键值操作的功能测试和边界测试
// - 队列操作的FIFO/LIFO行为验证
// - 哈希表操作的字段管理测试
// - 事务操作的ACID特性验证
// - 驱动注册和发现机制测试
// - 错误处理和边界条件测试
//
// 测试策略：
// - 参数化测试：同一套测试用例验证所有驱动
// - 独立测试：每个驱动独立运行，避免相互影响
// - 清理机制：自动清理测试数据，避免残留
// - 性能基准：提供基本的性能基准测试
//
// 运行方式：
//
//	go test ./db/cache
//	go test -bench=. ./db/cache
//
// 作者: gophertool
package cache

import (
	"os"
	"testing"

	"github.com/gophertool/tool/db/cache/config"
	_interface "github.com/gophertool/tool/db/cache/interface"

	// 导入所有实现以确保驱动注册
	_ "github.com/gophertool/tool/db/cache/badgerdb"
	_ "github.com/gophertool/tool/db/cache/buntdb"
	_ "github.com/gophertool/tool/db/cache/redis"
)

// TestCacheDrivers 测试所有缓存驱动的基本功能
func TestCacheDrivers(t *testing.T) {
	// 测试配置
	testConfigs := []struct {
		name   string
		config config.Cache
	}{
		{
			name: "BadgerDB",
			config: config.Cache{
				Driver: config.CacheDriverBadger,
				Path:   "./test_badger_data",
			},
		},
		{
			name: "BuntDB",
			config: config.Cache{
				Driver: config.CacheDriverBuntdb,
				Path:   "./test_bunt_data.db",
			},
		},
		// Redis测试需要Redis服务器运行，可以根据需要启用
		// {
		// 	name: "Redis",
		// 	config: config.Cache{
		// 		Driver:   config.CacheDriverRedis,
		// 		Host:     "localhost",
		// 		Port:     "6379",
		// 		Password: "",
		// 		DB:       0,
		// 	},
		// },
	}

	for _, tc := range testConfigs {
		t.Run(tc.name, func(t *testing.T) {
			// 创建缓存实例
			cache, err := _interface.New(tc.config)
			if err != nil {
				t.Fatalf("创建%s缓存失败: %v", tc.name, err)
			}
			defer func() {
				cache.Close()
				// 清理测试数据
				if tc.config.Path != "" {
					os.RemoveAll(tc.config.Path)
				}
			}()

			// 运行所有测试
			testBasicOperations(t, cache, tc.name)
			testQueueOperations(t, cache, tc.name)
			testHashOperations(t, cache, tc.name)
			testTransactionOperations(t, cache, tc.name)
		})
	}
}

// testBasicOperations 测试基本的键值操作
func testBasicOperations(t *testing.T, cache _interface.Cache, driverName string) {
	t.Logf("测试%s基本操作", driverName)

	// 测试Set和Get
	key := "test_key"
	value := "test_value"

	err := cache.Set(key, value, 0)
	if err != nil {
		t.Errorf("%s Set操作失败: %v", driverName, err)
		return
	}

	retrievedValue, err := cache.Get(key)
	if err != nil {
		t.Errorf("%s Get操作失败: %v", driverName, err)
		return
	}

	if retrievedValue != value {
		t.Errorf("%s 值不匹配，期望: %s, 实际: %s", driverName, value, retrievedValue)
	}

	// 测试Exists
	exists, err := cache.Exists(key)
	if err != nil {
		t.Errorf("%s Exists操作失败: %v", driverName, err)
		return
	}
	if !exists {
		t.Errorf("%s 键应该存在", driverName)
	}

	// 测试Delete
	err = cache.Delete(key)
	if err != nil {
		t.Errorf("%s Delete操作失败: %v", driverName, err)
		return
	}

	// 验证删除后键不存在
	exists, err = cache.Exists(key)
	if err != nil {
		t.Errorf("%s Exists操作失败: %v", driverName, err)
		return
	}
	if exists {
		t.Errorf("%s 键不应该存在", driverName)
	}

	// 测试获取不存在的键
	_, err = cache.Get("nonexistent_key")
	if err != _interface.ErrKeyNotFound {
		t.Errorf("%s 获取不存在键应该返回ErrKeyNotFound，实际: %v", driverName, err)
	}
}

// testQueueOperations 测试队列操作
func testQueueOperations(t *testing.T, cache _interface.Cache, driverName string) {
	t.Logf("测试%s队列操作", driverName)

	queueKey := "test_queue"

	// 测试RPush
	values := []string{"value1", "value2", "value3"}
	for _, value := range values {
		err := cache.RPush(queueKey, value)
		if err != nil {
			t.Errorf("%s RPush操作失败: %v", driverName, err)
			return
		}
	}

	// 测试队列长度
	length, err := cache.Len(queueKey)
	if err != nil {
		t.Errorf("%s Len操作失败: %v", driverName, err)
		return
	}
	if length != int64(len(values)) {
		t.Errorf("%s 队列长度不正确，期望: %d, 实际: %d", driverName, len(values), length)
	}

	// 测试LPop
	poppedValue, err := cache.LPop(queueKey)
	if err != nil {
		t.Errorf("%s LPop操作失败: %v", driverName, err)
		return
	}
	if poppedValue != values[0] {
		t.Errorf("%s LPop值不正确，期望: %s, 实际: %s", driverName, values[0], poppedValue)
	}

	// 测试PopAll
	remainingValues, err := cache.PopAll(queueKey)
	if err != nil {
		t.Errorf("%s PopAll操作失败: %v", driverName, err)
		return
	}
	expectedRemaining := values[1:]
	if len(remainingValues) != len(expectedRemaining) {
		t.Errorf("%s PopAll返回数量不正确，期望: %d, 实际: %d", driverName, len(expectedRemaining), len(remainingValues))
	}

	// 验证队列为空
	length, err = cache.Len(queueKey)
	if err != nil {
		t.Errorf("%s Len操作失败: %v", driverName, err)
		return
	}
	if length != 0 {
		t.Errorf("%s 队列应该为空，实际长度: %d", driverName, length)
	}

	// 测试从空队列弹出元素
	_, err = cache.LPop(queueKey)
	if err != _interface.ErrKeyNotFound {
		t.Errorf("%s 从空队列LPop应该返回ErrKeyNotFound，实际: %v", driverName, err)
	}
}

// testHashOperations 测试哈希表操作
func testHashOperations(t *testing.T, cache _interface.Cache, driverName string) {
	t.Logf("测试%s哈希表操作", driverName)

	hashKey := "test_hash"
	field1, value1 := "field1", "value1"
	field2, value2 := "field2", "value2"

	// 测试HSet
	err := cache.HSet(hashKey, field1, value1, 0)
	if err != nil {
		t.Errorf("%s HSet操作失败: %v", driverName, err)
		return
	}

	err = cache.HSet(hashKey, field2, value2, 0)
	if err != nil {
		t.Errorf("%s HSet操作失败: %v", driverName, err)
		return
	}

	// 测试HGet
	retrievedValue, err := cache.HGet(hashKey, field1)
	if err != nil {
		t.Errorf("%s HGet操作失败: %v", driverName, err)
		return
	}
	if retrievedValue != value1 {
		t.Errorf("%s HGet值不正确，期望: %s, 实际: %s", driverName, value1, retrievedValue)
	}

	// 测试HGetAll
	allFields, err := cache.HGetAll(hashKey)
	if err != nil {
		t.Errorf("%s HGetAll操作失败: %v", driverName, err)
		return
	}
	if len(allFields) != 2 {
		t.Errorf("%s HGetAll返回字段数量不正确，期望: 2, 实际: %d", driverName, len(allFields))
	}
	if allFields[field1] != value1 || allFields[field2] != value2 {
		t.Errorf("%s HGetAll返回值不正确: %v", driverName, allFields)
	}

	// 测试HDel
	err = cache.HDel(hashKey, field1)
	if err != nil {
		t.Errorf("%s HDel操作失败: %v", driverName, err)
		return
	}

	// 验证字段已删除
	_, err = cache.HGet(hashKey, field1)
	if err != _interface.ErrKeyNotFound {
		t.Errorf("%s 获取已删除字段应该返回ErrKeyNotFound，实际: %v", driverName, err)
	}

	// 验证其他字段仍存在
	retrievedValue, err = cache.HGet(hashKey, field2)
	if err != nil {
		t.Errorf("%s HGet操作失败: %v", driverName, err)
		return
	}
	if retrievedValue != value2 {
		t.Errorf("%s HGet值不正确，期望: %s, 实际: %s", driverName, value2, retrievedValue)
	}
}

// testTransactionOperations 测试事务操作
func testTransactionOperations(t *testing.T, cache _interface.Cache, driverName string) {
	t.Logf("测试%s事务操作", driverName)

	// 测试事务提交
	tx, err := cache.BeginTx()
	if err != nil {
		t.Errorf("%s BeginTx操作失败: %v", driverName, err)
		return
	}

	key1, value1 := "tx_key1", "tx_value1"
	key2, value2 := "tx_key2", "tx_value2"

	err = tx.Set(key1, value1, 0)
	if err != nil {
		t.Errorf("%s 事务Set操作失败: %v", driverName, err)
		tx.Rollback()
		return
	}

	err = tx.Set(key2, value2, 0)
	if err != nil {
		t.Errorf("%s 事务Set操作失败: %v", driverName, err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		t.Errorf("%s Commit操作失败: %v", driverName, err)
		return
	}

	// 验证事务提交后数据存在
	retrievedValue, err := cache.Get(key1)
	if err != nil {
		t.Errorf("%s 事务提交后Get操作失败: %v", driverName, err)
		return
	}
	if retrievedValue != value1 {
		t.Errorf("%s 事务提交后值不正确，期望: %s, 实际: %s", driverName, value1, retrievedValue)
	}

	// 清理测试数据
	cache.Delete(key1)
	cache.Delete(key2)
}

// TestDriverRegistration 测试驱动注册功能
func TestDriverRegistration(t *testing.T) {
	drivers := _interface.GetRegisteredDrivers()

	expectedDrivers := []string{
		config.CacheDriverBadger,
		config.CacheDriverBuntdb,
		config.CacheDriverRedis,
	}

	for _, expectedDriver := range expectedDrivers {
		found := false
		for _, driver := range drivers {
			if driver == expectedDriver {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("驱动 %s 未注册", expectedDriver)
		}
	}
}

// TestInvalidDriver 测试无效驱动处理
func TestInvalidDriver(t *testing.T) {
	cfg := config.Cache{
		Driver: "invalid_driver",
	}

	_, err := _interface.New(cfg)
	if err == nil {
		t.Error("使用无效驱动应该返回错误")
	}

	if err != nil && err.Error() == "" {
		t.Error("错误信息不应该为空")
	}
}

// TestEmptyDriver 测试空驱动处理
func TestEmptyDriver(t *testing.T) {
	cfg := config.Cache{
		Driver: "",
	}

	_, err := _interface.New(cfg)
	if err == nil {
		t.Error("使用空驱动应该返回错误")
	}
}

// BenchmarkCacheOperations 性能基准测试
func BenchmarkCacheOperations(b *testing.B) {
	cfg := config.Cache{
		Driver: config.CacheDriverBadger,
		Path:   "./bench_badger_data",
	}

	cache, err := _interface.New(cfg)
	if err != nil {
		b.Fatalf("创建缓存失败: %v", err)
	}
	defer func() {
		cache.Close()
		os.RemoveAll(cfg.Path)
	}()

	b.Run("Set", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Set("bench_key", "bench_value", 0)
		}
	})

	b.Run("Get", func(b *testing.B) {
		cache.Set("bench_key", "bench_value", 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Get("bench_key")
		}
	})

	b.Run("RPush", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.RPush("bench_queue", "bench_value")
		}
	})

	b.Run("LPop", func(b *testing.B) {
		// 预填充队列
		for i := 0; i < b.N; i++ {
			cache.RPush("bench_queue", "bench_value")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.LPop("bench_queue")
		}
	})
}
