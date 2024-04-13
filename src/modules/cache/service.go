package cache

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"anthropic-discord-bot/src/logger"
)

const (
	attachmentsDir = "attachments"
	messagesDir    = "messages"
)

type Service struct {
	logger *logger.Logger

	cacheDir    string
	cacheMaxAge time.Duration
}

func (s *Service) saveCache(subdir []string, filename string, content *[]byte) (err error) {
	dirPath := filepath.Join(append([]string{s.cacheDir}, subdir...)...)

	err = mkdir(dirPath)
	if err != nil {
		return
	}

	filePath := filepath.Join(dirPath, filename)

	err = os.WriteFile(filePath, *content, 777)

	return
}

func (s *Service) getCache(subdir []string, filename string) (result *[]byte, err error) {
	filePath := filepath.Join(append(append([]string{s.cacheDir}, subdir...), filename)...)

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = updateModTime(filePath)
		if err != nil {
			s.logger.Error("Error updating ModTime", err)
		}
	}()

	return &bytes, nil
}

func (s *Service) SaveAttachment(id string, content *[]byte) {
	err := s.saveCache([]string{attachmentsDir}, id, content)

	if err != nil {
		s.logger.Error("Error saving attachment "+id+" to cache", err)
		return
	}

	s.logger.Debug("Attachment " + id + " saved to cache")
}

func (s *Service) GetAttachment(id string) (content *[]byte) {
	content, err := s.getCache([]string{attachmentsDir}, id)
	if err != nil {
		if !os.IsNotExist(err) {
			s.logger.Error("Error get attachment "+id+" from cache", err)
		}
		return
	}

	s.logger.Debug("Attachment " + id + " loaded from cache")
	return
}

func (s *Service) SaveMessage(channelID, messageID string, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		s.logger.Error("Error marshal message "+messageID+" to cache", err)
		return
	}

	err = s.saveCache([]string{messagesDir, channelID}, messageID, &data)

	if err != nil {
		s.logger.Error("Error saving message "+messageID+" to cache", err)
		return
	}

	s.logger.Debug("Message " + messageID + " saved to cache")
}

func (s *Service) GetMessage(channelID, messageID string, message interface{}) {
	content, err := s.getCache([]string{messagesDir, channelID}, messageID)
	if err != nil {
		if !os.IsNotExist(err) {
			s.logger.Error("Error get message "+messageID+" from cache", err)
		}
		return
	}

	err = json.Unmarshal(*content, message)
	if err != nil {
		s.logger.Error("Error unmarshal message "+messageID+" from cache", err)
		return
	}

	s.logger.Debug("Message " + messageID + " loaded from cache")
	return
}

func (s *Service) startCacheCleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		var cleared = 0
		var errors = 0

		s.logger.Debug("Clearing expired cache...")
		filepath.Walk(s.cacheDir, func(path string, info fs.FileInfo, err error) (_ error) {
			if err != nil {
				s.logger.Error("Error walking cache directory", err)
				errors++
				return
			}
			if info.IsDir() {
				return
			}
			if err != nil {
				s.logger.Error("Error getting file info", err)
				errors++
				return
			}
			if time.Since(info.ModTime()) > s.cacheMaxAge {
				err = os.Remove(path)
				if err != nil {
					s.logger.Error("Error removing expired cache file", err)
					errors++
					return
				}
				s.logger.Debug("Removed expired cache file: ", path)
				cleared++
			}
			return
		})

		s.logger.Debug("Expired cache cleared. Cleared: " + strconv.Itoa(cleared) + ". Errors: " + strconv.Itoa(errors))
	}
}
