package thirdparty

import (
	"context"

	"github.com/xh-polaris/synapse/biz/domain/thirdparty/repo"
	"github.com/xh-polaris/synapse/biz/domain/thirdparty/service"
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"github.com/xh-polaris/synapse/biz/infra/contract/orm"
)

func InitService(ctx context.Context, db *orm.DB, idGen id.IDGenerator) *ThirdPartyService {
	ThirdPartySVC.DomainSVC = service.NewThirdPartyDomain(ctx, &service.Component{
		ThirdPartyUserRepo: repo.NewThirdPartyUserRepo(db),
		IdGen:              idGen,
	})
	return ThirdPartySVC
}
