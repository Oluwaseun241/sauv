package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func SetupLogger(logDir string) (*log.Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return log.New(f, "", log.LstdFlags), nil
} 