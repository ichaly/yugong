package data

import "time"

type Video struct {
	Id        int64      `gorm:"AUTO_INCREMENT;comment:ID"`
	From      Platform   `gorm:"size:50;comment:来源"`
	Vid       string     `gorm:"size:50;unique_index;comment:视频ID"`
	Url       string     `gorm:"size:500;comment:视频链接"`
	Width     int64      `gorm:"comment:宽"`
	Height    int64      `gorm:"comment:高"`
	Title     string     `gorm:"size:200;comment:视频标题"`
	Cover     string     `gorm:"size:500;comment:封面"`
	Fid       string     `gorm:"size:100;comment:外部账号"`
	Aid       string     `gorm:"size:100;comment:虚拟账号"`
	State     int        `gorm:"comment:状态"`
	UploadAt  *time.Time `gorm:"comment:上传时间"`
	SourceAt  time.Time  `gorm:"comment:原始时间"`
	CreatedAt time.Time  `gorm:"comment:创建时间"`
	UpdatedAt time.Time  `gorm:"comment:更新时间"`
}

func (Video) TableName() string {
	return "video"
}
