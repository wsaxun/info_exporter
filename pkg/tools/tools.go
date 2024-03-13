package tools

import (
	"time"
)

// CronJob 每隔一段时间执行一次
func CronJob(job func(), duration time.Duration) {
	go func() {
		// 定时前运行一次
		job()
		ticker := time.NewTicker(duration)
		for range ticker.C {
			job()
		}
	}()
}
