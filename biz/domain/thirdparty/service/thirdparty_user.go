package service

import (
	"context"

	"github.com/xh-polaris/synapse/biz/domain/thirdparty/entity"
)

type ThirdParty interface {
	BUPTLogin(ctx context.Context, ticket string) (*entity.ThirdPartyUser, error)
}
