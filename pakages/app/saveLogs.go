package app

import (
	"log"
	"os"
)

func StartSaveLogs() {
	// Create or open the log file and truncate its content
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Set the output of logs to the file
	log.SetOutput(logFile)

	// Example log entries
	log.Println("This is an info message")
	log.Println("This is another log message")
}
