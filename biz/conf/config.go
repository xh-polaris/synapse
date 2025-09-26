package conf

import (
	"os"
	"sync"

	confx "github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	service.ServiceConf
	ListenOn string
	State    string
	Auth     *Auth
	Mongo    *Mongo
	Cache    cache.CacheConf
	Redis    redis.RedisConf
}

type Auth struct {
	SecretKey    string
	PublicKey    string
	AccessExpire int64
}

type Mongo struct {
	URL string
	DB  string
}

func NewConfig() (*Config, error) {
	once.Do(func() {
		c := new(Config)
		path := os.Getenv("CONFIG_PATH")
		if path == "" {
			path = "etc/config.yaml"
		}
		if err := confx.Load(path, c); err != nil {
			panic(err)
		}
		if err := c.SetUp(); err != nil {
			panic(err)
		}
		config = c
	})
	return config, nil
}

func GetConfig() *Config {
	once.Do(func() {
		_, _ = NewConfig()
	})
	return config
}
