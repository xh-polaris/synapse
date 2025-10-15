package internal

import (
	"github.com/xh-polaris/synapse/biz/api/model/thirdparty"
	"github.com/xh-polaris/synapse/biz/domain/thirdparty/entity"
)

func ThirdPartyUser2BasicUserVO(eu *entity.ThirdPartyUser) *thirdparty.ThirdPartyBasicUser {
	return &thirdparty.ThirdPartyBasicUser{
		BasicUserId: eu.ID,
	}
}
