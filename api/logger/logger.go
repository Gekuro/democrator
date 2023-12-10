package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func NewLoggerWriter() (io.Writer, error) {
	if _, err := os.Stat(os.Getenv("LOGS_DIR")); os.IsNotExist(err) {
		err := os.MkdirAll(os.Getenv("LOGS_DIR"), os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("error opening/creating logs directory: %s", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("error opening/creating logs directory: %s", err)
	}

	year, month, day := time.Now().Date()
	logPath := fmt.Sprintf("%s/%v-%v-%v.log", os.Getenv("LOGS_DIR"), day, month, year)

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %s", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	return mw, nil
}