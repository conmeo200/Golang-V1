package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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
			return logs, nil // Return empty list if directory doesn't exist
		}
		return nil, fmt.Errorf("failed to read log directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories if any
		}

		info, err := entry.Info()
		if err != nil {
			continue // Skip if we can't get file info
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
	// Security check to prevent path traversal
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
