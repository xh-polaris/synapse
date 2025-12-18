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

func NewSchoolRepo(db *orm.DB) SchoolRepo {
	return dal.NewSchoolDAO(db)
}

type BasicUserRepo interface {
	FindByID(ctx context.Context, id string) (*model.BasicUser, error)
	FindByPhone(ctx context.Context, phone string) (*model.BasicUser, error)
	FindManyBySchoolID(ctx context.Context, schoolId string) ([]*model.BasicUser, error)
	FindByCode(ctx context.Context, schoolId, code string) (*model.BasicUser, error)
	Create(ctx context.Context, nu *model.BasicUser) (*model.BasicUser, error)
	ResetPassword(ctx context.Context, basicUserId, password string) error
}

type AuthRepo interface {
}

type SchoolRepo interface {
	FindByID(ctx context.Context, id string) (*model.School, error)
}
