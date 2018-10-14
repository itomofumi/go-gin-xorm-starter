package util

import "time"

// GetTimeNowFunc is current time function.
// Default function is "time.Now".
var GetTimeNowFunc = time.Now

// GetTimeNow は現在時刻を返す
func GetTimeNow() time.Time {
	now := GetTimeNowFunc()
	return now
}

// GetFormatedTimeNow はDB用フォーマット済みの現在時刻を返す
func GetFormatedTimeNow() string {
	return GetTimeNow().Format("2006-01-02 15:04:05")
}
