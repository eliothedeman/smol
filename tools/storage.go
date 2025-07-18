package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/eliothedeman/smol/unit"
)

type Storage struct {
	basePath string
	mu       sync.RWMutex
	ctx      unit.Ctx
}

func NewStorage(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s *Storage) Init(ctx unit.Ctx) {
	s.ctx = ctx
}

func (s *Storage) Handle(ctx unit.Ctx, from unit.UnitRef, message any) error {
	switch msg := message.(type) {
	case string:
		return s.handleStringCommand(ctx, from, msg)
	case map[string]interface{}:
		return s.handleMapCommand(ctx, from, msg)
	default:
		return fmt.Errorf("unsupported message type: %T", message)
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

func (s *Storage) handleStringCommand(ctx unit.Ctx, from unit.UnitRef, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "save":
		if len(parts) < 3 {
			return fmt.Errorf("save requires key and value")
		}
		key := parts[1]
		value := strings.Join(parts[2:], " ")

		// Try to parse as JSON first
		var data interface{}
		if err := json.Unmarshal([]byte(value), &data); err != nil {
			// If not JSON, store as string
			data = value
		}

		if err := s.Save(key, data); err != nil {
			return err
		}
		from.Send(fmt.Sprintf("Saved %s", key))

	case "load":
		if len(parts) != 2 {
			return fmt.Errorf("load requires key")
		}
		key := parts[1]
		var data interface{}
		if err := s.Load(key, &data); err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%s = %v", key, data))

	case "delete", "del":
		if len(parts) != 2 {
			return fmt.Errorf("delete requires key")
		}
		key := parts[1]
		if err := s.Delete(key); err != nil {
			return err
		}
		from.Send(fmt.Sprintf("Deleted %s", key))

	case "list":
		keys, err := s.List()
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			from.Send("No items")
		} else {
			result := "Stored items:\n"
			for _, key := range keys {
				result += fmt.Sprintf("  %s\n", key)
			}
			from.Send(result)
		}

	case "info":
		from.Send(s.Info())

	case "help":
		from.Send(`Storage commands:
  save <key> <value> - Save data
  load <key> - Load data
  delete <key> - Delete data
  list - List all keys
  info - Storage info
  help - Show this help`)

	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return nil
}

func (s *Storage) handleMapCommand(ctx unit.Ctx, from unit.UnitRef, command map[string]interface{}) error {
	action, ok := command["action"].(string)
	if !ok {
		return fmt.Errorf("missing action field")
	}

	switch action {
	case "save":
		key, ok1 := command["key"].(string)
		value, ok2 := command["value"]
		if !ok1 || !ok2 {
			return fmt.Errorf("save requires key and value")
		}
		if err := s.Save(key, value); err != nil {
			return err
		}
		from.Send(fmt.Sprintf("Saved %s", key))

	case "load":
		key, ok := command["key"].(string)
		if !ok {
			return fmt.Errorf("load requires key")
		}
		var data interface{}
		if err := s.Load(key, &data); err != nil {
			return err
		}
		from.Send(data)

	case "delete":
		key, ok := command["key"].(string)
		if !ok {
			return fmt.Errorf("delete requires key")
		}
		if err := s.Delete(key); err != nil {
			return err
		}
		from.Send(fmt.Sprintf("Deleted %s", key))

	case "list":
		keys, err := s.List()
		if err != nil {
			return err
		}
		from.Send(keys)

	case "info":
		from.Send(s.Info())

	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	return nil
}
