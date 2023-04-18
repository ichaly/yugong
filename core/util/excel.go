package util

import "github.com/xuri/excelize/v2"

func ReadExcel(name string, sheets ...string) (map[string][][]string, error) {
	f, err := excelize.OpenFile(name) //
	if err != nil {
		return nil, err
	}
	res := make(map[string][][]string)
	for _, s := range sheets {
		rows, err := f.GetRows(s)
		if err != nil {
			continue
		}
		res[s] = rows
	}
	return res, err
}
