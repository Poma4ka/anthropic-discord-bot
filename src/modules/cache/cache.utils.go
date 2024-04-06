package cache

import (
	"os"
	"time"
)

func mkdir(dir string) (err error) {
	err = os.MkdirAll(dir, 0600)

	return
}

func parseDuration(cacheMaxAge string) time.Duration {
	duration, err := time.ParseDuration(cacheMaxAge)
	if err != nil {
		return time.Hour
	}
	return duration
}
