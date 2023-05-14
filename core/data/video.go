package data

import "time"

type Video struct {
	Id        int64      `gorm:"AUTO_INCREMENT;comment:ID"`
	Vid       string     `gorm:"size:50;uniqueIndex;comment:视频ID"`
	Fid       string     `gorm:"size:100;comment:外部账号"`
	Aid       string     `gorm:"size:100;comment:虚拟账号"`
	Url       string     `gorm:"size:500;comment:视频链接"`
	From      Platform   `gorm:"size:50;comment:来源"`
	Title     string     `gorm:"size:200;comment:视频标题"`
	Cover     string     `gorm:"size:500;comment:封面"`
	State     int        `gorm:"comment:状态"`
	Sticky    bool       `gorm:"comment:是否置顶"`
	Remark    string     `gorm:"type:text;comment:备注"`
	UploadAt  *time.Time `gorm:"comment:上传时间"`
	SourceAt  *time.Time `gorm:"comment:原始时间"`
	CreatedAt time.Time  `gorm:"comment:创建时间"`
	UpdatedAt time.Time  `gorm:"comment:更新时间"`
}

func (Video) TableName() string {
	return "video"
}
