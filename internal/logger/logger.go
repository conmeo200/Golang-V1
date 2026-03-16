package logger

import (
	"log"
	"os"
	"path/filepath"
)

var (
	AppLogger   *log.Logger
	ErrorLogger *log.Logger
)

func Init() {
	// 1. Create log directory if not exists
	err := os.MkdirAll("log", os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// 2. Open app.log
	appFile, err := os.OpenFile(filepath.Join("log", "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open app.log: %v", err)
	}

	// 3. Open error.log
	errorFile, err := os.OpenFile(filepath.Join("log", "error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open error.log: %v", err)
	}

	// 4. Initialize global loggers
	AppLogger = log.New(appFile, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
