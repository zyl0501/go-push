package utils

import "time"

const(
	FullTimeFormat = "2006-01-02 15:04:05"
)

func MillisecondToDuration(millisecond int64) time.Duration {
	return time.Duration(millisecond * time.Millisecond.Nanoseconds())
}

func DurationToMillisecond(duration time.Duration) int64 {
	return int64(duration / time.Millisecond)
}
