package cache

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"anthropic-discord-bot/src/logger"
)

type Service struct {
	logger *logger.Logger

	cacheDir    string
	cacheMaxAge time.Duration
}

func (s *Service) SaveAttachment(id string, content *[]byte) {
	s.logger.Debug("Saving attachment " + id + " to cache")
	filePath := filepath.Join(s.cacheDir, id)
	err := os.WriteFile(filePath, *content, 0600)
	if err != nil {
		s.logger.Error("Error saving attachment to cache", err)
		return
	}
	return
}

func (s *Service) GetAttachment(id string) (content *[]byte) {

	filePath := filepath.Join(s.cacheDir, id)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			s.logger.Error("Error getting attachment from cache", err)
		}
		return
	}

	defer func() {
		file, err := os.Stat(filePath)
		if err != nil {
			s.logger.Error("Error getting file info", err)
		}
		err = os.Chtimes(filePath, file.ModTime(), time.Now())
		if err != nil {
			s.logger.Error("Error updating ModTime", err)
		}
	}()

	s.logger.Debug("Loading attachment " + id + " from cache")

	return &bytes
}

func (s *Service) startCacheCleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.logger.Debug("Clearing expired cache...")
		err := filepath.Walk(s.cacheDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				s.logger.Error("Error walking cache directory", err)
				return err
			}
			if info.IsDir() {
				return nil
			}
			if err != nil {
				s.logger.Error("Error getting file info", err)
				return err
			}
			if time.Since(info.ModTime()) > s.cacheMaxAge {
				err = os.Remove(path)
				if err != nil {
					s.logger.Error("Error removing expired cache file", err)
					return err
				}
				s.logger.Debug("Removed expired cache file: ", path)
			}
			return nil
		})
		if err != nil {
			s.logger.Error("Error clearing expired cache", err)
		}
		s.logger.Debug("Expired cache cleared")
	}
}
