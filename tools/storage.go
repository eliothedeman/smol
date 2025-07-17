package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Storage struct {
	basePath string
	mu       sync.RWMutex
}

func NewStorage(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s *Storage) Save(key string, data interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.basePath, key+".json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (s *Storage) Load(key string, result interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.basePath, key+".json")
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(result)
}

func (s *Storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.basePath, key+".json")
	return os.Remove(path)
}

func (s *Storage) List() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	files, err := filepath.Glob(filepath.Join(s.basePath, "*.json"))
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, file := range files {
		key := filepath.Base(file[:len(file)-5]) // Remove .json extension
		keys = append(keys, key)
	}

	return keys, nil
}

func (s *Storage) Info() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	files, _ := filepath.Glob(filepath.Join(s.basePath, "*.json"))
	return fmt.Sprintf("Storage: %d items in %s", len(files), s.basePath)
}
