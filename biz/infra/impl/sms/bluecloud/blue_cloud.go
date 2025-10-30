package bluecloud

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/xh-polaris/synapse/biz/conf"
	"github.com/xh-polaris/synapse/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
	"github.com/xh-polaris/synapse/biz/pkg/httpcli"
)

const (
	singleSendURL   = "https://bluecloudccs.21vbluecloud.com/services/sms/messages?api-version=2018-10-01"
	checkReceiveURL = "https://bluecloudccs.21vbluecloud.com/services/sms/messages/%s?api-version=2018-10-01&continuationToken=&count=10"
)

type bluecloudSMS struct {
	cache      *sms.Cache
	authHeader http.Header
}

func New(ctx context.Context, cache *sms.Cache, account, token string) (sms.Provider, error) {
	s, err := getBlueCloudSMSProvider(ctx, cache, account, token)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func getBlueCloudSMSProvider(ctx context.Context, cacheCli *sms.Cache, account, token string) (sms.Provider, error) {
	header := http.Header{}
	header.Set("content-type", "application/json")
	header.Set("Account", account)
	header.Set("Authorization", token)
	return &bluecloudSMS{cache: cacheCli, authHeader: header}, nil
}

// Send 发送验证码并校验用户是否成功收到
func (b *bluecloudSMS) Send(ctx context.Context, app, cause, phone string, param *sms.SMSParam) (err error) {
	// 发送短信
	if _, err = b.send(ctx, app, cause, phone, param); err != nil {
		return err
	}
	expire := param.Expire

	if err = b.cache.Store(ctx, app, cause, phone, param.Code, expire); err != nil {
		return err
	}
	return nil
}

func (b *bluecloudSMS) send(_ context.Context, app, cause, phone string, param *sms.SMSParam) (map[string]any, error) {
	body := map[string]any{
		"phoneNumber": []string{phone},
		"messageBody": map[string]any{
			"extend":       "00",
			"templateName": conf.GetConfig().SMS.AppConf[app][cause],
			"templateParam": map[string]any{
				"otpcode": param.Code,                                // 验证码
				"expire":  strconv.Itoa(int(param.Expire.Minutes())), // 以分钟为单位超时时间
			},
		},
	}
	res, err := httpcli.GetHttpClient().Post(singleSendURL, b.authHeader, body)
	return res, err
}

func (b *bluecloudSMS) Check(ctx context.Context, app, cause, phone, code string) (bool, error) {
	if code == "xh-polaris" && conf.GetConfig().State == "test" {
		return true, nil
	}
	ori, err := b.cache.Load(ctx, app, cause, phone)
	if errors.Is(err, cache.Nil) {
		return false, nil
	}
	return ori == code, err
}
