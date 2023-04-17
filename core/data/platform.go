package data

import "database/sql/driver"

type Platform string

const (
	DouYin      Platform = "DOUYIN"
	XiaoHongShu Platform = "XIAOHONGSHU"
)

func (my *Platform) Scan(value interface{}) error {
	*my = Platform(value.(string))
	return nil
}

func (my Platform) Value() (driver.Value, error) {
	return string(my), nil
}
