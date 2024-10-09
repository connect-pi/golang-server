package api

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"project/pakages/app" // Ensure this path is correct
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// // upgrader is used to upgrade an HTTP connection to a WebSocket connection.
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // Allow connections from any origin
// 	},
// }

// var (
// 	logMutex   sync.Mutex          // Mutex to synchronize access to shared resources
// 	clients    []*websocket.Conn   // Slice to hold connected WebSocket clients
// 	logChannel = make(chan string) // Channel for log messages
// )

// // startProcess starts the log process and sends a log message to clients.
// func startProcess() {
// 	fmt.Println("startProcess")
// 	logMutex.Lock()
// 	defer logMutex.Unlock()

// 	if app.IsRun() {
// 		fmt.Println("Process is already running.")
// 		return
// 	}

// 	fmt.Println("Starting process...")
// 	app.Run()
// 	logChannel <- "START!" // Send log message to channel
// }

// // stopProcess stops the log process and sends a log message to clients.
// func stopProcess() {
// 	fmt.Println("stopProcess")
// 	logMutex.Lock()
// 	defer logMutex.Unlock()

// 	if !app.IsRun() {
// 		fmt.Println("Process is not running.")
// 		return
// 	}

// 	fmt.Println("Stopping process...")
// 	app.Stop()
// 	logChannel <- "STOP!" // Send log message to channel
// }

// // broadcastLog sends a log message to all connected clients.
// func broadcastLog(message string) {
// 	logMutex.Lock()
// 	defer logMutex.Unlock()

// 	// Send the log message to all clients
// 	for _, conn := range clients {
// 		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
// 			log.Println("Write error:", err) // Print error if sending fails
// 		}
// 	}
// }

// // handleWebSocket manages WebSocket connections and messages.
// func handleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil) // Upgrade HTTP connection to WebSocket
// 	if err != nil {
// 		log.Println("Upgrade error:", err) // Log upgrade errors
// 		return
// 	}
// 	defer conn.Close() // Ensure connection is closed when done

// 	logMutex.Lock()
// 	clients = append(clients, conn) // Add the new client to the list
// 	logMutex.Unlock()

// 	// Send the current process status to the new client
// 	logMutex.Lock()
// 	if app.IsRun() {
// 		conn.WriteMessage(websocket.TextMessage, []byte("Process is currently running"))
// 	} else {
// 		conn.WriteMessage(websocket.TextMessage, []byte("No process is running"))
// 	}
// 	logMutex.Unlock()

// 	for {
// 		// Wait for messages from the client
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Read error:", err) // Log read errors
// 			break
// 		}

// 		fmt.Println(msg)

// 		// Check if the start or stop command was received
// 		if string(msg) == "start" {
// 			go startProcess() // Start the log process in a new goroutine
// 		} else if string(msg) == "stop" {
// 			stopProcess() // Stop the log process
// 		}
// 	}

// 	// Remove the client from the list after the connection is closed
// 	logMutex.Lock()
// 	for i, c := range clients {
// 		if c == conn {
// 			clients = append(clients[:i], clients[i+1:]...) // Remove client from slice
// 			break
// 		}
// 	}
// 	logMutex.Unlock()
// }

// // logSender listens for log messages and broadcasts them to clients.
// func logSender() {
// 	for {
// 		message := <-logChannel                        // Receive log messages from channel
// 		log.Println("Log message received: ", message) // Log the received message
// 		broadcastLog(message)                          // Broadcast the message to clients
// 	}
// }

// // serveHTML serves the HTML page for the WebSocket client.
// func serveHTML(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "run.html") // Serve the run.html file
// }

// // Run starts the HTTP server and handles incoming requests.
// func Run() {
// 	fmt.Println("Starting server...")
// 	go logSender()                          // Start the log sender in a new goroutine
// 	http.HandleFunc("/", serveHTML)         // Handle requests to the root URL
// 	http.HandleFunc("/ws", handleWebSocket) // Handle WebSocket requests

// 	fmt.Println("Server started at :8080")       // Print server start message
// 	log.Fatal(http.ListenAndServe(":8080", nil)) // Start the server
// }
