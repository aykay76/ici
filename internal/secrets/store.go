package secrets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Store defines the interface for secret storage backends
type Store interface {
	// Get retrieves a secret by name, returns empty string if not found
	Get(name string) (string, error)
	// Set stores a secret
	Set(name string, value string) error
	// Delete removes a secret
	Delete(name string) error
	// List returns all secret names
	List() ([]string, error)
	// GetAll returns all secrets as a map
	GetAll() (map[string]string, error)
}

// FileStore is a simple file-based secret store (JSON format)
type FileStore struct {
	filePath string
}

// NewFileStore creates a new file-based secret store
// Uses ~/.ici/secrets.json by default
func NewFileStore(filePath string) *FileStore {
	if filePath == "" {
		homeDir, _ := os.UserHomeDir()
		filePath = filepath.Join(homeDir, ".ici", "secrets.json")
	}
	return &FileStore{filePath: filePath}
}

// DefaultStore creates a file store at the default location
func DefaultStore() *FileStore {
	return NewFileStore("")
}

// Get retrieves a secret by name
func (fs *FileStore) Get(name string) (string, error) {
	secrets, err := fs.loadSecrets()
	if err != nil {
		return "", err
	}
	if value, exists := secrets[name]; exists {
		return value, nil
	}
	return "", nil // Return empty string if not found (not an error)
}

// Set stores a secret
func (fs *FileStore) Set(name string, value string) error {
	secrets, err := fs.loadSecrets()
	if err != nil {
		return err
	}
	secrets[name] = value
	return fs.saveSecrets(secrets)
}

// Delete removes a secret
func (fs *FileStore) Delete(name string) error {
	secrets, err := fs.loadSecrets()
	if err != nil {
		return err
	}
	delete(secrets, name)
	return fs.saveSecrets(secrets)
}

// List returns all secret names
func (fs *FileStore) List() ([]string, error) {
	secrets, err := fs.loadSecrets()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(secrets))
	for name := range secrets {
		names = append(names, name)
	}
	return names, nil
}

// GetAll returns all secrets as a map
func (fs *FileStore) GetAll() (map[string]string, error) {
	return fs.loadSecrets()
}

// loadSecrets reads secrets from the file
func (fs *FileStore) loadSecrets() (map[string]string, error) {
	secrets := make(map[string]string)

	// If file doesn't exist, return empty map
	if _, err := os.Stat(fs.filePath); os.IsNotExist(err) {
		return secrets, nil
	}

	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets file: %w", err)
	}

	if err := json.Unmarshal(data, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse secrets file: %w", err)
	}

	return secrets, nil
}

// saveSecrets writes secrets to the file
func (fs *FileStore) saveSecrets(secrets map[string]string) error {
	// Ensure directory exists
	if dir := filepath.Dir(fs.filePath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create secrets directory: %w", err)
		}
	}

	data, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal secrets: %w", err)
	}

	// Write with restricted permissions for security
	if err := os.WriteFile(fs.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write secrets file: %w", err)
	}

	return nil
}
