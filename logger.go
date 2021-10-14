package sap_segmentation

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

const (
	logFilename = "log/segmentation_import.log"
)

var (
	sevenDays = time.Hour * 24 * 7
)

func NewLogger(maxAge time.Duration) (*log.Logger, error) {

	if err := removeOld(logFilename, maxAge); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(
		logFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return log.New(
		io.MultiWriter(os.Stdout, logFile),
		"",
		log.LstdFlags), nil
}

func removeOld(fileName string, maxAge time.Duration) error {

	if err := ensureDir(fileName); err != nil {
		return err
	}

	info, err := os.Lstat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s is directory", fileName)
	}

	if time.Since(info.ModTime()) > maxAge {
		return os.Remove(fileName)
	}

	return nil
}

func ensureDir(filePath string) error {
	dir := path.Dir(filePath)
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}
