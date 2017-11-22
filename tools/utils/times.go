package utils

import "time"

func MillisecondToDuration(millisecond int64) time.Duration {
	return time.Duration(millisecond * time.Millisecond.Nanoseconds())
}

func DurationToMillisecond(duration time.Duration) int64 {
	return int64(duration / time.Millisecond)
}
