package utils

import "time"

const timeFmt = "2006-01-02 15:04:05"

func TimestampToStr(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(timeFmt)
}
