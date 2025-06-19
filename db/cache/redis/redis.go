// redis包：基于Redis的分布式缓存实现
// 提供键值存储、哈希表操作、队列操作和事务支持
//
// Redis是一个高性能的内存数据结构存储，支持多种数据结构
// 本包实现了Cache接口，提供统一的缓存操作API
//
// 主要特性：
// - 高性能的内存存储
// - 支持持久化
// - 原生队列操作支持
// - 丰富的数据结构
// - 事务支持（Pipeline）
// - 集群支持
// - 分布式缓存
//
// 作者: gophertool
package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/gophertool/tool/db/cache/config"
	_interface "github.com/gophertool/tool/db/cache/interface"
)

// 包初始化时注册Redis驱动
func init() {
	_interface.RegisterDriver(config.CacheDriverRedis, NewRedisClient)
}

// RedisDb Redis缓存实现结构体
type RedisDb struct {
	db *redis.Client // Redis客户端实例
}

// LPush 将元素插入到列表左边
// 参数：
//
//	key - 列表键名
//	value - 要插入的值
//
// 返回值：
//
//	error - 操作错误
func (r *RedisDb) LPush(key string, value string) error {
	return r.db.LPush(key, value).Err()
}

// RPush 将元素插入到列表右边
// 参数：
//
//	key - 列表键名
//	value - 要插入的值
//
// 返回值：
//
//	error - 操作错误
func (r *RedisDb) RPush(key string, value string) error {
	return r.db.RPush(key, value).Err()
}

// LPop 弹出列表最左边的元素
// 参数：
//
//	key - 列表键名
//
// 返回值：
//
//	string - 弹出的元素值
//	error - 操作错误，列表为空时返回ErrKeyNotFound
func (r *RedisDb) LPop(key string) (string, error) {
	val, err := r.db.LPop(key).Result()
	// 统一错误处理：将Redis特定错误转换为接口标准错误
	if errors.Is(err, redis.Nil) {
		return "", _interface.ErrKeyNotFound
	}
	return val, err
}

// RPop 弹出列表最右边的元素
// 参数：
//
//	key - 列表键名
//
// 返回值：
//
//	string - 弹出的元素值
//	error - 操作错误，列表为空时返回ErrKeyNotFound
func (r *RedisDb) RPop(key string) (string, error) {
	val, err := r.db.RPop(key).Result()
	// 统一错误处理：将Redis特定错误转换为接口标准错误
	if errors.Is(err, redis.Nil) {
		return "", _interface.ErrKeyNotFound
	}
	return val, err
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
func (r *RedisDb) Get(key string) (string, error) {
	val, err := r.db.Get(key).Result()
	// 统一错误处理：将Redis特定错误转换为接口标准错误
	if errors.Is(err, redis.Nil) {
		return "", _interface.ErrKeyNotFound
	}
	return val, err
}

// HGet 获取哈希表中指定field的值
// 参数：
//
//	key - 哈希表键名
//	field - 字段名
//
// 返回值：
//
//	string - 字段对应的值
//	error - 操作错误，字段不存在时返回ErrKeyNotFound
func (r *RedisDb) HGet(key, field string) (string, error) {
	val, err := r.db.HGet(key, field).Result()
	// 统一错误处理：将Redis特定错误转换为接口标准错误
	if errors.Is(err, redis.Nil) {
		return "", _interface.ErrKeyNotFound
	}
	return val, err
}

type RedisTx struct {
	pipe redis.Pipeliner
}

func (tx *RedisTx) Commit() error {
	_, err := tx.pipe.Exec()
	return err
}

func (tx *RedisTx) Rollback() error {
	return tx.pipe.Discard()
}

func (tx *RedisTx) Set(key string, value string, ttl time.Duration) error {
	return tx.pipe.Set(key, value, ttl).Err()
}

func (tx *RedisTx) Delete(key string) error {
	return tx.pipe.Del(key).Err()
}

func (tx *RedisTx) Expire(key string, ttl time.Duration) error {
	return tx.pipe.Expire(key, ttl).Err()
}

func (tx *RedisTx) HSet(key, field, value string, ttl time.Duration) error {
	if err := tx.pipe.HSet(key, field, value).Err(); err != nil {
		return err
	}
	if ttl > 0 {
		return tx.pipe.Expire(key, ttl).Err()
	}
	return nil
}

func (tx *RedisTx) HDel(key, field string) error {
	return tx.pipe.HDel(key, field).Err()
}

func (r *RedisDb) Close() {
	_ = r.db.Close()
}

func (r *RedisDb) Set(key string, value string, ttl time.Duration) error {
	return r.db.Set(key, value, ttl).Err()
}

func (r *RedisDb) Delete(key string) error {
	return r.db.Del(key).Err()
}

func (r *RedisDb) Exists(key string) (bool, error) {
	count, err := r.db.Exists(key).Result()
	return count > 0, err
}

func (r *RedisDb) Expire(key string, ttl time.Duration) error {
	return r.db.Expire(key, ttl).Err()
}

// HSet 设置哈希表中的field-value，并设置过期时间
// 参数：
//
//	key - 哈希表键名
//	field - 字段名
//	value - 字段值
//	ttl - 过期时间
//
// 返回值：
//
//	error - 操作错误
func (r *RedisDb) HSet(key, field, value string, ttl time.Duration) error {
	if err := r.db.HSet(key, field, value).Err(); err != nil {
		return err
	}
	if ttl > 0 {
		return r.db.Expire(key, ttl).Err()
	}
	return nil
}

// HDel 删除哈希表中的一个或多个field
// 参数：
//
//	key - 哈希表键名
//	field - 字段名
//
// 返回值：
//
//	error - 操作错误
func (r *RedisDb) HDel(key, field string) error {
	return r.db.HDel(key, field).Err()
}

// HGetAll 获取哈希表中所有的field和value
// 参数：
//
//	key - 哈希表键名
//
// 返回值：
//
//	map[string]string - 所有字段和值的映射
//	error - 操作错误
func (r *RedisDb) HGetAll(key string) (map[string]string, error) {
	return r.db.HGetAll(key).Result()
}

func (r *RedisDb) Push(key string, value string) error {
	return r.RPush(key, value)
}

func (r *RedisDb) Pop(key string) (string, error) {
	return r.LPop(key)
}

func (r *RedisDb) PopAll(key string) ([]string, error) {
	pipe := r.db.TxPipeline()
	lrange := pipe.LRange(key, 0, -1)
	pipe.Del(key)
	_, err := pipe.Exec()
	if err != nil {
		return nil, err
	}
	return lrange.Val(), nil
}

func (r *RedisDb) Len(key string) (int64, error) {
	return r.db.LLen(key).Result()
}

func (r *RedisDb) BeginTx() (_interface.Tx, error) {
	txPipe := r.db.TxPipeline()
	return &RedisTx{pipe: txPipe}, nil
}

func NewRedisClient(config config.Cache) (_interface.Cache, error) {
	addr := config.Host + ":" + config.Port
	redisDb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     200,
		MinIdleConns: 100,
	})
	if _, err := redisDb.Ping().Result(); err != nil {
		return nil, err
	}
	return &RedisDb{db: redisDb}, nil
}
