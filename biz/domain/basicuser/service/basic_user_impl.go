package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/xh-polaris/synapse/biz/conf"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/dal/model"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/entity"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/repo"
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"github.com/xh-polaris/synapse/biz/infra/contract/risk"
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
	"github.com/xh-polaris/synapse/biz/pkg/lang/crypt"
	"github.com/xh-polaris/synapse/biz/pkg/logs"
	"github.com/xh-polaris/synapse/biz/types/cst"
	"github.com/xh-polaris/synapse/biz/types/errno"
)

type Component struct {
	BasicUserRepo repo.BasicUserRepo
	AuthRepo      repo.AuthRepo
	SchoolRepo    repo.SchoolRepo
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
	if u == nil { // 未注册过
		return nil, errorx.New(errno.PhoneNotExisted)
	}
	if requirePassword {
		if err = loginLimiter(ctx, verify, u.Password, phone); err != nil {
			return nil, err
		}
	}
	return basicUserModel2Entity(u)
}

func (i *userImpl) LoginByStudentID(ctx context.Context, schoolId, studentId, verify string) (*entity.BasicUser, error) {
	u, err := i.BasicUserRepo.FindByStudentID(ctx, schoolId, studentId)
	if err != nil {
		return nil, err
	}
	if u == nil { // 未注册过
		return nil, errorx.New(errno.StudentIDNotExisted)
	}
	if err = loginLimiter(ctx, verify, u.Password, schoolId, studentId); err != nil {
		return nil, err
	}
	return basicUserModel2Entity(u)
}

func (i *userImpl) PhoneExist(ctx context.Context, phone string) (is bool, err error) {
	mu, err := i.BasicUserRepo.FindByPhone(ctx, phone)
	if err != nil {
		return false, err
	}
	if mu != nil {
		return true, nil
	}
	return false, nil
}

func (i *userImpl) StudentIDExist(ctx context.Context, schoolId, studentId string) (is bool, err error) {
	mus, err := i.BasicUserRepo.FindManyBySchoolID(ctx, schoolId)
	if err != nil {
		return false, err
	}
	for _, mu := range mus {
		if mu.StudentID != nil && *mu.StudentID == studentId {
			return true, nil
		}
	}
	return false, nil
}

func (i *userImpl) Register(ctx context.Context, authType, authId, extraAuthId, password string) (u *entity.BasicUser, err error) {
	var hashed string
	if password != "" {
		hashed, err = crypt.Hash(password)
		if err != nil {
			return nil, err
		}
	}

	nu := &model.BasicUser{
		ID:       i.IdGen.GenID(ctx),
		Password: &hashed,
	}

	switch authType {
	case cst.AuthTypePhoneVerify:
		nu.Phone = &authId
	case cst.AuthTypeIdPassword:
		aid, err := id.FromHex(authId)
		if err != nil {
			return nil, err
		}
		nu.SchoolID = &aid
		nu.StudentID = &extraAuthId

	default:
		return nil, errorx.New(errno.UnSupportAuthType, errorx.KV("type", authType))
	}
	nu, err = i.BasicUserRepo.Create(ctx, nu)
	if err != nil {
		return nil, errorx.New(errno.ErrRegister)
	}
	return basicUserModel2Entity(nu)
}

func (i *userImpl) ResetPassword(ctx context.Context, basicUserId string, password string) error {
	if password == "" {
		return errorx.New(errno.MustPassword)
	}
	hashed, err := crypt.Hash(password)
	if err != nil {
		return err
	}
	return i.BasicUserRepo.ResetPassword(ctx, basicUserId, hashed)
}

func loginLimiter(ctx context.Context, password string, hashed *string, parts ...string) error {
	key := "risk:login:passport:" + strings.Join(parts, ",")
	limit, _, err := risk.CheckUpperLimit(ctx, key, conf.GetConfig().Token.MaxInPeriod)
	if err != nil {
		return err
	}
	if limit { // 达到上限, 不允许校验
		return errorx.New(errno.TooOftenLoginError, errorx.KV("period", strconv.Itoa(conf.GetConfig().SMS.Period/60)))
	}
	if hashed == nil || *hashed == "" {
		return errorx.New(errno.NoPassword)
	}
	if !crypt.Check(password, *hashed) {
		if err = risk.AddOnce(ctx, key, conf.GetConfig().Token.Period); err != nil {
			logs.Errorf("record send verify err:%s", err)
		}
		return errorx.New(errno.ErrPassword)
	}
	return nil
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
