package util

import (
	"strings"
)

func Join(elem ...string) string {
	b := strings.Builder{}
	for _, e := range elem {
		b.WriteString(e)
	}
	return b.String()
}
