// badgerdb包：基于BadgerDB的高性能本地缓存实现
// 提供键值存储、哈希表操作、队列操作和事务支持
//
// BadgerDB是一个高性能的LSM树Key-Value存储引擎，专为SSD优化设计
// 本包实现了Cache接口，提供统一的缓存操作API
//
// 主要特性：
// - 高性能读写操作，基于LSM树结构
// - 支持TTL过期机制
// - 队列操作（FIFO/LIFO）
// - 哈希表操作（通过复合键实现）
// - 事务支持（读写事务）
// - 线程安全的并发访问
// - 自动垃圾回收和压缩
// - 本地文件存储，无需外部依赖
//
// 使用场景：
// - 本地高性能缓存
// - 嵌入式应用存储
// - 单机应用的持久化缓存
//
// 作者: gophertool
package badgerdb

import (
	"bytes"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/gophertool/tool/db/cache/config"
	_interface "github.com/gophertool/tool/db/cache/interface"

	"github.com/dgraph-io/badger"
)

// 包初始化时注册BadgerDB驱动
func init() {
	_interface.RegisterDriver(config.CacheDriverBadger, NewBadgerStore)
}

// BadgerDb BadgerDB缓存实现结构体
type BadgerDb struct {
	db         *badger.DB // BadgerDB实例
	queueMutex sync.Map   // 用于队列操作的互斥锁映射
}

