package util

import "time"

// GetTimeNow は現在時刻を返す
// TODO: テストに使う時にはビルドタグとかで現在時刻を操作できるようする
func GetTimeNow() *time.Time {
	now := time.Now()
	return &now
}

// GetFormatedTimeNow はDB用フォーマット済みの現在時刻を返す
func GetFormatedTimeNow() string {
	return GetTimeNow().Format("2006-01-02 15:04:05")
}
