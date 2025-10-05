package service

import (
	"context"

	"github.com/xh-polaris/synapse/biz/domain/basicuser/entity"
)

type BasicUser interface {
	PhoneNotExist(ctx context.Context, phone string) (err error)
	Register(ctx context.Context, authType, authID, password string) (*entity.BasicUser, error)
	LoginByPhone(ctx context.Context, requirePassword bool, phone, verify string) (*entity.BasicUser, error)
}
