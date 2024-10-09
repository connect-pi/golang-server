package clog

import (
	"fmt"
	"sync"
)

var (
	Logs               []string   // برش برای ذخیره لاگ‌ها
	mu                 sync.Mutex // قفل برای اطمینان از ایمنی همزمانی
	maxLogs            = 300      // حداکثر تعداد لاگ‌ها
	showLogsInTerminal = true     // آیا لاگ‌ها در ترمینال نشان داده شوند
)

// Custom Println that writes to the log slice and optionally to the terminal
func Println(a ...interface{}) {
	// Lock to ensure safe concurrent access
	mu.Lock()
	defer mu.Unlock()

	// Write to the log message
	logMessage := fmt.Sprintln(a...)
	Logs = append(Logs, logMessage) // ذخیره لاگ در آرایه

	// محدود کردن تعداد لاگ‌ها به ۳۰۰
	if len(Logs) > maxLogs {
		Logs = Logs[1:] // حذف قدیمی‌ترین لاگ
	}

	if showLogsInTerminal {
		go fmt.Print(logMessage) // ارسال به ترمینال
	}
}
