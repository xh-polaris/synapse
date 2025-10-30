package sms

import (
	"context"
	"fmt"
	"strconv"

	"github.com/xh-polaris/synapse/biz/conf"
	"github.com/xh-polaris/synapse/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse/biz/infra/contract/risk"
	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
	"github.com/xh-polaris/synapse/biz/pkg/logs"
	"github.com/xh-polaris/synapse/biz/types/errno"
)

type SafeSMSProvider struct {
	sms.Provider
	Cache cache.Cmdable
}

func NewSafeSMSProvider(provider sms.Provider, cacheCli cache.Cmdable) (*SafeSMSProvider, error) {
	return &SafeSMSProvider{
		Provider: provider,
		Cache:    cacheCli,
	}, nil
}

// Send 发送验证码
func (s *SafeSMSProvider) Send(ctx context.Context, app, cause, phone string, param *sms.SMSParam) error {
	// 判断是否到上限
	key := fmt.Sprintf("risk:sendVerifyCode:%s:%s:%S" + app + cause + phone)
	limit, _, err := risk.CheckUpperLimit(ctx, key, conf.GetConfig().SMS.MaxInPeriod)
	if err != nil {
		return err
	}
	if limit { // 达到上限, 不允许发送
		return errorx.New(errno.ErrSendUpperLimit, errorx.KV("period", strconv.Itoa(conf.GetConfig().SMS.Period)))
	}
	// 发送验证码
	err = s.Provider.Send(ctx, app, cause, phone, param)
	if err != nil {
		return err
	}
	// 记录操作
	if err = risk.AddOnce(ctx, key, conf.GetConfig().SMS.Period); err != nil {
		logs.Errorf("record send verify err:%s", err)
	}
	return nil
}

// Check 校验验证码
func (s *SafeSMSProvider) Check(ctx context.Context, app, cause, phone, code string) (bool, error) {
	return s.Provider.Check(ctx, app, cause, phone, code)
}