// LPush 将元素插入到列表头部
// 参数：
//
//	key - 列表键名
//	value - 要插入的值
//
// 返回值：
//
//	error - 操作错误
func (b *BadgerDb) LPush(key string, value string) error {
	b.lock(key)
	defer b.unlock(key)

	headKey := key + ":head"
	tailKey := key + ":tail"

	headVal, err := b.Get(headKey)
	if errors.Is(err, _interface.ErrKeyNotFound) {
		// 列表不存在，初始化
		if err := b.Set(headKey, "0", 0); err != nil {
			return err
		}
		if err := b.Set(tailKey, "1", 0); err != nil {
			return err
		}
		if err := b.Set(key+":0", value, 0); err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	// 解析头索引
	headIndex, err := strconv.ParseInt(headVal, 10, 64)
	if err != nil {
		return err
	}

	// 头索引减一
	headIndex--

	// 存储新元素
	if err := b.Set(key+":"+strconv.FormatInt(headIndex, 10), value, 0); err != nil {
		return err
	}

	// 更新头索引
	return b.Set(headKey, strconv.FormatInt(headIndex, 10), 0)
}

// RPush 将元素插入到列表尾部
// 参数：
//
//	key - 列表键名
//	value - 要插入的值
//
// 返回值：
//
//	error - 操作错误
func (b *BadgerDb) RPush(key string, value string) error {
	b.lock(key)
	defer b.unlock(key)

	headKey := key + ":head"
	tailKey := key + ":tail"

	tailVal, err := b.Get(tailKey)
	if errors.Is(err, _interface.ErrKeyNotFound) {
		// 列表不存在，初始化
		if err := b.Set(headKey, "0", 0); err != nil {
			return err
		}
		if err := b.Set(tailKey, "1", 0); err != nil {
			return err
		}
		if err := b.Set(key+":0", value, 0); err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	// 解析尾索引
	tailIndex, err := strconv.ParseInt(tailVal, 10, 64)
	if err != nil {
		return err
	}

	// 存储新元素
	if err := b.Set(key+":"+strconv.FormatInt(tailIndex, 10), value, 0); err != nil {
		return err
	}

	// 增加尾索引
	tailIndex++

	// 更新尾索引
	return b.Set(tailKey, strconv.FormatInt(tailIndex, 10), 0)
}

// LPop 弹出列表头部元素
// 参数：
//
//	key - 列表键名
//
// 返回值：
//
//	string - 弹出的元素值
//	error - 操作错误
func (b *BadgerDb) LPop(key string) (string, error) {
	b.lock(key)
	defer b.unlock(key)

	headKey := key + ":head"
	tailKey := key + ":tail"

	headVal, err := b.Get(headKey)
	if errors.Is(err, _interface.ErrKeyNotFound) {
		return "", _interface.ErrKeyNotFound
	} else if err != nil {
		return "", err
	}

	tailVal, err := b.Get(tailKey)
	if err != nil {
		return "", err
	}

	// 解析索引
	headIndex, err := strconv.ParseInt(headVal, 10, 64)
	if err != nil {
		return "", err
	}

	tailIndex, err := strconv.ParseInt(tailVal, 10, 64)
	if err != nil {
		return "", err
	}

	// 检查列表是否为空
	if headIndex >= tailIndex {
		return "", _interface.ErrKeyNotFound
	}

	// 获取头部元素，修复key格式问题
	elementKey := key + ":" + strconv.FormatInt(headIndex, 10)
	value, err := b.Get(elementKey)
	if err != nil {
		return "", err
	}

	// 删除元素
	if err := b.Delete(elementKey); err != nil {
		return "", err
	}

	// 增加头索引
	headIndex++

	// 更新头索引
	if err := b.Set(headKey, strconv.FormatInt(headIndex, 10), 0); err != nil {
		return "", err
	}

	return value, nil
}

// RPop 弹出列表尾部元素
func (b *BadgerDb) RPop(key string) (string, error) {
	b.lock(key)
	defer b.unlock(key)

	headKey := key + ":head"
	tailKey := key + ":tail"

	headVal, err := b.Get(headKey)
	if errors.Is(err, _interface.ErrKeyNotFound) {
		return "", _interface.ErrKeyNotFound
	} else if err != nil {
		return "", err
	}

	tailVal, err := b.Get(tailKey)
	if err != nil {
		return "", err
	}

	// 解析索引
	headIndex, err := strconv.ParseInt(headVal, 10, 64)
	if err != nil {
		return "", err
	}

	tailIndex, err := strconv.ParseInt(tailVal, 10, 64)
	if err != nil {
		return "", err
	}

	// 检查列表是否为空
	if headIndex >= tailIndex {
		return "", _interface.ErrKeyNotFound
	}

	// 减少尾索引
	tailIndex--

	// 获取尾部元素
	elementKey := key + ":" + strconv.FormatInt(tailIndex, 10)
	value, err := b.Get(elementKey)
	if err != nil {
		return "", err
	}

	// 删除元素
	if err := b.Delete(elementKey); err != nil {
		return "", err
	}

	// 更新尾索引
	if err := b.Set(tailKey, strconv.FormatInt(tailIndex, 10), 0); err != nil {
		return "", err
	}

	return value, nil
}

func (b *BadgerDb) lock(key string) {
	actual, _ := b.queueMutex.LoadOrStore(key, &sync.Mutex{})
	mutex := actual.(*sync.Mutex)
	mutex.Lock()
}

func (b *BadgerDb) unlock(key string) {
	if actual, ok := b.queueMutex.Load(key); ok {
		mutex := actual.(*sync.Mutex)
		mutex.Unlock()
	}
}

// Push 添加元素到列表尾部
func (b *BadgerDb) Push(key string, value string) error {
	return b.RPush(key, value)
}

// Pop 移除并返回列表最后一个元素
func (b *BadgerDb) Pop(key string) (string, error) {
	return b.LPop(key)
}

// PopAll 取出并清空整个列表
func (b *BadgerDb) PopAll(key string) ([]string, error) {
	b.lock(key)
	defer b.unlock(key)

	headKey := key + ":head"
	tailKey := key + ":tail"

	headVal, err := b.Get(headKey)
	if errors.Is(err, _interface.ErrKeyNotFound) {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}

	tailVal, err := b.Get(tailKey)
	if err != nil {
		return nil, err
	}

	// 解析索引
	headIndex, err := strconv.ParseInt(headVal, 10, 64)
	if err != nil {
		return nil, err
	}

	tailIndex, err := strconv.ParseInt(tailVal, 10, 64)
	if err != nil {
		return nil, err
	}

	// 检查列表是否为空
	if headIndex >= tailIndex {
		return []string{}, nil
	}

	// 获取所有元素
	result := make([]string, 0, tailIndex-headIndex)
	for i := headIndex; i < tailIndex; i++ {
		elementKey := key + ":" + strconv.FormatInt(i, 10)
		value, err := b.Get(elementKey)
		if err != nil {
			continue // 跳过获取失败的元素
		}
		result = append(result, value)

		// 删除元素
		_ = b.Delete(elementKey)
	}

	// 重置列表
	_ = b.Delete(headKey)
	_ = b.Delete(tailKey)

	return result, nil
}

// Len 获取列表长度
func (b *BadgerDb) Len(key string) (int64, error) {
	headKey := key + ":head"
	tailKey := key + ":tail"

	headVal, err := b.Get(headKey)
	if errors.Is(err, _interface.ErrKeyNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	tailVal, err := b.Get(tailKey)
	if err != nil {
		return 0, err
	}

	// 解析索引
	headIndex, err := strconv.ParseInt(headVal, 10, 64)
	if err != nil {
		return 0, err
	}

	tailIndex, err := strconv.ParseInt(tailVal, 10, 64)
	if err != nil {
		return 0, err
	}

	// 计算长度
	return tailIndex - headIndex, nil
}

func (b *BadgerDb) Close() {
	_ = b.db.Close()
}

// Get 获取指定key的值
// 参数：
//
//	key - 键名
//
// 返回值：
//
//	string - 键对应的值
//	error - 操作错误，键不存在时返回ErrKeyNotFound
func (b *BadgerDb) Get(key string) (string, error) {
	var val []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		return err
	})
	// 统一错误处理：将BadgerDB特定错误转换为接口标准错误
	if errors.Is(err, badger.ErrKeyNotFound) {
		return "", _interface.ErrKeyNotFound
	}
	return string(val), err
}

func (b *BadgerDb) Set(key string, value string, ttl time.Duration) error {
	return b.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(value))
		if ttl > 0 {
			e.WithTTL(ttl)
		}
		return txn.SetEntry(e)
	})
}

