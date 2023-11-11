package account

import (
	"time"

	"github.com/ogreks/meeseeks-box/internal/model"
)

type AccountModel struct {
	model.Model

	AID          string    `gorm:"column:aid;not null;uniqueIndex:uidx_aid;type:varchar(32);default:0;comment:用户唯一ID"`
	Type         uint8     `gorm:"column:type;not null;type:tinyint(4);default:0;comment:用户类型1.超级管理员/2.普通管理员/3.普通用户"`
	CountryCode  string    `gorm:"column:country_code;type:varchar(8);default:86;comment:国家区号例如中国86/美国则是1"`
	PurePhone    string    `gorm:"column:pure_phone;type:varchar(128);default:null;comment:不带区号的手机号"`
	Phone        string    `gorm:"column:phone;type:varchar(128);default:null;comment:带区号的手机号"`
	Email        string    `gorm:"column:email;type:varchar(128);default:null;comment:邮箱"`
	Password     string    `gorm:"column:password;type:varchar(128);default:null;comment:密码"`
	LastLoginAt  time.Time `gorm:"column:last_login_at;type:datetime;default:null;comment:最后登录时间"`
	IsEnabled    uint8     `gorm:"column:is_enabled;type:tinyint(1);default:1;comment:是否启用0.禁用/1.正常"`
	WaitDelete   uint8     `gorm:"column:wait_delete;type:tinyint(1);default:0;comment:等待删除0.否/1.是/2.永久"`
	WaitDeleteAt time.Time `gorm:"column:wait_delete_at;type:datetime;default:null;comment:等待删除时间"`
}
