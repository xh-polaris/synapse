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

func NewSchoolDAO(db *orm.DB) *SchoolDAO {
	return &SchoolDAO{
		query: query.Use(db),
	}
}

type SchoolDAO struct {
	query *query.Query
}

func (d *SchoolDAO) FindByID(ctx context.Context, idStr string) (*model.School, error) {
	sid, err := id.FromHex(idStr)
	if err != nil {
		return nil, err
	}
	school, err := d.query.WithContext(ctx).School.Where(d.query.School.ID.Eq(sid)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return school, err
}
