package base

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	OK        = newResult(http.StatusOK)                  // 通用成功
	FORBIDDEN = newResult(http.StatusForbidden)           // 无权限
	ERROR     = newResult(http.StatusInternalServerError) // 通用错误
)

type result struct {
	Code   int                      `json:"code"`             // 错误码
	Data   interface{}              `json:"data,omitempty"`   // 返回数据
	Errors []map[string]interface{} `json:"errors,omitempty"` // 错误信息
}

// WithError 自定义错误信息
func (res *result) WithError(errors ...error) result {
	var maps []map[string]interface{}
	if errors != nil && len(errors) > 0 {
		for _, e := range errors {
			maps = append(maps, map[string]interface{}{"message": e.Error()})
		}
	}
	return result{
		Code:   res.Code,
		Errors: maps,
	}
}

func (res *result) WithMessage(msgs ...string) result {
	var errs []error
	for _, m := range msgs {
		errs = append(errs, errors.New(m))
	}
	return res.WithError(errs...)
}

// WithData 追加响应数据
func (res *result) WithData(data interface{}) result {
	return result{
		Code: res.Code,
		Data: data,
	}
}

// ToString 返回 JSON 格式的错误详情
func (res *result) ToString() string {
	err := &struct {
		Code   int                      `json:"code"`
		Data   interface{}              `json:"data,omitempty"`
		Errors []map[string]interface{} `json:"errors,omitempty"`
	}{
		Code:   res.Code,
		Data:   res.Data,
		Errors: res.Errors,
	}
	raw, _ := json.Marshal(err)
	return string(raw)
}

// newResult 构造函数
func newResult(code int) *result {
	return &result{
		Code: code,
		Data: nil,
	}
}
