// buntdb包：基于BuntDB的高性能内存缓存实现
// 提供键值存储、哈希表操作、队列操作和事务支持
//
// BuntDB是一个快速的内存数据库，适用于需要高性能读写的场景
// 本包实现了Cache接口，提供统一的缓存操作API
//
// 主要特性：
// - 纯内存存储，读写性能极佳
// - 支持持久化到文件
// - 支持TTL过期
// - 队列操作（FIFO/LIFO）
// - 哈希表操作
// - 事务支持
// - 线程安全
//
// 作者: gophertool
package buntdb

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/gophertool/tool/db/cache/config"
	_interface "github.com/gophertool/tool/db/cache/interface"

	"github.com/tidwall/buntdb"
)

// 包初始化时注册BuntDB驱动
func init() {
	_interface.RegisterDriver(config.CacheDriverBuntdb, NewBuntStore)
}

// BuntDb BuntDB缓存实现结构体
type BuntDb struct {
	db         *buntdb.DB // BuntDB实例
	queueMutex sync.Map   // 用于队列操作的互斥锁映射
}

// Close 关闭数据库连接
func (b *BuntDb) Close() {
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
func (b *BuntDb) Get(key string) (string, error) {
	var val string
	err := b.db.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(key)
		if err != nil {
			return err
		}
		val = v
		return nil
	})
	// 统一错误处理：将BuntDB特定错误转换为接口标准错误
	if errors.Is(err, buntdb.ErrNotFound) {
		return "", _interface.ErrKeyNotFound
	}
	return val, err
}

func (b *BuntDb) Set(key string, value string, ttl time.Duration) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		var opts *buntdb.SetOptions
		if ttl > 0 {
			opts = &buntdb.SetOptions{Expires: true, TTL: ttl}
		}
		_, _, err := tx.Set(key, value, opts)
		return err
	})
}

func (b *BuntDb) Delete(key string) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})
}

func (b *BuntDb) Exists(key string) (bool, error) {
	err := b.db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(key)
		return err
	})
	if errors.Is(err, buntdb.ErrNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (b *BuntDb) Expire(key string, ttl time.Duration) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(key, val, &buntdb.SetOptions{Expires: true, TTL: ttl})
		return err
	})
}

func (b *BuntDb) HGet(key, field string) (string, error) {
	compositeKey := key + ":" + field
	return b.Get(compositeKey)
}

func (b *BuntDb) HSet(key, field, value string, ttl time.Duration) error {
	compositeKey := key + ":" + field
	return b.Set(compositeKey, value, ttl)
}

func (b *BuntDb) HDel(key, field string) error {
	compositeKey := key + ":" + field
	return b.Delete(compositeKey)
}

func (b *BuntDb) HGetAll(key string) (map[string]string, error) {
	result := make(map[string]string)
	prefix := key + ":"

	err := b.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(prefix+"*", func(k, v string) bool {
			field := k[len(prefix):]
			result[field] = v
			return true
		})
	})

	return result, err
}

func (b *BuntDb) Push(key string, value string) error {
	return b.RPush(key, value)
}

func (b *BuntDb) LPush(key string, value string) error {
	b.lock(key)
	defer b.unlock(key)

	return b.db.Update(func(tx *buntdb.Tx) error {
		headKey := key + ":head"
		tailKey := key + ":tail"

		var head int64 = 0
		var tail int64 = 0

		if headVal, err := tx.Get(headKey); err == nil {
			if h, err := strconv.ParseInt(headVal, 10, 64); err == nil {
				head = h
			}
		}

		if tailVal, err := tx.Get(tailKey); err == nil {
			if t, err := strconv.ParseInt(tailVal, 10, 64); err == nil {
				tail = t
			}
		}

		head--

		elemKey := key + ":elem:" + strconv.FormatInt(head, 10)
		_, _, err := tx.Set(elemKey, value, nil)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(headKey, strconv.FormatInt(head, 10), nil)
		if err != nil {
			return err
		}

		if head == tail-1 {
			_, _, err = tx.Set(tailKey, strconv.FormatInt(tail, 10), nil)
		}

		return err
	})
}

