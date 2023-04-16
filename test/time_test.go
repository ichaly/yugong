package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	duration := 2*time.Hour + 1*time.Minute + 1*time.Second +
		1*time.Millisecond + 1*time.Microsecond + 1*time.Nanosecond
	fmt.Printf("duration:%v\n", duration.Minutes())
	fmt.Printf("%.0f时%.0f分%.0f秒%.3d毫秒",
		duration.Hours(),
		float64(int64(duration.Minutes())%60),
		float64(int64(duration.Seconds())%60),
		duration.Milliseconds()%1000,
	)
}
