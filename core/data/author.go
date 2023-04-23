package data

import "time"

type Author struct {
	Id        uint       `gorm:"AUTO_INCREMENT;comment:ID"`
	From      Platform   `gorm:"size:50;comment:来源"`
	OpenId    string     `gorm:"size:100;comment:抖音OpenId"`
	Aid       string     `gorm:"size:100;comment:虚拟账号"`
	Fid       string     `gorm:"size:100;comment:抖音账号"`
	Url       string     `gorm:"size:200;comment:分享链接"`
	Avatar    string     `gorm:"size:200;comment:头像"`
	Nickname  string     `gorm:"size:50;comment:昵称"`
	Total     int64      `gorm:"size:50;comment:作品数"`
	Cron      string     `gorm:"size:50;comment:定时任务"`
	CreatedAt time.Time  `gorm:"comment:创建时间"`
	UpdatedAt time.Time  `gorm:"comment:更新时间"`
	MaxTime   *time.Time `gorm:"comment:最后同步时间"`
	MinTime   *time.Time `gorm:"comment:起始同步时间"`
}

func (Author) TableName() string {
	return "author"
}
