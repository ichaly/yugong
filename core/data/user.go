package data

import "time"

type User struct {
	Id        uint      `gorm:"AUTO_INCREMENT;comment:ID"`
	Aid       string    `gorm:"size:100;comment:阿里号"`
	Did       string    `gorm:"size:100;comment:抖音号"`
	Url       string    `gorm:"size:50;unique_index;comment:分享链接"`
	Avatar    string    `gorm:"size:150;comment:头像"`
	Nickname  string    `gorm:"size:50;comment:昵称"`
	ItemCount int64     `gorm:"size:50;comment:作品数"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
}
