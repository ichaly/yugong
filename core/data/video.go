package data

import "time"

type Video struct {
	Id        uint       `gorm:"AUTO_INCREMENT;comment:ID"`
	Title     string     `gorm:"size:200;comment:视频标题"`
	Cover     string     `gorm:"size:500;comment:封面"`
	From      string     `gorm:"size:50;comment:来源"`
	Did       string     `gorm:"size:100;comment:抖音号"`
	Aid       string     `gorm:"size:100;comment:阿里号"`
	Url       string     `gorm:"size:500;unique_index;comment:视频链接"`
	Target    string     `gorm:"size:200;comment:产物链接"`
	UploadAt  *time.Time `gorm:"comment:上传时间"`
	SourceAt  time.Time  `gorm:"comment:原始时间"`
	CreatedAt time.Time  `gorm:"comment:创建时间"`
	UpdatedAt time.Time  `gorm:"comment:更新时间"`
}
