package user

import (
	"time"

	"github.com/ogreks/meeseeks-box/internal/model"
)

type UserModel struct {
	model.Model
	AccountID     uint64    `gorm:"column:account_id;type:bigint(20);not null;uniqueIndex;comment:用户所属账户"`
	UID           string    `gorm:"column:id;type:bigint(20);not null;uniqueIndex;comment:用户唯一编号"`
	UserName      string    `gorm:"column:user_name;type:varchar(64);uniqueIndex;comment:用户名"`
	NickName      string    `gorm:"column:nick_name;type:varchar(64);comment:昵称"`
	Password      string    `gorm:"column:password;type:char(64);comment:密码"`
	Gender        uint8     `gorm:"column:gender;type:tinyint(1);default:1;comment:性别"`
	Bio           string    `gorm:"column:bio;type:text;comment:个人简介"`
	LastAtivityAt time.Time `gorm:"column:last_activity_at;type:datetime;default:null;comment:最后活跃时间"`
	IsEnabled     uint8     `gorm:"column:is_enabled;type:tinyint(1);default:1;comment:是否启用0.封号/1.正常"`
	WaitDelete    uint8     `gorm:"column:wait_delete;type:tinyint(1);default:0;comment:是否等待删除0.否/1.是/2.永久"`
	WaitDeleteAt  time.Time `gorm:"column:wait_delete_at;type:datetime;default:null;comment:等待删除时间"`
}
