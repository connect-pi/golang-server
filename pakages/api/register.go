package api

import (
	"fmt"
	"net/http"
	"project/pakages/app"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}

	fmt.Println("Client connected")

	// SEND message every time
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			err := conn.WriteMessage(websocket.TextMessage, []byte("{isRun:"+fmt.Sprint(app.IsRun())+", newLog: '-'}"))
			if err != nil {
				fmt.Println("Error while writing message:", err)
				break
			}
		}
	}()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}

		// GET
		if string(msg) == "start" { // START
			fmt.Println("-- START --")
			go app.Start()
		} else if string(msg) == "stop" { // STOP
			fmt.Println("-- STOP --")
			app.Stop()
		}

		fmt.Printf("Received message: %s\n", msg)

		// SEND
		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			fmt.Println("Error while writing message:", err)
			break
		}
	}

	conn.Close()
}

func Register() {
	http.HandleFunc("/ws", handleConnection)
	fmt.Println("Server started on :8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
