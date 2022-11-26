package data

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGreeterRepo,
)

// Data .
type Data struct {
	DBMap map[string]*gorm.DB
	Cache map[string]*redis.Client
	lock  sync.RWMutex
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	database := c.GetDatabases()
	caches := c.GetRedis()

	data := &Data{
		DBMap: make(map[string]*gorm.DB),
		Cache: make(map[string]*redis.Client),
	}
	data.lock.RLock()
	// gorm.DB
	for _, dbConf := range database {
		// TODO
		data.DBMap[dbConf.GetDbName()] = nil
	}

	// redis.Client
	for _, cacheConf := range caches {
		// TODO
		data.Cache[cacheConf.GetDbName()] = nil
	}

	data.lock.RUnlock()

	return data, cleanup, nil
}
