package system

import (
	"context"

	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
)

func InitService(ctx context.Context, sms sms.Provider) *SystemService {
	SystemSVC.sms = sms
	return SystemSVC
}
