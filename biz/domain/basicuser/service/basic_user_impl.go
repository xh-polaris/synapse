package service

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/dal/model"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/entity"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/repo"
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
	"github.com/xh-polaris/synapse/biz/pkg/lang/crypt"
	"github.com/xh-polaris/synapse/biz/types/cst"
	"github.com/xh-polaris/synapse/biz/types/errno"
)

type Component struct {
	BasicUserRepo repo.BasicUserRepo
	AuthRepo      repo.AuthRepo
	IdGen         id.IDGenerator
}

func NewBasicUserDomain(ctx context.Context, c *Component) BasicUser {
	return &userImpl{Component: c}
}

type userImpl struct {
	*Component
}

func (i *userImpl) LoginByPhone(ctx context.Context, requirePassword bool, phone, verify string) (*entity.BasicUser, error) {
	u, err := i.BasicUserRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errorx.New(errno.PhoneNotExisted, errorx.KV("phone", phone))
	}
	if requirePassword {
		if u.Password == nil {
			return nil, errorx.New(errno.NoPassword)
		}
		if !crypt.Check(verify, *u.Password) {
			return nil, errorx.New(errno.ErrPassword)
		}
	}
	return basicUserModel2Entity(u)
}

func (i *userImpl) PhoneNotExist(ctx context.Context, phone string) (err error) {
	mu, err := i.BasicUserRepo.FindByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if mu != nil {
		return errorx.New(errno.PhoneHasExisted, errorx.KV("phone", phone))
	}
	return nil
}

func (i *userImpl) Register(ctx context.Context, authType, authID, password string) (*entity.BasicUser, error) {
	hashed, err := crypt.Hash(password)
	if err != nil {
		return nil, err
	}

	nu := &model.BasicUser{
		ID:       i.IdGen.GenID(ctx),
		Password: &hashed,
	}

	switch authType {
	case cst.AuthTypePhoneVerify:
		nu.Phone = &authID
	}
	nu, err = i.BasicUserRepo.Create(ctx, nu)
	if err != nil {
		return nil, errorx.New(errno.ErrRegister)
	}
	return basicUserModel2Entity(nu)
}

func basicUserModel2Entity(u *model.BasicUser) (*entity.BasicUser, error) {
	eu := &entity.BasicUser{
		ID:        u.ID.Hex(),
		Code:      u.Code,
		Phone:     u.Phone,
		Password:  u.Password,
		Name:      u.Name,
		Gender:    u.Gender,
		CreatedAt: time.UnixMilli(u.CreatedAt),
		UpdatedAt: time.UnixMilli(u.UpdatedAt),
	}
	if len(u.Extra) != 0 {
		extra := map[string]any{}
		err := sonic.Unmarshal([]byte(u.Extra.String()), &extra)
		if err != nil {
			return nil, err
		}
		eu.Extra = &extra
	}
	return eu, nil
}
