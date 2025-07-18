package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestStorageInit(t *testing.T) {
	storage := NewStorage("/tmp/test_storage")
	ctx := &mockCtx{Context: context.Background()}
	storage.Init(ctx)
}

func TestStorageSaveLoad(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(tempDir)
	ctx := &mockCtx{Context: context.Background()}
	storage.Init(ctx)

	// Test saving and loading data
	testData := map[string]interface{}{
		"name":  "test",
		"value": 42,
		"list":  []int{1, 2, 3},
	}

	if err := storage.Save("test_key", testData); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var loadedData map[string]interface{}
	if err := storage.Load("test_key", &loadedData); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loadedData["name"] != "test" || loadedData["value"] != 42.0 {
		t.Errorf("Unexpected loaded data: %v", loadedData)
	}
}

func TestStorageSaveLoadString(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(tempDir)

	// Test saving and loading string data
	if err := storage.Save("string_key", "hello world"); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var loadedData string
	if err := storage.Load("string_key", &loadedData); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loadedData != "hello world" {
		t.Errorf("Expected 'hello world', got %v", loadedData)
	}
}

func TestStorageDelete(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(tempDir)

	// Save some data
	if err := storage.Save("to_delete", "data"); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Delete it
	if err := storage.Delete("to_delete"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it's gone
	var result string
	if err := storage.Load("to_delete", &result); err == nil {
		t.Error("Expected error when loading deleted key")
	}
}

func TestStorageList(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(tempDir)

	// Save multiple items
	items := []string{"item1", "item2", "item3"}
	for _, item := range items {
		if err := storage.Save(item, "data"); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	// List them
	keys, err := storage.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(keys) != 3 {
		t.Errorf("Expected 3 items, got %d", len(keys))
	}

	// Check if all items are present
	found := make(map[string]bool)
	for _, key := range keys {
		found[key] = true
	}
	for _, item := range items {
		if !found[item] {
			t.Errorf("Item %s not found in list", item)
		}
	}
}

func TestStorageInfo(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(tempDir)

	// Save some items
	for i := 0; i < 5; i++ {
		if err := storage.Save(fmt.Sprintf("item%d", i), "data"); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	info := storage.Info()
	if info == "" {
		t.Error("Info should not be empty")
	}
}

func TestStorageHandleStringCommands(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(tempDir)
	ctx := &mockCtx{Context: context.Background()}
	storage.Init(ctx)

	from := &testMessageHandler{}

	// Test save command
	err := storage.Handle(ctx, from, "save test_key {\"name\":\"test\",\"value\":42}")
	if err != nil {
		t.Errorf("Handle save command failed: %v", err)
	}

	// Test load command
	err = storage.Handle(ctx, from, "load test_key")
	if err != nil {
		t.Errorf("Handle load command failed: %v", err)
	}

	// Test list command
	err = storage.Handle(ctx, from, "list")
	if err != nil {
		t.Errorf("Handle list command failed: %v", err)
	}

	// Test delete command
	err = storage.Handle(ctx, from, "delete test_key")
	if err != nil {
		t.Errorf("Handle delete command failed: %v", err)
	}

	// Test help command
	err = storage.Handle(ctx, from, "help")
	if err != nil {
		t.Errorf("Handle help command failed: %v", err)
	}
}

func TestStorageHandleMapCommands(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(tempDir)
	ctx := &mockCtx{Context: context.Background()}
	storage.Init(ctx)

	from := &testMessageHandler{}

	// Test save action
	err := storage.Handle(ctx, from, map[string]interface{}{
		"action": "save",
		"key":    "test_key",
		"value":  map[string]interface{}{"name": "test", "value": 42},
	})
	if err != nil {
		t.Errorf("Handle save action failed: %v", err)
	}

	// Test load action
	err = storage.Handle(ctx, from, map[string]interface{}{
		"action": "load",
		"key":    "test_key",
	})
	if err != nil {
		t.Errorf("Handle load action failed: %v", err)
	}

	// Test list action
	err = storage.Handle(ctx, from, map[string]interface{}{
		"action": "list",
	})
	if err != nil {
		t.Errorf("Handle list action failed: %v", err)
	}
}

func TestStorageErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(tempDir)
	ctx := &mockCtx{Context: context.Background()}
	storage.Init(ctx)

	from := &testMessageHandler{}

	// Test invalid command
	err := storage.Handle(ctx, from, "invalid command")
	if err == nil {
		t.Error("Expected invalid command error")
	}

	// Test missing action
	err = storage.Handle(ctx, from, map[string]interface{}{
		"key": "test",
	})
	if err == nil {
		t.Error("Expected missing action error")
	}

	// Test load non-existent key
	err = storage.Handle(ctx, from, "load non_existent")
	if err == nil {
		t.Error("Expected load error for non-existent key")
	}
}

func TestStorageSaveStringCommand(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(tempDir)
	ctx := &mockCtx{Context: context.Background()}
	storage.Init(ctx)

	from := &testMessageHandler{}

	// Test saving string data
	err := storage.Handle(ctx, from, "save greeting Hello World")
	if err != nil {
		t.Errorf("Handle save string command failed: %v", err)
	}

	// Verify it was saved
	var result string
	if err := storage.Load("greeting", &result); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", result)
	}
}

func TestStorageDirectoryCreation(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewStorage(filepath.Join(tempDir, "subdir", "storage"))

	// Save should create directories
	if err := storage.Save("test", "data"); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	path := filepath.Join(tempDir, "subdir", "storage", "test.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("File %s should exist", path)
	}
}
