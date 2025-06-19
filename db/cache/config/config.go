// config包：缓存配置管理
// 定义了所有缓存驱动的统一配置结构和常量定义
//
// 本包提供了缓存系统的配置管理功能，支持多种缓存驱动的配置参数
// 通过统一的配置结构，简化了不同缓存实现的配置管理
//
// 支持的缓存驱动：
// - Redis：分布式内存缓存，支持集群和持久化
// - BadgerDB：高性能本地LSM树存储
// - BuntDB：快速内存数据库，支持持久化
//
// 配置参数说明：
// - Driver：缓存驱动类型标识
// - Path：本地存储路径（BadgerDB/BuntDB使用）
// - Host：服务器地址（Redis使用）
// - Port：服务器端口（Redis使用）
// - Password：认证密码（Redis使用）
// - DB：数据库编号（Redis使用）
//
// 使用示例：
//
//	cfg := config.Cache{
//	    Driver: config.CacheDriverRedis,
//	    Host:   "localhost",
//	    Port:   "6379",
//	}
//
// 作者: gophertool
package config

const (
	CacheDriverRedis  = "redis"
	CacheDriverBadger = "badger"
	CacheDriverBuntdb = "buntdb"
)

type Cache struct {
	Driver   string
	Path     string
	Host     string
	Port     string
	Password string
	DB       int
}
