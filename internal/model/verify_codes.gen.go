// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameVerifyCode = "verify_codes"

// VerifyCode mapped from table <verify_codes>
type VerifyCode struct {
	ID         uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true" json:"id"`
	TemplateID uint64         `gorm:"column:template_id;type:bigint unsigned;not null;comment:模板编号" json:"template_id"`                                             // 模板编号
	Type       *uint64        `gorm:"column:type;type:bigint unsigned;not null;default:1;comment:类型:1.邮件 / 2.短信" json:"type"`                                       // 类型:1.邮件 / 2.短信
	Account    *string        `gorm:"column:account;type:varchar(128);not null;index:idx_account,priority:1;default:0;comment:邮箱或手机，手机号带国际区号，区号无加号" json:"account"` // 邮箱或手机，手机号带国际区号，区号无加号
	Code       *string        `gorm:"column:code;type:char(8);not null;default:0;comment:验证码" json:"code"`                                                          // 验证码
	IsEnabled  *uint32        `gorm:"column:is_enabled;type:tinyint unsigned;not null;default:1;comment:是否有效" json:"is_enabled"`                                    // 是否有效
	ExpiredAt  time.Time      `gorm:"column:expired_at;type:datetime;not null;comment:失效时间" json:"expired_at"`                                                      // 失效时间
	CreatedAt  *time.Time     `gorm:"column:created_at;type:datetime;not null;index:idx_created_at,priority:1;default:now();comment:创建时间" json:"created_at"`        // 创建时间
	UpdatedAt  *time.Time     `gorm:"column:updated_at;type:datetime;not null;default:now();comment:更新时间" json:"updated_at"`                                        // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deleted_at"`                                                               // 删除时间
}

// TableName VerifyCode's table name
func (*VerifyCode) TableName() string {
	return TableNameVerifyCode
}