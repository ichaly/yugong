package util

import (
	"strconv"
	"strings"
)

func JoinString(elem ...string) string {
	b := strings.Builder{}
	for _, e := range elem {
		b.WriteString(e)
	}
	return b.String()
}

func FormatLong(v int64) string {
	return strconv.FormatInt(v, 10)
}
