package system

import (
	"context"

	"github.com/xh-polaris/synapse/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
)

func InitService(ctx context.Context, sms sms.Provider, cache cache.Cmdable) *SystemService {
	SystemSVC.sms = sms
	SystemSVC.cache = cache
	return SystemSVC
}
