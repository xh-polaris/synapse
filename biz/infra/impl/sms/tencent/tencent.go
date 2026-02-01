package tencent

import (
	"context"
	"errors"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencenterrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/xh-polaris/synapse/biz/conf"
	"github.com/xh-polaris/synapse/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
	"github.com/xh-polaris/synapse/biz/pkg/logs"
)

type tencentSMS struct {
	cache  *sms.Cache
	client *tencentsms.Client
}

func New(ctx context.Context, cache *sms.Cache, account, token string) (sms.Provider, error) {
	s, err := getTencentSMSProvider(ctx, cache, account, token)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func getTencentSMSProvider(_ context.Context, cache *sms.Cache, secretId, secretKey string) (sms.Provider, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tencentsms.NewClient(credential, "ap-guangzhou", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	return &tencentSMS{cache: cache, client: client}, nil
}

func (t *tencentSMS) Send(ctx context.Context, app, cause, phone string, param *sms.SMSParam) (err error) {
	// 发送短信
	if _, err = t.send(ctx, app, cause, phone, param); err != nil {
		return err
	}
	expire := param.Expire

	if err = t.cache.Store(ctx, app, cause, phone, param.Code, expire); err != nil {
		return err
	}
	return nil
}

func (t *tencentSMS) send(_ context.Context, app, cause, phone string, param *sms.SMSParam) (map[string]any, error) {
	// 参数设置
	req := tencentsms.NewSendSmsRequest()
	req.SmsSdkAppId = common.StringPtr(conf.GetConfig().SMS.Extra["AppId"])     // 应用ID
	req.SignName = common.StringPtr(conf.GetConfig().SMS.Extra["Sign"])         // 签名内容
	req.TemplateId = common.StringPtr(conf.GetConfig().SMS.AppConf[app][cause]) //模板ID
	req.TemplateParamSet = common.StringPtrs([]string{param.Code, strconv.Itoa(int(param.Expire.Minutes()))})
	req.PhoneNumberSet = common.StringPtrs([]string{phone})

	// 发送响应
	resp, err := t.client.SendSms(req)

	// SDK 错误
	var tencentCloudSDKError *tencenterrors.TencentCloudSDKError
	if errors.As(err, &tencentCloudSDKError) {
		logs.Errorf("An Tencent API error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	// 正常响应
	r, respStr := map[string]any{}, resp.ToJsonString()
	if err = json.Unmarshal([]byte(respStr), &r); err != nil {
		logs.Infof("Tencent sms return resp %s but unmarsharl failed %s", respStr, err)
	}
	return r, nil
}
func (t *tencentSMS) Check(ctx context.Context, app, cause, phone, code string) (bool, error) {
	if code == "xh-polaris" && conf.GetConfig().State == "test" {
		return true, nil
	}
	ori, err := t.cache.Load(ctx, app, cause, phone)
	if errors.Is(err, cache.Nil) {
		return false, nil
	}
	return ori == code, err
}
