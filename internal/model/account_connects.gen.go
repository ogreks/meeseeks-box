// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameAccountConnect = "account_connects"

// AccountConnect mapped from table <account_connects>
type AccountConnect struct {
	ID                   uint64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true" json:"id"`
	AccountID            uint64     `gorm:"column:account_id;type:bigint unsigned;not null;index:idx_account_id,priority:1;comment:账号主键ID" json:"account_id"`                                    // 账号主键ID
	ConnectPlatformID    uint32     `gorm:"column:connect_platform_id;type:tinyint unsigned;not null;index:idx_connect_platform_id,priority:1;comment:连接平台ID" json:"connect_platform_id"`        // 连接平台ID
	ConnectAccountID     string     `gorm:"column:connect_account_id;type:varchar(128);index:idx_connect_account_id,priority:1;comment:关联账号ID唯一值例如：微信openid或GitHubID" json:"connect_account_id"` // 关联账号ID唯一值例如：微信openid或GitHubID
	ConnectToken         string     `gorm:"column:connect_token;type:longtext" json:"connect_token"`
	ConnectRefreshToken  string     `gorm:"column:connect_refresh_token;type:longtext" json:"connect_refresh_token"`
	ConnectUserName      string     `gorm:"column:connect_user_name;type:varchar(128);comment:关联用户名" json:"connect_user_name"`                                                                   // 关联用户名
	ConnectNickName      string     `gorm:"column:connect_nick_name;type:varchar(128);comment:关联昵称" json:"connect_nick_name"`                                                                    // 关联昵称
	IsEnabled            *uint32    `gorm:"column:is_enabled;type:tinyint unsigned;not null;default:1;comment:是否有效0.无效/1.有效" json:"is_enabled"`                                                  // 是否有效0.无效/1.有效
	MoreJSON             string     `gorm:"column:more_json;type:json;comment:扩展字段" json:"more_json"`                                                                                            // 扩展字段
	RefreshTokenExpireAt *time.Time `gorm:"column:refresh_token_expire_at;type:timestamp;index:idx_refresh_token_expire_at,priority:1;comment:授权RefreshToken有效期" json:"refresh_token_expire_at"` // 授权RefreshToken有效期
	CreatedAt            *time.Time `gorm:"column:created_at;type:timestamp;not null;index:idx_created_at,priority:1;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                  // 创建时间
	UpdatedAt            *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                                  // 更新时间
	DeletedAt            *time.Time `gorm:"column:deleted_at;type:timestamp;index:idx_account_connects_deleted_at,priority:1;comment:删除时间" json:"deleted_at"`                                    // 删除时间
}

// TableName AccountConnect's table name
func (*AccountConnect) TableName() string {
	return TableNameAccountConnect
}
