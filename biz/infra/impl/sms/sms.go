package sms

import (
	"context"
	"fmt"

	"github.com/xh-polaris/synapse/biz/conf"
	"github.com/xh-polaris/synapse/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
	"github.com/xh-polaris/synapse/biz/infra/impl/sms/bluecloud"
)

const (
	BlueCloud = "blue-cloud"
)

func New(ctx context.Context, cacheCli cache.Cmdable) (sms.Provider, error) {
	c := conf.GetConfig().SMS
	ch := sms.NewSMSCache(ctx, cacheCli)

	switch c.Provider {
	case BlueCloud:
		return bluecloud.New(ctx, ch, c.Account, c.Token)
	}
	return nil, fmt.Errorf("no such SMS provider: %s", c.Provider)
}
