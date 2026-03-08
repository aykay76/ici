package secrets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileStore_SetAndGet(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "ici-secrets-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	secretFile := filepath.Join(tmpDir, "secrets.json")
	store := NewFileStore(secretFile)

	// Test Set and Get
	err = store.Set("MY_SECRET", "super-secret-value")
	if err != nil {
		t.Fatalf("failed to set secret: %v", err)
	}

	value, err := store.Get("MY_SECRET")
	if err != nil {
		t.Fatalf("failed to get secret: %v", err)
	}

	if value != "super-secret-value" {
		t.Errorf("expected 'super-secret-value', got '%s'", value)
	}

	// Test Get non-existent secret
	value, err = store.Get("NON_EXISTENT")
	if err != nil {
		t.Fatalf("failed to get non-existent secret: %v", err)
	}
	if value != "" {
		t.Errorf("expected empty string for non-existent secret, got '%s'", value)
	}
}

func TestFileStore_Delete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ici-secrets-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	secretFile := filepath.Join(tmpDir, "secrets.json")
	store := NewFileStore(secretFile)

	// Set a secret
	err = store.Set("MY_SECRET", "value")
	if err != nil {
		t.Fatalf("failed to set secret: %v", err)
	}

	// Delete it
	err = store.Delete("MY_SECRET")
	if err != nil {
		t.Fatalf("failed to delete secret: %v", err)
	}

	// Verify it's gone
	value, err := store.Get("MY_SECRET")
	if err != nil {
		t.Fatalf("failed to get secret: %v", err)
	}
	if value != "" {
		t.Errorf("expected empty string after delete, got '%s'", value)
	}
}

func TestFileStore_List(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ici-secrets-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	secretFile := filepath.Join(tmpDir, "secrets.json")
	store := NewFileStore(secretFile)

	// Set multiple secrets
	store.Set("SECRET1", "value1")
	store.Set("SECRET2", "value2")
	store.Set("SECRET3", "value3")

	// List secrets
	names, err := store.List()
	if err != nil {
		t.Fatalf("failed to list secrets: %v", err)
	}

	if len(names) != 3 {
		t.Errorf("expected 3 secrets, got %d", len(names))
	}

	// Check that all expected names are present
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	for _, expected := range []string{"SECRET1", "SECRET2", "SECRET3"} {
		if !nameMap[expected] {
			t.Errorf("expected secret '%s' in list", expected)
		}
	}
}

func TestFileStore_GetAll(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ici-secrets-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	secretFile := filepath.Join(tmpDir, "secrets.json")
	store := NewFileStore(secretFile)

	// Set multiple secrets
	store.Set("KEY1", "val1")
	store.Set("KEY2", "val2")

	// Get all
	all, err := store.GetAll()
	if err != nil {
		t.Fatalf("failed to get all secrets: %v", err)
	}

	if len(all) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(all))
	}

	if all["KEY1"] != "val1" {
		t.Errorf("expected KEY1=val1, got KEY1=%s", all["KEY1"])
	}

	if all["KEY2"] != "val2" {
		t.Errorf("expected KEY2=val2, got KEY2=%s", all["KEY2"])
	}
}

func TestFileStore_EmptyStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ici-secrets-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	secretFile := filepath.Join(tmpDir, "secrets.json")
	store := NewFileStore(secretFile)

	// File doesn't exist yet, should return empty
	names, err := store.List()
	if err != nil {
		t.Fatalf("failed to list empty secrets: %v", err)
	}

	if len(names) != 0 {
		t.Errorf("expected 0 secrets, got %d", len(names))
	}
}

func TestFileStore_Persistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ici-secrets-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	secretFile := filepath.Join(tmpDir, "secrets.json")

	// Create first store and set a secret
	store1 := NewFileStore(secretFile)
	err = store1.Set("PERSISTENT", "persisted-value")
	if err != nil {
		t.Fatalf("failed to set secret: %v", err)
	}

	// Create a new store instance with same file
	store2 := NewFileStore(secretFile)
	value, err := store2.Get("PERSISTENT")
	if err != nil {
		t.Fatalf("failed to get secret from new store: %v", err)
	}

	if value != "persisted-value" {
		t.Errorf("expected 'persisted-value', got '%s'", value)
	}
}
