package cache

import (
	"os"
	"time"
)

func updateModTime(filePath string) error {
	file, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	err = os.Chtimes(filePath, file.ModTime(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func mkdir(dir string) (err error) {
	err = os.MkdirAll(dir, 0777)

	return
}

func parseDuration(cacheMaxAge string) time.Duration {
	duration, err := time.ParseDuration(cacheMaxAge)
	if err != nil {
		return time.Hour
	}
	return duration
}