func (b *BuntDb) RPush(key string, value string) error {
	b.lock(key)
	defer b.unlock(key)

	return b.db.Update(func(tx *buntdb.Tx) error {
		headKey := key + ":head"
		tailKey := key + ":tail"

		var head int64 = 0
		var tail int64 = 0

		if headVal, err := tx.Get(headKey); err == nil {
			if h, err := strconv.ParseInt(headVal, 10, 64); err == nil {
				head = h
			}
		}

		if tailVal, err := tx.Get(tailKey); err == nil {
			if t, err := strconv.ParseInt(tailVal, 10, 64); err == nil {
				tail = t
			}
		}

		elemKey := key + ":elem:" + strconv.FormatInt(tail, 10)
		_, _, err := tx.Set(elemKey, value, nil)
		if err != nil {
			return err
		}

		tail++

		_, _, err = tx.Set(tailKey, strconv.FormatInt(tail, 10), nil)
		if err != nil {
			return err
		}

		if tail-1 == head {
			_, _, err = tx.Set(headKey, strconv.FormatInt(head, 10), nil)
		}

		return err
	})
}
func (b *BuntDb) Pop(key string) (string, error) {
	return b.LPop(key)
}

// LPop 弹出列表头部元素
// 参数：
//
//	key - 列表键名
//
// 返回值：
//
//	string - 弹出的元素值
//	error - 操作错误，列表为空时返回ErrKeyNotFound
func (b *BuntDb) LPop(key string) (string, error) {
	b.lock(key)
	defer b.unlock(key)

	var result string

	err := b.db.Update(func(tx *buntdb.Tx) error {
		headKey := key + ":head"
		tailKey := key + ":tail"

		headVal, err := tx.Get(headKey)
		if err != nil {
			return err
		}

		tailVal, err := tx.Get(tailKey)
		if err != nil {
			return err
		}

		head, err := strconv.ParseInt(headVal, 10, 64)
		if err != nil {
			return err
		}

		tail, err := strconv.ParseInt(tailVal, 10, 64)
		if err != nil {
			return err
		}

		if head >= tail {
			return buntdb.ErrNotFound
		}

		elemKey := key + ":elem:" + strconv.FormatInt(head, 10)
		val, err := tx.Get(elemKey)
		if err != nil {
			return err
		}

		result = val

		_, err = tx.Delete(elemKey)
		if err != nil {
			return err
		}

		head++

		_, _, err = tx.Set(headKey, strconv.FormatInt(head, 10), nil)

		return err
	})

	// 统一错误处理：将BuntDB特定错误转换为接口标准错误
	if errors.Is(err, buntdb.ErrNotFound) {
		return "", _interface.ErrKeyNotFound
	}

	return result, err
}

// RPop 弹出列表尾部元素
// 参数：
//
//	key - 列表键名
//
// 返回值：
//
//	string - 弹出的元素值
//	error - 操作错误，列表为空时返回ErrKeyNotFound
func (b *BuntDb) RPop(key string) (string, error) {
	b.lock(key)
	defer b.unlock(key)

	var result string

	err := b.db.Update(func(tx *buntdb.Tx) error {
		headKey := key + ":head"
		tailKey := key + ":tail"

		headVal, err := tx.Get(headKey)
		if err != nil {
			return err
		}

		tailVal, err := tx.Get(tailKey)
		if err != nil {
			return err
		}

		head, err := strconv.ParseInt(headVal, 10, 64)
		if err != nil {
			return err
		}

		tail, err := strconv.ParseInt(tailVal, 10, 64)
		if err != nil {
			return err
		}

		if head >= tail {
			return buntdb.ErrNotFound
		}

		tail--

		elemKey := key + ":elem:" + strconv.FormatInt(tail, 10)
		val, err := tx.Get(elemKey)
		if err != nil {
			return err
		}

		result = val

		_, err = tx.Delete(elemKey)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(tailKey, strconv.FormatInt(tail, 10), nil)

		return err
	})

	// 统一错误处理：将BuntDB特定错误转换为接口标准错误
	if errors.Is(err, buntdb.ErrNotFound) {
		return "", _interface.ErrKeyNotFound
	}

	return result, err
}

