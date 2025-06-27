package model

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
	"gorm.io/gorm/schema"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/crypto"
)

// AESEncryptString 为 Gorm 自定义字段类型，用 AES-GCM 算法加密
type AESEncryptString string

// Scan 解析 driver 提供的数据
func (s *AESEncryptString) Scan(value any) error {
	if value == nil {
		*s = ""
		return nil
	}

	var data string
	switch v := value.(type) {
	case string:
		data = v
	case []byte:
		data = string(v)
	default:
		return errors.New("invalid value type")
	}

	decryptedData, err := crypto.AESDecrypt(config.G.Service.EncryptSecret, data)
	if err != nil {
		return err
	}
	*s = AESEncryptString(decryptedData)
	return nil
}

// Value 提供 driver.Value
func (s AESEncryptString) Value() (driver.Value, error) {
	encryptedData, err := crypto.AESEncrypt(config.G.Service.EncryptSecret, string(s))
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

// GormDataType 提供 Gorm 需要的数据类型
func (s *AESEncryptString) GormDataType() string {
	return "varchar(128)"
}

var (
	_ driver.Valuer                = (*AESEncryptString)(nil)
	_ sql.Scanner                  = (*AESEncryptString)(nil)
	_ schema.GormDataTypeInterface = (*AESEncryptString)(nil)
)
