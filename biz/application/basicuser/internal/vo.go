package internal

import (
	"github.com/xh-polaris/synapse/biz/api/model/basicuser"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/entity"
)

func BasicUserPO2VO(u *entity.BasicUser) *basicuser.BasicUser {
	return &basicuser.BasicUser{
		BasicUserId: u.ID,
	}
}
