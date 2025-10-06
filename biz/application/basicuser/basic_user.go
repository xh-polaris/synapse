package basicuser

import (
	"context"

	model "github.com/xh-polaris/synapse/biz/api/model/basicuser"
	"github.com/xh-polaris/synapse/biz/application/base/token"
	"github.com/xh-polaris/synapse/biz/application/basicuser/internal"
	application "github.com/xh-polaris/synapse/biz/application/internal"
	"github.com/xh-polaris/synapse/biz/conf"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/entity"
	basicuser "github.com/xh-polaris/synapse/biz/domain/basicuser/service"
	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
	ctxcache "github.com/xh-polaris/synapse/biz/pkg/ctxcache/ctx_cache"
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
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
	if err = conf.ValidApp(req.App.Name); err != nil {
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
	default:
		return nil, errorx.New(errno.UnSupportAuthType, errorx.KV("type", req.AuthType))
	}

	if req.GetPassword() == "" {
		return nil, errorx.New(errno.MustPassword)
	}
	var u *entity.BasicUser
	if u, err = s.DomainSVC.Register(ctx, req.AuthType, req.AuthId, *req.Password); err != nil {
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
	ok, err := s.sms.Check(ctx, app, "passport", phone, code)
	if err != nil {
		return err
	}
	if !ok {
		return errorx.New(errno.ErrVerifyCode)
	}
	return err
}

func (s *BasicUserService) Login(ctx context.Context, req *model.BasicUserLoginReq) (resp *model.BasicUserLoginResp, err error) {
	if err = conf.ValidApp(req.App.Name); err != nil {
		return nil, err
	}

	var ok bool
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
			u, err = s.DomainSVC.Register(ctx, req.AuthType, req.AuthId, "")
		}
	case cst.AuthTypePhonePassword:
		u, err = s.DomainSVC.LoginByPhone(ctx, true, req.AuthId, req.Verify)
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
	if err = conf.ValidApp(req.App.Name); err != nil {
		return nil, err
	}
	info, _ := ctxcache.Get[*token.Info](ctx, cst.TokenInfo)

	if err = s.DomainSVC.ResetPassword(ctx, info.BasicUserId, req.NewPassword); err != nil {
		return nil, errorx.New(errno.ErrResetPassword)
	}

	return &model.BasicUserResetPasswordResp{
		Resp: application.Success(),
	}, nil
}
