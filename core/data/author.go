package data

import "time"

type Author struct {
	Id        int64      `gorm:"AUTO_INCREMENT;comment:ID"`
	From      Platform   `gorm:"size:50;comment:来源"`
	OpenId    string     `gorm:"size:100;comment:抖音OpenId"`
	Aid       string     `gorm:"size:100;comment:虚拟账号"`
	Fid       string     `gorm:"size:100;comment:抖音账号"`
	Url       string     `gorm:"size:200;comment:分享链接"`
	Avatar    string     `gorm:"size:200;comment:头像"`
	Nickname  string     `gorm:"size:50;comment:昵称"`
	Signature string     `gorm:"size:200;comment:个性签名"`
	Total     int64      `gorm:"size:50;comment:初始数量"`
	Cron      string     `gorm:"size:50;comment:定时任务"`
	MaxTime   *time.Time `gorm:"comment:最后同步时间"`
	MinTime   *time.Time `gorm:"comment:起始同步时间"`
	CreatedAt time.Time  `gorm:"comment:创建时间"`
	UpdatedAt time.Time  `gorm:"comment:更新时间"`
}

func (Author) TableName() string {
	return "author"
}
