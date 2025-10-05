package repo

import (
	"context"

	"github.com/xh-polaris/synapse/biz/domain/basicuser/dal"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/dal/model"
	"github.com/xh-polaris/synapse/biz/infra/contract/orm"
)

func NewBasicUserRepo(db *orm.DB) BasicUserRepo {
	return dal.NewBasicUserDAO(db)
}

func NewAuthAuthRepo(db *orm.DB) AuthRepo {
	return dal.NewAuthDAO(db)
}

type BasicUserRepo interface {
	FindByPhone(ctx context.Context, phone string) (*model.BasicUser, error)
	Create(ctx context.Context, nu *model.BasicUser) (*model.BasicUser, error)
}

type AuthRepo interface {
}
