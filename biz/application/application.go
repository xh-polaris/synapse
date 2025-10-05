package application

import (
	"context"

	"github.com/xh-polaris/synapse/biz/application/base/appinfra"
	"github.com/xh-polaris/synapse/biz/application/basicuser"
	"github.com/xh-polaris/synapse/biz/application/system"
)

type BasicService struct {
	infra        *appinfra.AppDependencies
	systemSVC    *system.SystemService
	basicUserSVC *basicuser.BasicUserService
}

func InitApplication(ctx context.Context) error {
	infra, err := appinfra.Init(ctx)
	if err != nil {
		return err
	}
	initBasicServices(ctx, infra)
	return nil
}

func initBasicServices(ctx context.Context, infra *appinfra.AppDependencies) *BasicService {
	systemSVC := system.InitService(ctx, infra.SMS)
	basicUserSVC := basicuser.InitService(ctx, infra.SMS, infra.DB, infra.IDGen)
	return &BasicService{infra: infra, systemSVC: systemSVC, basicUserSVC: basicUserSVC}
}
