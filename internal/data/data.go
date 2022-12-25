package data

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
)

// Data .
type Data struct {
	db    *gorm.DB
	cache *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{
		db:    GetMysqlDB(c.GetDatabase()), // TODO 需要替换GORM logger, 需要基于log.Logger实现gorm logger.Interface
		cache: GetRedisClient(c.GetRedis()),
	}

	cleanup := func() {
		aLog := log.NewHelper(logger)
		aLog.Info("closing the data resources")
		// 关闭数据库连接
		sqlDB, err := data.db.DB()
		if err != nil {
			aLog.Errorf("closing the data resources error: %v", err)
		}
		err = sqlDB.Close()
		if err != nil {
			aLog.Errorf("closing the data resources error: %v", err)
		}
		// 关闭redis连接
		err = data.cache.Close()
		if err != nil {
			aLog.Errorf("closing the data resources error: %v", err)
		}
	}
	return data, cleanup, nil
}

// GetMysqlDB 获取mysql数据库连接
func GetMysqlDB(cfg *conf.Data_Database, logger ...logger.Interface) *gorm.DB {
	var opts []gorm.Option
	if len(logger) > 0 {
		opts = append(opts, &gorm.Config{Logger: logger[0]})
	}

	conn, err := gorm.Open(mysql.Open(cfg.GetDsn()), opts...)
	if err != nil {
		panic(err)
	}

	if cfg.GetDebug() {
		conn = conn.Debug()
	}

	return conn
}

// GetRedisClient 获取redis客户端
func GetRedisClient(cfg *conf.Data_Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.GetAddr(),
		Password:     cfg.GetPassword(),
		DB:           int(cfg.GetDb()),
		WriteTimeout: cfg.GetWriteTimeout().AsDuration(),
		ReadTimeout:  cfg.GetReadTimeout().AsDuration(),
		DialTimeout:  cfg.GetDialTimeout().AsDuration(),
	})
}
