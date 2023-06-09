package data

import "time"

type Author struct {
	Id        int64      `gorm:"AUTO_INCREMENT;comment:ID"`
	From      Platform   `gorm:"size:50;comment:来源"`
	OpenId    string     `gorm:"size:100;comment:外部ID"`
	Aid       string     `gorm:"size:100;unique_index;comment:虚拟账号"`
	Url       string     `gorm:"size:200;comment:分享链接"`
	Avatar    string     `gorm:"size:200;comment:头像"`
	Nickname  string     `gorm:"size:50;comment:昵称"`
	Signature string     `gorm:"size:200;comment:个性签名"`
	Cron      string     `gorm:"size:50;comment:定时任务"`
	Total     int        `gorm:"comment:初始数量"`
	Start     *time.Time `gorm:"comment:初始时间"`
	CreatedAt time.Time  `gorm:"comment:创建时间"`
	UpdatedAt time.Time  `gorm:"comment:更新时间"`
	Disable   bool       `gorm:"comment:是否禁用"`
}

func (Author) TableName() string {
	return "author"
}