func (b *BuntDb) PopAll(key string) ([]string, error) {
	b.lock(key)
	defer b.unlock(key)

	var result []string

	err := b.db.Update(func(tx *buntdb.Tx) error {
		headKey := key + ":head"
		tailKey := key + ":tail"

		headVal, err := tx.Get(headKey)
		if err != nil {
			return nil
		}

		tailVal, err := tx.Get(tailKey)
		if err != nil {
			return nil
		}

		head, err := strconv.ParseInt(headVal, 10, 64)
		if err != nil {
			return err
		}

		tail, err := strconv.ParseInt(tailVal, 10, 64)
		if err != nil {
			return err
		}

		// Check if queue is empty
		if head >= tail {
			return nil // Queue is empty, return empty result
		}

		// Get all elements from head to tail-1
		for i := head; i < tail; i++ {
			elemKey := key + ":elem:" + strconv.FormatInt(i, 10)
			val, err := tx.Get(elemKey)
			if err != nil {
				continue // Skip if element doesn't exist
			}

			result = append(result, val)

			// Delete the element
			_, err = tx.Delete(elemKey)
			if err != nil {
				return err
			}
		}

		// Reset indices (we'll use 0 for both head and tail to indicate empty queue)
		_, _, err = tx.Set(headKey, "0", nil)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(tailKey, "0", nil)
		return err
	})

	return result, err
}

func (b *BuntDb) Len(key string) (int64, error) {
	var length int64 = 0

	err := b.db.View(func(tx *buntdb.Tx) error {
		// Get current head and tail indices
		headKey := key + ":head"
		tailKey := key + ":tail"

		headVal, err := tx.Get(headKey)
		if err != nil {
			return nil // Queue doesn't exist, length is 0
		}

		tailVal, err := tx.Get(tailKey)
		if err != nil {
			return nil // Queue doesn't exist, length is 0
		}

		head, err := strconv.ParseInt(headVal, 10, 64)
		if err != nil {
			return err
		}

		tail, err := strconv.ParseInt(tailVal, 10, 64)
		if err != nil {
			return err
		}

		// Calculate length
		if tail > head {
			length = tail - head
		}

		return nil
	})

	return length, err
}

func (b *BuntDb) lock(key string) {
	actual, _ := b.queueMutex.LoadOrStore(key, &sync.Mutex{})
	mutex := actual.(*sync.Mutex)
	mutex.Lock()
}

func (b *BuntDb) unlock(key string) {
	if actual, ok := b.queueMutex.Load(key); ok {
		mutex := actual.(*sync.Mutex)
		mutex.Unlock()
	}
}

type buntTx struct {
	tx *buntdb.Tx
}

func (tx *buntTx) Set(key string, value string, ttl time.Duration) error {
	var opts *buntdb.SetOptions
	if ttl > 0 {
		opts = &buntdb.SetOptions{Expires: true, TTL: ttl}
	}
	_, _, err := tx.tx.Set(key, value, opts)
	return err
}

func (tx *buntTx) Delete(key string) error {
	_, err := tx.tx.Delete(key)
	return err
}

func (tx *buntTx) Commit() error {
	return tx.tx.Commit()
}

func (tx *buntTx) Rollback() error {
	return tx.tx.Rollback()
}

func (b *BuntDb) BeginTx() (_interface.Tx, error) {
	tx, err := b.db.Begin(true)
	if err != nil {
		return nil, err
	}
	return &buntTx{tx: tx}, nil
}

func NewBuntStore(config config.Cache) (_interface.Cache, error) {
	db, err := buntdb.Open(config.Path)
	if err != nil {
		return nil, err
	}
	return &BuntDb{db: db}, nil
}
