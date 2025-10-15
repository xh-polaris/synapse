package model

import (
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"gorm.io/gorm"
)

type ThirdPartyUser struct {
	ID        id.ID          `gorm:"column:id;primaryKey;autoIncrement:false;type:binary(12)" json:"id"` // Primary Key ID
	Code      string         `gorm:"column:code;uniqueIndex;type:varchar(256)" json:"code"`              // Code 唯一标识
	App       string         `gorm:"column:app;type:varchar(12);" json:"app"`                            // App 第三方应用标识
	CreatedAt int64          `gorm:"column:created_at;not null;autoCreateTime:milli;" json:"created_at"` // Creation Time (Milliseconds)
	UpdatedAt int64          `gorm:"column:updated_at;not null;autoUpdateTime:milli;" json:"updated_at"` // Update Time (Milliseconds)
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at"`                               // Deletion Time (Milliseconds)
}
