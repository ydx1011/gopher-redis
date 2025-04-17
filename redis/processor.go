package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/xfali/xlog"
	"github.com/ydx1011/gopher-core/bean"
	"github.com/ydx1011/yfig"
	"time"
)

const (
	BuildinValueRedisSources = "gopher.redisSources"
)

type Sources struct {
	Prefix   string
	Addrs    []string
	Password string
	// 单节点模式生效
	DB              int
	PoolSize        int
	IsCluster       string
	RouteRandomly   bool
	ReadOnly        bool
	MinIdleConn     int
	DialTimeoutSec  int
	ReadTimeoutSec  int
	WriteTimeoutSec int
	PoolFIFO        bool
	// 最大生命周�����, 超过释放回收
	MaxConnAgeSec int
	// 连接池最大超时时间, 超过重置
	PoolTimeoutSec int
	// 最大空闲时间, 超过释放
	IdleTimeoutSec int
	// 检测空闲连接频率, 默认一分钟
	IdleCheckFrequencySec int
}

type Processor struct {
	logger xlog.Logger
}

func (p *Processor) Init(conf yfig.Properties, container bean.Container) error {
	dss := map[string]*Sources{}
	err := conf.GetValue(BuildinValueRedisSources, &dss)
	if len(dss) == 0 {
		p.logger.Errorln("No Database")
		return nil
	}
	for k, v := range dss {
		if v.IsCluster == "true" {
			client := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:         v.Addrs,
				Password:      v.Password,
				PoolSize:      v.PoolSize,
				ReadOnly:      v.ReadOnly,
				RouteRandomly: v.RouteRandomly,
				MinIdleConns:  v.MinIdleConn,
				PoolFIFO:      v.PoolFIFO,
				DialTimeout:   time.Duration(v.DialTimeoutSec) * time.Second,
				ReadTimeout:   time.Duration(v.ReadTimeoutSec) * time.Second,
				WriteTimeout:  time.Duration(v.WriteTimeoutSec) * time.Second,
				PoolTimeout:   time.Duration(v.PoolTimeoutSec) * time.Second,
			})
			//添加到注入容器
			container.RegisterByName(k, client)
		} else {
			client := redis.NewClient(&redis.Options{
				Addr:         v.Addrs[0],
				Password:     v.Password,
				DB:           v.DB,
				PoolSize:     v.PoolSize,
				MinIdleConns: v.MinIdleConn,
				DialTimeout:  time.Duration(v.DialTimeoutSec) * time.Second,
				ReadTimeout:  time.Duration(v.ReadTimeoutSec) * time.Second,
				WriteTimeout: time.Duration(v.WriteTimeoutSec) * time.Second,
				PoolFIFO:     v.PoolFIFO,
				PoolTimeout:  time.Duration(v.PoolTimeoutSec) * time.Second,
			})
			//添加到注入容器
			container.RegisterByName(k, client)
		}

	}

	return err
}

func (p *Processor) Process() error {
	return nil
}

func (p *Processor) Classify(o interface{}) (bool, error) {
	//switch v := o.(type) {
	//case redis.Client:
	//	err := p.parseBean(v)
	//	return true, err
	//}
	return false, nil
}

func (p *Processor) BeanDestroy() error {
	return nil
}
