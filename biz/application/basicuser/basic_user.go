package basicuser

import (
	"context"
	"fmt"
	"strconv"

	model "github.com/xh-polaris/synapse/biz/api/model/basicuser"
	"github.com/xh-polaris/synapse/biz/application/base/token"
	"github.com/xh-polaris/synapse/biz/application/basicuser/internal"
	application "github.com/xh-polaris/synapse/biz/application/internal"
	"github.com/xh-polaris/synapse/biz/conf"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/entity"
	basicuser "github.com/xh-polaris/synapse/biz/domain/basicuser/service"
	"github.com/xh-polaris/synapse/biz/infra/contract/risk"
	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
	ctxcache "github.com/xh-polaris/synapse/biz/pkg/ctxcache/ctx_cache"
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
	"github.com/xh-polaris/synapse/biz/pkg/logs"
	"github.com/xh-polaris/synapse/biz/types/cst"
	"github.com/xh-polaris/synapse/biz/types/errno"
)

var BasicUserSVC = &BasicUserService{}

type BasicUserService struct {
	sms       sms.Provider
	DomainSVC basicuser.BasicUser
}

// RegisterNewBasicUser 注册一个新用户
func (s *BasicUserService) RegisterNewBasicUser(ctx context.Context, req *model.BasicUserRegisterReq) (resp *model.BasicUserRegisterResp, err error) {
	if err = conf.ValidApp(req.GetApp()); err != nil {
		return nil, err
	}

	switch req.AuthType {
	case cst.AuthTypePhoneVerify:
		if err = s.validPhoneVerify(ctx, req.App.Name, req.AuthId, req.Verify); err != nil {
			return nil, err
		}
		ok, err := s.DomainSVC.PhoneExist(ctx, req.AuthId)
		if err != nil {
			return nil, err
		}
		if ok {
			return nil, errorx.New(errno.PhoneHasExisted, errorx.KV("phone", req.AuthId))
		}
	case cst.AuthTypeIdPassword:
		if req.ExtraAuthId == nil {
			return nil, errorx.New(errno.MissingParameter, errorx.KV("parameter", "学号"))
		}
		ok, err := s.DomainSVC.StudentIDExist(ctx, req.AuthId, *req.ExtraAuthId)
		if err != nil {
			return nil, err
		}
		if ok {
			return nil, errorx.New(errno.StudentIDHasExisted, errorx.KV("studentId", *req.ExtraAuthId))
		}
	default:
		return nil, errorx.New(errno.UnSupportAuthType, errorx.KV("type", req.AuthType))
	}

	if req.GetPassword() == "" {
		return nil, errorx.New(errno.MustPassword)
	}

	var extraAuthId string
	if req.ExtraAuthId != nil {
		extraAuthId = *req.ExtraAuthId
	}
	var u *entity.BasicUser
	if u, err = s.DomainSVC.Register(ctx, req.AuthType, req.AuthId, extraAuthId, *req.Password); err != nil {
		return nil, err
	}

	info := &token.Info{BasicUserId: u.ID}
	jwt, err := token.SignJWT(conf.GetConfig().Token, info)
	if err != nil {
		return nil, err
	}

	resp = &model.BasicUserRegisterResp{
		Resp:      application.Success(),
		Token:     jwt,
		BasicUser: internal.BasicUserPO2VO(u),
	}
	return
}

func (s *BasicUserService) validPhoneVerify(ctx context.Context, app, phone, code string) error {
	// 判断是否到上限
	key := fmt.Sprintf("risk:login:passport:%s", phone)
	limit, _, err := risk.CheckUpperLimit(ctx, key, conf.GetConfig().Token.MaxInPeriod)
	if err != nil {
		return err
	}
	if limit { // 达到上限, 不允许校验
		return errorx.New(errno.TooOftenLoginError, errorx.KV("period", strconv.Itoa(conf.GetConfig().SMS.Period/60)))
	}
	ok, err := s.sms.Check(ctx, app, "passport", phone, code)
	if err != nil {
		return err
	}
	if !ok {
		if err = risk.AddOnce(ctx, key, conf.GetConfig().Token.Period); err != nil {
			logs.Errorf("record send verify err:%s", err)
		}
		return errorx.New(errno.ErrVerifyCode)
	}
	return err
}

func (s *BasicUserService) Login(ctx context.Context, req *model.BasicUserLoginReq) (resp *model.BasicUserLoginResp, err error) {
	if req.App == nil {
		return nil, errorx.New(errno.MissingParameter, errorx.KV("parameter", "app"))
	}
	if err = conf.ValidApp(req.GetApp()); err != nil {
		return nil, err
	}

	var ok = true
	var u *entity.BasicUser
	switch req.AuthType {
	case cst.AuthTypePhoneVerify:
		err = s.validPhoneVerify(ctx, req.App.Name, req.AuthId, req.Verify)
		if err != nil {
			return nil, err
		}
		ok, err = s.DomainSVC.PhoneExist(ctx, req.AuthId)
		if err != nil {
			return nil, err
		}
		if ok {
			u, err = s.DomainSVC.LoginByPhone(ctx, false, req.AuthId, req.Verify)
		} else {
			u, err = s.DomainSVC.Register(ctx, req.AuthType, req.AuthId, "", "")
		}
	case cst.AuthTypePhonePassword:
		u, err = s.DomainSVC.LoginByPhone(ctx, true, req.AuthId, req.Verify)
	case cst.AuthTypeIdPassword:
		if req.ExtraAuthId == nil {
			return nil, errorx.New(errno.MissingParameter, errorx.KV("parameter", "学号"))
		}
		u, err = s.DomainSVC.LoginByStudentID(ctx, req.AuthId, *req.ExtraAuthId, req.Verify)
	default:
		return nil, errorx.New(errno.UnSupportAuthType, errorx.KV("type", req.AuthType))
	}
	if err != nil {
		return nil, err
	}

	info := &token.Info{BasicUserId: u.ID}
	jwt, err := token.SignJWT(conf.GetConfig().Token, info)
	if err != nil {
		return nil, err
	}

	resp = &model.BasicUserLoginResp{
		Resp:      application.Success(),
		Token:     jwt,
		BasicUser: internal.BasicUserPO2VO(u),
		New:       !ok,
	}
	return
}

func (s *BasicUserService) ResetPassword(ctx context.Context, req *model.BasicUserResetPasswordReq) (resp *model.BasicUserResetPasswordResp, err error) {
	if err = conf.ValidApp(req.GetApp()); err != nil {
		return nil, err
	}
	info, _ := ctxcache.Get[*token.Info](ctx, cst.TokenInfo)

	switch req.ResetCode {
	case "":
		if err = s.DomainSVC.ResetPassword(ctx, info.BasicUserId, req.NewPassword); err != nil {
			return nil, errorx.New(errno.ErrResetPassword)
		}
	default:
		if err, ok := conf.VerifyResetCode(req.GetApp(), req.ResetCode); err != nil {
			return nil, err
		} else if !ok {
			return nil, errorx.New(errno.ErrResetPassword)
		}
		if err = s.DomainSVC.ResetPassword(ctx, info.BasicUserId, req.NewPassword); err != nil {
			return nil, errorx.New(errno.ErrResetPassword)
		}
	}

	return &model.BasicUserResetPasswordResp{
		Resp: application.Success(),
	}, nil
}
