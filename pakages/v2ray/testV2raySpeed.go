package v2ray

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

// تست سرعت اینترنت از طریق پروکسی SOCKS5 با محدودیت زمانی
func TestV2raySpeed(proxyPort int) (float64, error) {
	fmt.Println("Connect to: SOCKS5", "127.0.0.1:"+strconv.Itoa(proxyPort))

	// تنظیم پروکسی SOCKS5
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+strconv.Itoa(proxyPort), nil, proxy.Direct)
	if err != nil {
		return 0, fmt.Errorf("failed to create SOCKS5 dialer: %v", err)
	}

	// تنظیم کلاینت HTTP با پروکسی SOCKS5
	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
	}

	// ایجاد کانتکست با محدودیت زمانی
	ctx, cancel := context.WithTimeout(context.Background(), 1300*time.Millisecond)
	defer cancel()

	// ایجاد درخواست با کانتکست
	req, err := http.NewRequestWithContext(ctx, "GET", "https://raw.githubusercontent.com/BitDoctor/speed-test-file/master/5mb.txt", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// ارسال درخواست HTTP برای تست سرعت
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// خواندن داده‌ها در مدت زمان ۱ ثانیه
	start := time.Now()
	var totalBytesRead int64
	buf := make([]byte, 8192) // بافر 8 کیلوبایتی

	for {
		select {
		case <-ctx.Done(): // اگر زمان تمام شد
			duration := time.Since(start)
			speedMbps := float64(totalBytesRead) / duration.Seconds() / (1024 * 1024) // سرعت بر حسب مگابیت بر ثانیه
			return speedMbps, nil
		default:
			n, err := resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				return 0, fmt.Errorf("failed to read response body: %v", err)
			}
			totalBytesRead += int64(n)
			if err == io.EOF {
				break
			}
		}
	}

	// در صورتی که حلقه به پایان برسد
	duration := time.Since(start)
	speedMbps := float64(totalBytesRead) / duration.Seconds() / (1024 * 1024) // سرعت بر حسب مگابیت بر ثانیه

	return speedMbps, nil
}
