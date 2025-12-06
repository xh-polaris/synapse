package model

import (
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type School struct {
	ID        id.ID          `gorm:"column:id;primaryKey;autoIncrement:false;type:binary(12)" json:"id"` // Primary Key ID
	Name      string         `gorm:"column:name;type:varchar(60);comment:School Name" json:"name"`       // Name 学校名
	Extra     datatypes.JSON `gorm:"column:extra" json:"extra"`                                          // Extra json字符串存储可能存在的额外信息
	CreatedAt int64          `gorm:"column:created_at;not null;autoCreateTime:milli;" json:"created_at"` // Creation Time (Milliseconds)
	UpdatedAt int64          `gorm:"column:updated_at;not null;autoUpdateTime:milli;" json:"updated_at"` // Update Time (Milliseconds)
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at"`                               // Deletion Time (Milliseconds)
}
