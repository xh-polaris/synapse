package basicuser

import (
	"context"

	"github.com/xh-polaris/synapse/biz/domain/basicuser/repo"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/service"
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"github.com/xh-polaris/synapse/biz/infra/contract/orm"
	"github.com/xh-polaris/synapse/biz/infra/contract/sms"
)

func InitService(ctx context.Context, sms sms.Provider, db *orm.DB, idGen id.IDGenerator) *BasicUserService {
	BasicUserSVC.sms = sms
	BasicUserSVC.DomainSVC = service.NewBasicUserDomain(ctx, &service.Component{
		BasicUserRepo: repo.NewBasicUserRepo(db),
		AuthRepo:      repo.NewAuthAuthRepo(db),
		SchoolRepo:    repo.NewSchoolRepo(db),
		IdGen:         idGen,
	})
	return BasicUserSVC
}
