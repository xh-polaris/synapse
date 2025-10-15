package repo

import (
	"context"

	"github.com/xh-polaris/synapse/biz/domain/thirdparty/dal"
	"github.com/xh-polaris/synapse/biz/domain/thirdparty/dal/model"
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"github.com/xh-polaris/synapse/biz/infra/contract/orm"
)

func NewThirdPartyUserRepo(db *orm.DB) ThirdPartyUserRepo {
	return dal.NewThirdPartyUserDAO(db)
}

type ThirdPartyUserRepo interface {
	FindOrCreate(ctx context.Context, id id.ID, code, app string) (*model.ThirdPartyUser, error)
}
