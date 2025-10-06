package dal

import (
	"context"
	"errors"

	"github.com/xh-polaris/synapse/biz/domain/basicuser/dal/model"
	"github.com/xh-polaris/synapse/biz/domain/basicuser/dal/query"
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"github.com/xh-polaris/synapse/biz/infra/contract/orm"
	"gorm.io/gorm"
)

func NewBasicUserDAO(db *orm.DB) *BasicUserDAO {
	return &BasicUserDAO{query: query.Use(db)}
}

type BasicUserDAO struct {
	query *query.Query
}

func (d *BasicUserDAO) FindByPhone(ctx context.Context, phone string) (*model.BasicUser, error) {
	user, err := d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.Phone.Eq(phone)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, err
}

func (d *BasicUserDAO) Create(ctx context.Context, nu *model.BasicUser) (*model.BasicUser, error) {
	err := d.query.WithContext(ctx).BasicUser.Create(nu)
	return nu, err
}

func (d *BasicUserDAO) ResetPassword(ctx context.Context, basicUserId, password string) error {
	buid, err := id.FromHex(basicUserId)
	if err != nil {
		return err
	}
	_, err = d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.ID.Eq(buid)).
		Update(d.query.BasicUser.Password, password)
	return err
}
