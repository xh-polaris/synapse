package model

import "gorm.io/gorm"

type Auth struct {
	AuthId      string         `gorm:"column:auth_id;primaryKey" json:"auth_id"`                           // 认证id
	AuthType    string         `gorm:"column:auth_type;primaryKey;" json:"auth_type"`                      // 认证类型
	BasicUserID int64          `gorm:"column:basic_user_id;" json:"basic_user_id"`                         // 用户id
	CreatedAt   int64          `gorm:"column:created_at;not null;autoCreateTime:milli;" json:"created_at"` // Creation Time (Milliseconds)
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at"`                               // Deletion Time (Milliseconds)
}
