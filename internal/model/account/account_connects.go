package account

import (
	"time"

	"github.com/ogreks/meeseeks-box/internal/model"
)

type AccountConnectModel struct {
	model.Model

	AccountID            int64     `gorm:"column:account_id;not null;index:idx_account_id;type:bigint(20);default:0;comment:账号主键ID"`
	ConnectPlatformID    uint8     `gorm:"column:connect_platform_id;not null;index:idx_connect_platform_id;type:tinyint(4);default:0;comment:连接平台ID"`
	ConnectAccountID     string    `gorm:"column:connect_account_id;not null;index:idx_connect_account_id;type:varchar(128);default:null;comment:关联账号ID唯一值例如：微信openid或GitHubID"`
	ConnectToken         string    `grom:"column:connect_token;type:varchar(128);default:null;comment:授权Token"`
	ConnectRefreshToken  string    `grom:"column:connect_refresh_token;type:varchar(128);default:null;comment:授权RefreshToken"`
	ConnectUserName      string    `gorm:"column:connect_user_name;type:varchar(128);default:null;comment:关联用户名"`
	ConnectNickname      string    `gorm:"column:connect_nickname;type:varchar(128);default:null;comment:关联昵称"`
	IsEnabled            uint8     `gorm:"column:is_enabled;type:tinyint(1);default:1;comment:是否有效0.无效/1.有效"`
	MoreJson             string    `gorm:"column:more_json;type:json;default:null;comment:扩展字段"`
	RefreshTokenExpireAt time.Time `gorm:"refresh_token_expire_at;type:datetime;default:null;comment:授权RefreshToken有效期"`
}

func (a *AccountConnectModel) TableName() string {
	return "account_connects"
}
