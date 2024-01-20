// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameAccount = "accounts"

// Account mapped from table <accounts>
type Account struct {
	ID           uint64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true" json:"id"`
	Aid          string     `gorm:"column:aid;type:varchar(32);not null;uniqueIndex:uidx_aid,priority:1;comment:用户唯一ID" json:"aid"`                // 用户唯一ID
	Type         uint32     `gorm:"column:type;type:tinyint unsigned;not null;comment:用户类型1.超级管理员/2.普通管理员/3.普通用户" json:"type"`                     // 用户类型1.超级管理员/2.普通管理员/3.普通用户
	CountryCode  *string    `gorm:"column:country_code;type:varchar(8);default:86;comment:国家区号例如中国86/美国则是1" json:"country_code"`                   // 国家区号例如中国86/美国则是1
	PurePhone    string     `gorm:"column:pure_phone;type:varchar(128);index:idx_pure_phone,priority:1;comment:不带区号的手机号" json:"pure_phone"`        // 不带区号的手机号
	Phone        string     `gorm:"column:phone;type:varchar(128);index:idx_phone,priority:1;comment:带区号的手机号" json:"phone"`                        // 带区号的手机号
	Email        string     `gorm:"column:email;type:varchar(128);index:idx_email,priority:1;comment:邮箱" json:"email"`                             // 邮箱
	UserName     string     `gorm:"column:user_name;type:varchar(64);index:idx_user_name,priority:1;comment:用户账号(非账号注册随机生成)" json:"user_name"`     // 用户账号(非账号注册随机生成)
	Password     string     `gorm:"column:password;type:varchar(128);comment:密码" json:"password"`                                                  // 密码
	LastLoginAt  *time.Time `gorm:"column:last_login_at;type:timestamp;index:idx_last_login_at,priority:1;comment:最后登录时间" json:"last_login_at"`    // 最后登录时间
	IsEnabled    *uint32    `gorm:"column:is_enabled;type:tinyint unsigned;not null;default:1;comment:是否启用0.禁用/1.正常" json:"is_enabled"`            // 是否启用0.禁用/1.正常
	WaitDelete   uint32     `gorm:"column:wait_delete;type:tinyint unsigned;not null;comment:等待删除0.否/1.是/2.永久" json:"wait_delete"`                 // 等待删除0.否/1.是/2.永久
	WaitDeleteAt *time.Time `gorm:"column:wait_delete_at;type:timestamp;index:idx_wait_delete_at,priority:1;comment:等待删除时间" json:"wait_delete_at"` // 等待删除时间
	CreatedAt    *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`            // 创建时间
	UpdatedAt    *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`            // 更新时间
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`                                               // 删除时间
}

// TableName Account's table name
func (*Account) TableName() string {
	return TableNameAccount
}