func (b *BadgerDb) Delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (b *BadgerDb) Exists(key string) (bool, error) {
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (b *BadgerDb) Expire(key string, ttl time.Duration) error {
	// 实现逻辑：先获取旧值，再重新设置 TTL
	return b.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(key), val).WithTTL(ttl)
		return txn.SetEntry(e)
	})
}

func (b *BadgerDb) HGet(key, field string) (string, error) {
	compositeKey := key + ":" + field
	return b.Get(compositeKey)
}
func (b *BadgerDb) HSet(key, field, value string, ttl time.Duration) error {
	compositeKey := key + ":" + field
	return b.Set(compositeKey, value, ttl)
}

func (b *BadgerDb) HDel(key, field string) error {
	compositeKey := key + ":" + field
	return b.Delete(compositeKey)
}

func (b *BadgerDb) HGetAll(key string) (map[string]string, error) {
	result := make(map[string]string)
	prefix := []byte(key + ":")

	err := b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			field := string(bytes.TrimPrefix(item.Key(), prefix))
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			result[field] = string(val)
		}
		return nil
	})

	return result, err
}

type badgerTx struct {
	txn *badger.Txn
}

func (tx *badgerTx) Set(key string, value string, ttl time.Duration) error {
	e := badger.NewEntry([]byte(key), []byte(value))
	if ttl > 0 {
		e.WithTTL(ttl)
	}
	return tx.txn.SetEntry(e)
}
func (tx *badgerTx) Delete(key string) error {
	return tx.txn.Delete([]byte(key))
}
func (tx *badgerTx) Commit() error {
	return tx.txn.Commit()
}

func (tx *badgerTx) Rollback() error {
	tx.txn.Discard()
	return nil
}

func (b *BadgerDb) BeginTx() (_interface.Tx, error) {
	return &badgerTx{txn: b.db.NewTransaction(true)}, nil // 读写事务
}

// NewBadgerStore 创建BadgerDB缓存实例的工厂函数
// 参数：
//
//	config - 缓存配置
//
// 返回值：
//
//	Cache - 缓存接口实例
//	error - 创建错误
func NewBadgerStore(config config.Cache) (_interface.Cache, error) {
	opts := badger.DefaultOptions(config.Path).
		WithLogger(nil).       // 禁用日志以提高性能
		WithSyncWrites(false). // 异步写入提高性能
		WithTruncate(true)     // 启动时清理损坏的数据
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerDb{db: db}, nil
}
