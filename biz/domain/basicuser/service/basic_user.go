package service

import (
	"context"

	"github.com/xh-polaris/synapse/biz/domain/basicuser/entity"
)

type BasicUser interface {
	PhoneExist(ctx context.Context, phone string) (is bool, err error)
	CodeExist(ctx context.Context, schoolId, code string) (is bool, err error)
	Register(ctx context.Context, authType, authId, extraAuthId, password string) (*entity.BasicUser, error)
	LoginByPhone(ctx context.Context, requirePassword bool, phone, verify string) (*entity.BasicUser, error)
	LoginByCode(ctx context.Context, schoolId, code, verify string) (*entity.BasicUser, error)
	ResetPassword(ctx context.Context, basicUserId string, password string) error
}
