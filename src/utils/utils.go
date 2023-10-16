package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"time"
)

func GetTimestampFileName(index int, fileExtention string) string {
	var t time.Time = time.Now()
	return fmt.Sprintf("%d.%02d.%02d.%02d.%02d.%02d.%03d.%s.txt",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), index+1, fileExtention)
}

func GetTimestampedFileName(dir string, filename string) string {
	var (
		t    time.Time = time.Now()
		path string
	)
	filename = filepath.Base(filename)
	path = filepath.Join(dir, fmt.Sprintf("%d.%02d.%02d.%02d.%02d.%02d-%s",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), filename))
	log.Printf("ut> Generating new timestamp filename:\ndir: %s\nfilename: %s\ngenerated: %s", dir, filename, path)
	return path
}
