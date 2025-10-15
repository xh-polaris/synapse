package dal

import (
	"context"

	"github.com/xh-polaris/synapse/biz/domain/thirdparty/dal/model"
	"github.com/xh-polaris/synapse/biz/domain/thirdparty/dal/query"
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"github.com/xh-polaris/synapse/biz/infra/contract/orm"
	"gorm.io/gen/field"
)

func NewThirdPartyUserDAO(db *orm.DB) *ThirdPartyUserDAO {
	return &ThirdPartyUserDAO{query: query.Use(db)}
}

type ThirdPartyUserDAO struct {
	query *query.Query
}

func (t *ThirdPartyUserDAO) FindOrCreate(ctx context.Context, id id.ID, code, app string) (*model.ThirdPartyUser, error) {
	u, err := t.query.WithContext(ctx).ThirdPartyUser.Attrs(field.Attrs(&model.ThirdPartyUser{ID: id, App: app})).
		Where(t.query.ThirdPartyUser.Code.Eq(code)).FirstOrCreate()
	return u, err
}
