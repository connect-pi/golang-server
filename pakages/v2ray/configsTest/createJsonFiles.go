package configsTest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"project/pakages/v2ray"
)

func CreateJsonFiles() {
	// آرایه‌ای از رشته‌ها به عنوان ورودی
	data := v2ray.Uris

	// مسیر پوشه‌ای که فایل‌ها در آن قرار می‌گیرند
	dirPath := ".v2rayConfig/testConfigs"

	// پاک کردن پوشه‌ی testing اگر وجود دارد
	if err := os.RemoveAll(dirPath); err != nil {
		fmt.Println("Error removing directory:", err)
		return
	}

	// ایجاد مجدد پوشه
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// برای هر آیتم در آرایه، یک پوشه و فایل JSON ایجاد کنید
	for i, item := range data {
		// ایجاد پوشه با نام ایندکس
		subDirPath := filepath.Join(dirPath, fmt.Sprintf("%d", i))
		if err := os.MkdirAll(subDirPath, os.ModePerm); err != nil {
			fmt.Println("Error creating subdirectory:", err)
			continue
		}

		// نام فایل JSON
		fileName := filepath.Join(subDirPath, "config.json")

		// محتوای JSON را ایجاد کنید
		jsonData := v2ray.UriToJson(item, 3281+i)

		// فایل را باز کنید یا ایجاد کنید
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating file:", err)
			continue
		}

		// نوشتن محتوای JSON به فایل
		if _, err := io.WriteString(file, jsonData); err != nil {
			fmt.Println("Error writing to file:", err)
			file.Close()
			continue
		}

		// بسته کردن فایل بعد از نوشتن
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
			continue
		}

		fmt.Println("Created file:", fileName)
	}
}
