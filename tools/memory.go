package tools

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Memory struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		data: make(map[string]interface{}),
	}
}

func (m *Memory) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *Memory) Get(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, exists := m.data[key]
	return val, exists
}

func (m *Memory) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *Memory) List() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]interface{})
	for k, v := range m.data {
		result[k] = v
	}
	return result
}

func (m *Memory) Save() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

func (m *Memory) Load(data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return json.Unmarshal(data, &m.data)
}

func (m *Memory) Info() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("Memory: %d items", len(m.data))
}
