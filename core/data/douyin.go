package data

import "time"

type Douyin struct {
	Id        uint       `gorm:"AUTO_INCREMENT;comment:ID"`
	From      Platform   `gorm:"size:50;comment:来源"`
	OpenId    string     `gorm:"size:100;comment:抖音OpenId"`
	Aid       string     `gorm:"size:100;comment:虚拟账号"`
	Fid       string     `gorm:"size:100;comment:抖音账号"`
	Url       string     `gorm:"size:200;comment:分享链接"`
	Avatar    string     `gorm:"size:200;comment:头像"`
	Nickname  string     `gorm:"size:50;comment:昵称"`
	ItemCount int64      `gorm:"size:50;comment:作品数"`
	CreatedAt time.Time  `gorm:"comment:创建时间"`
	UpdatedAt time.Time  `gorm:"comment:更新时间"`
	LastTime  *time.Time `gorm:"comment:最后时间"`
	FirstTime *time.Time `gorm:"comment:开始时间"`
}

func (Douyin) TableName() string {
	return "douyin"
}
