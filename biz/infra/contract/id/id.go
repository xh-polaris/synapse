package id

import (
	"context"
	"database/sql/driver"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ID 是mongodb中的ObjectID, 为了和其他服务使用MongoDB数据库相适应
type ID primitive.ObjectID

type IDGenerator interface {
	GenID(ctx context.Context) ID
	GenMultiIDs(ctx context.Context, counts int) []ID
}

// Hex 方法返回 ObjectID 的十六进制字符串表示
func (i ID) Hex() string {
	return primitive.ObjectID(i).Hex()
}

func FromHex(str string) (ID, error) {
	oid, err := primitive.ObjectIDFromHex(str)
	return ID(oid), err
}

// Value 实现 driver.Valuer 接口，用于将 ID 转换为数据库可存储的值
func (i ID) Value() (driver.Value, error) {
	// 将 ObjectID 转换为 12 字节的二进制数据
	return i[:], nil
}

// Scan 实现 sql.Scanner 接口，用于从数据库读取值到 ID
func (i *ID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("ID cannot be nil")
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// 确保字节长度正确 (ObjectID 是 12 字节)
	if len(b) != 12 {
		return errors.New("invalid ObjectID length")
	}

	// 直接转换为 ID 类型
	*i = ID(primitive.ObjectID(b))
	return nil
}
