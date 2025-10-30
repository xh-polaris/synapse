package appinfra

import (
	"context"
	"fmt"

	"github.com/xh-polaris/synapse/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"github.com/xh-polaris/synapse/biz/infra/contract/orm"
	"github.com/xh-polaris/synapse/biz/infra/contract/risk"
	contsms "github.com/xh-polaris/synapse/biz/infra/contract/sms"
	"github.com/xh-polaris/synapse/biz/infra/impl/cache/redis"
	"github.com/xh-polaris/synapse/biz/infra/impl/mongoid"
	"github.com/xh-polaris/synapse/biz/infra/impl/mysql"
	"github.com/xh-polaris/synapse/biz/infra/impl/sms"
)

type AppDependencies struct {
	DB          *orm.DB
	IDGen       id.IDGenerator
	RiskManager *risk.RiskManager
	CacheCli    cache.Cmdable
	SMS         contsms.Provider
}

// Init 初始化
func Init(ctx context.Context) (_ *AppDependencies, err error) {
	infra := &AppDependencies{}
	infra.DB, err = mysql.New()
	if err != nil {
		return nil, fmt.Errorf("init [db] failed err:%v", err)
	}
	infra.IDGen = mongoid.New(ctx)
	infra.CacheCli = redis.New()
	infra.RiskManager = risk.New(infra.CacheCli)
	infra.SMS, err = sms.New(ctx, infra.CacheCli)
	if err != nil {
		return nil, fmt.Errorf("init infra [SMS] failed, err:%v", err)
	}
	return infra, nil
}
