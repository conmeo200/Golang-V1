package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	AppLogger    *log.Logger
	ErrorLogger  *log.Logger
	StripeLogger *log.Logger
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

	// 4. Open stripe.log
	stripeFile, err := os.OpenFile(filepath.Join("log", "stripe.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open stripe.log: %v", err)
	}

	// 5. Initialize global loggers
	AppLogger = log.New(appFile, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	StripeLogger = log.New(stripeFile, "STRIPE: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogFileInfo holds metadata about a log file
type LogFileInfo struct {
	Filename  string
	Size      int64
	UpdatedAt time.Time
}

// LogServiceInterface defines the interface for log management
type LogServiceInterface interface {
	ListLogs() ([]LogFileInfo, error)
	GetLogContent(filename string) (string, error)
}

// FileLogService implements LogServiceInterface reading from a local directory
type FileLogService struct {
	logDirectory string
}

// NewFileLogService creates a new FileLogService
func NewFileLogService(logDirectory string) *FileLogService {
	return &FileLogService{
		logDirectory: logDirectory,
	}
}

// ListLogs returns a list of log files in the configured directory
func (s *FileLogService) ListLogs() ([]LogFileInfo, error) {
	var logs []LogFileInfo

	entries, err := os.ReadDir(s.logDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return logs, nil 
		}
		return nil, fmt.Errorf("failed to read log directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue 
		}

		info, err := entry.Info()
		if err != nil {
			continue 
		}

		logs = append(logs, LogFileInfo{
			Filename:  entry.Name(),
			Size:      info.Size(),
			UpdatedAt: info.ModTime(),
		})
	}

	return logs, nil
}

// GetLogContent returns the content of a specific log file
func (s *FileLogService) GetLogContent(filename string) (string, error) {
	cleanPath := filepath.Clean(filepath.Join(s.logDirectory, filename))
	if filepath.Dir(cleanPath) != filepath.Clean(s.logDirectory) {
		return "", fmt.Errorf("invalid log filename")
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("log file not found")
		}
		return "", fmt.Errorf("failed to read log file: %w", err)
	}

	return string(content), nil
}
