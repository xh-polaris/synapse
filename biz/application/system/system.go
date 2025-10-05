package system

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/xh-polaris/synapse/biz/api/model/system"
	"github.com/xh-polaris/synapse/biz/application/internal"
	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
	"github.com/xh-polaris/synapse/biz/pkg/logs"
	"github.com/xh-polaris/synapse/biz/types/cst"
	"github.com/xh-polaris/synapse/biz/types/errno"
)

var SystemSVC = &SystemService{}

type SystemService struct {
	sms sms.Provider
}

func (s *SystemService) Send(ctx context.Context, req *system.SendVerifyCodeReq) (*system.SendVerifyCodeResp, error) {
	if req.AuthType != cst.AuthTypePhoneVerify {
		return nil, errorx.New(errno.ErrInvalidAuthType)
	}

	param := &sms.SMSParam{Code: genCode(), Expire: time.Duration(req.Expire) * time.Second}
	if err := s.sms.Send(ctx, req.App.Name, req.Cause, req.AuthId, param); err != nil {
		logs.Error(err)
		return nil, errorx.New(errno.ErrSendPhoneVerify)
	}
	return &system.SendVerifyCodeResp{Resp: internal.Success()}, nil
}

// 生成n位随机验证码
func genCode() string {
	return strconv.Itoa(rand.Intn(999999-100000) + 100000)
}

func (s *SystemService) Check(ctx context.Context, req *system.CheckVerifyCodeReq) (*system.CheckVerifyCodeResp, error) {
	if req.AuthType != cst.AuthTypePhoneVerify {
		return nil, errorx.New(errno.ErrInvalidAuthType)
	}

	check, err := s.sms.Check(ctx, req.App.Name, req.Cause, req.AuthId, req.Verify)
	if err != nil {
		return nil, err
	}
	return &system.CheckVerifyCodeResp{Resp: internal.Success(), Verify: check}, nil
}
