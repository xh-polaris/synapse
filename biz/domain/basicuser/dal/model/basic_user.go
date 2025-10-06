package model

import (
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BasicUser struct {
	ID        id.ID          `gorm:"column:id;primaryKey;autoIncrement:false;type:binary(12)" json:"id"`            // Primary Key ID
	Code      *string        `gorm:"column:code;uniqueIndex;type:varchar(24)" json:"code"`                          // Code 身份证(Encrypted)
	Phone     *string        `gorm:"column:phone;uniqueIndex;type:varchar(16)" json:"phone"`                        // Phone 手机号
	Password  *string        `gorm:"column:password;type:varchar(60);comment:Password (Encrypted)" json:"password"` // Password (Encrypted)
	Name      string         `gorm:"column:name;comment:User Nickname" json:"name"`                                 // User Nickname
	Gender    uint8          `gorm:"column:gender" json:"gender"`                                                   // Gender 性别
	Extra     datatypes.JSON `gorm:"column:extra" json:"extra"`                                                     // Extra json字符串存储可能存在的额外信息
	CreatedAt int64          `gorm:"column:created_at;not null;autoCreateTime:milli;" json:"created_at"`            // Creation Time (Milliseconds)
	UpdatedAt int64          `gorm:"column:updated_at;not null;autoUpdateTime:milli;" json:"updated_at"`            // Update Time (Milliseconds)
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at"`                                          // Deletion Time (Milliseconds)
}
