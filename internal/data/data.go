package data

import (
	"context"
	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	kGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	clientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	GetRPCConn,
	GetGreeterClient,
	NewData,
	NewGreeterRepo,
)

// Data .
type Data struct {
	DBMap map[string]*gorm.DB
	Cache map[string]*redis.Client
	lock  sync.RWMutex
	hc    v1.GreeterClient
}

// NewData .
func NewData(c *conf.Data, hc v1.GreeterClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	database := c.GetDatabases()
	caches := c.GetRedis()

	data := &Data{
		DBMap: make(map[string]*gorm.DB),
		Cache: make(map[string]*redis.Client),
		hc:    hc,
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

// GetRPCConn 获取rpc连接
func GetRPCConn(etcdClient *clientV3.Client, discovery *conf.Discovery) *grpc.ClientConn {
	// new dis with etcd client
	dis := etcd.New(etcdClient)
	endpoint := "discovery:///provider"
	conn, err := kGrpc.Dial(context.Background(), kGrpc.WithEndpoint(endpoint), kGrpc.WithDiscovery(dis))
	if err != nil {
		panic(err)
	}

	return conn
}

// GetGreeterClient 获取GreeterClient
func GetGreeterClient(conn *grpc.ClientConn) v1.GreeterClient {
	return v1.NewGreeterClient(conn)
}
