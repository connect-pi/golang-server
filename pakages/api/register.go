package api

// import (
// 	"encoding/json" // اضافه کردن پکیج json
// 	"fmt"
// 	"net/http"
// 	"project/pakages/app"
// 	"project/pakages/clog"
// 	"strings"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func handleConnection(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Error while upgrading connection:", err)
// 		return
// 	}

// 	// ارسال پیام به صورت دوره‌ای
// 	go func() {
// 		for {
// 			time.Sleep(500 * time.Millisecond)

// 			// ساختار داده برای ارسال به کلاینت
// 			data := struct {
// 				IsRun  bool   `json:"isRun"`
// 				NewLog string `json:"newLog"`
// 			}{
// 				IsRun:  app.IsRun(),
// 				NewLog: strings.Join(clog.Logs, "\n"), // فرض کنید این متغیر شامل لاگ‌هاست
// 			}

// 			// تبدیل داده‌ها به JSON
// 			jsonData, err := json.Marshal(data)
// 			if err != nil {
// 				fmt.Println("Error marshaling JSON:", err)
// 				break
// 			}

// 			// ارسال داده‌های JSON به کلاینت
// 			err = conn.WriteMessage(websocket.TextMessage, jsonData)
// 			if err != nil {
// 				fmt.Println("Error while writing message:", err)
// 				break
// 			}
// 		}
// 	}()

// 	for {
// 		messageType, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Error while reading message:", err)
// 			break
// 		}

// 		// GET
// 		if string(msg) == "start" { // START
// 			fmt.Println("-- START --")
// 			go app.Start()
// 		} else if string(msg) == "stop" { // STOP
// 			fmt.Println("-- STOP --")
// 			app.Stop()
// 		}

// 		fmt.Printf("Received message: %s\n", msg)

// 		// SEND
// 		err = conn.WriteMessage(messageType, msg)
// 		if err != nil {
// 			fmt.Println("Error while writing message:", err)
// 			break
// 		}
// 	}

// 	conn.Close()
// }

// func handleAdmin(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "html/admin.html")
// }

// func Register() {
// 	http.HandleFunc("/ws", handleConnection)
// 	http.HandleFunc("/admin", handleAdmin)
// 	fmt.Println("✨ Admin server started on 0.0.0.0:8080")

// 	err := http.ListenAndServe(":8080", nil)

// 	if err != nil {
// 		fmt.Println("Error starting server:", err)
// 	}
// }
