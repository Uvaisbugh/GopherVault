package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Uvaisbugh/gophervault/internal/models"
)

// 📚 Go concept: Declaring sentinel errors using errors.New().
// This allows callers to check specifically for errors like ErrNotFound using errors.Is().
var (
	ErrNotFound = errors.New("entry not found")
)

// JSONStore implements the Repository interface by saving records to a JSON file.
//
// 📚 Go concept: sync.RWMutex provides read-write thread safety.
// Even in simple CLIs, handling concurrent reads/writes is good Go practice.
type JSONStore struct {
	filePath string
	mu       sync.RWMutex
}

// NewJSONStore creates a JSONStore. If filePath is empty, it defaults to:
// ~/.gophervault/vault.json in the user's home directory.
func NewJSONStore(filePath string) (*JSONStore, error) {
	if filePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not locate user home directory: %w", err)
		}
		dir := filepath.Join(homeDir, ".gophervault")
		// 📚 Go concept: os.MkdirAll creates parent directories if they don't exist.
		// 0755 is POSIX file permissions: rwxr-xr-x.
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create vault directory: %w", err)
		}
		filePath = filepath.Join(dir, "vault.json")
	} else {
		// Ensure directory for custom path exists
		dir := filepath.Dir(filePath)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory for %s: %w", filePath, err)
			}
		}
	}

	store := &JSONStore{filePath: filePath}
	return store, nil
}

// FilePath returns the underlying storage file path.
func (s *JSONStore) FilePath() string {
	return s.filePath
}

// load reads and parses the JSON file.
//
// 📚 Go concept: Unexported helper method.
// Notice defer: s.mu.RUnlock() runs when the function returns!
func (s *JSONStore) load() ([]*models.Entry, error) {
	// If the file doesn't exist yet, return an empty slice rather than an error.
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []*models.Entry{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage file: %w", err)
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		return []*models.Entry{}, nil
	}

	var entries []*models.Entry
	// 📚 Go concept: json.Unmarshal converts raw JSON bytes into Go structs.
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %w", s.filePath, err)
	}

	return entries, nil
}

// persist writes entries to disk as pretty-printed JSON.
func (s *JSONStore) persist(entries []*models.Entry) error {
	// 📚 Go concept: json.MarshalIndent outputs formatted, indented JSON.
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize entries to JSON: %w", err)
	}

	// 0644 gives read/write to owner, read-only to others.
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write to storage file: %w", err)
	}

	return nil
}

// Save inserts a new entry or updates an existing one if the ID matches.
func (s *JSONStore) Save(entry *models.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := s.load()
	if err != nil {
		return err
	}

	updated := false
	for i, existing := range entries {
		if existing.ID == entry.ID {
			entries[i] = entry
			updated = true
			break
		}
	}

	if !updated {
		// 📚 Go concept: append() adds an element to a slice and returns the updated slice.
		entries = append(entries, entry)
	}

	return s.persist(entries)
}

// FindByID searches for an entry by exact ID or prefix match (e.g. first 6-8 chars).
func (s *JSONStore) FindByID(id string) (*models.Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := s.load()
	if err != nil {
		return nil, err
	}

	idLower := strings.ToLower(id)
	var matches []*models.Entry

	for _, entry := range entries {
		if strings.ToLower(entry.ID) == idLower {
			return entry, nil
		}
		if strings.HasPrefix(strings.ToLower(entry.ID), idLower) {
			matches = append(matches, entry)
		}
	}

	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("ambiguous id prefix '%s', matched %d entries", id, len(matches))
	}

	return nil, ErrNotFound
}

// FindAll returns all entries.
func (s *JSONStore) FindAll() ([]*models.Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.load()
}

// Delete removes an entry with matching ID or ID prefix.
func (s *JSONStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := s.load()
	if err != nil {
		return err
	}

	targetIndex := -1
	idLower := strings.ToLower(id)

	for i, entry := range entries {
		if strings.ToLower(entry.ID) == idLower || strings.HasPrefix(strings.ToLower(entry.ID), idLower) {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return ErrNotFound
	}

	// 📚 Go concept: Deleting an element from a slice in Go.
	// append(slice[:i], slice[i+1:]...) takes elements before i and appends
	// elements after i. The "..." unpacks the sub-slice into arguments.
	entries = append(entries[:targetIndex], entries[targetIndex+1:]...)

	return s.persist(entries)
}
