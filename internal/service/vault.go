package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Uvaisbugh/gophervault/internal/models"
	"github.com/Uvaisbugh/gophervault/internal/storage"
)

var (
	ErrEmptyTitle = errors.New("title cannot be empty")
	ErrEmptyBody  = errors.New("body cannot be empty")
)

// TagCount represents a tag and its usage count.
type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// VaultStats contains aggregate metrics for saved entries.
type VaultStats struct {
	TotalEntries int               `json:"total_entries"`
	TotalTags    int               `json:"total_tags"`
	Categories   map[string]int    `json:"categories"`
	PopularTags  []TagCount        `json:"popular_tags"`
}

// VaultService contains the business logic for managing knowledge notes.
//
// 📚 Go concept: Dependency Injection. VaultService accepts the storage.Repository
// interface, not a concrete struct. This keeps the service decoupled from disk details.
type VaultService struct {
	repo storage.Repository
}

// NewVaultService constructs a new VaultService.
func NewVaultService(repo storage.Repository) *VaultService {
	return &VaultService{repo: repo}
}

// AddEntry validates inputs, generates a new entry, and persists it.
func (s *VaultService) AddEntry(title, body, category string, tags []string) (*models.Entry, error) {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)

	if title == "" {
		return nil, ErrEmptyTitle
	}
	if body == "" {
		return nil, ErrEmptyBody
	}

	if category == "" {
		category = "general"
	}

	// Clean up tags (trim and remove empties)
	cleanTags := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			cleanTags = append(cleanTags, t)
		}
	}

	entry := models.NewEntry(title, body, category, cleanTags)
	if err := s.repo.Save(entry); err != nil {
		return nil, fmt.Errorf("could not save entry: %w", err)
	}

	return entry, nil
}

// GetEntry retrieves an entry by exact ID or prefix.
func (s *VaultService) GetEntry(id string) (*models.Entry, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.FindByID(id)
}

// SearchEntries searches all entries for matches in title, body, or tags.
func (s *VaultService) SearchEntries(query string) ([]*models.Entry, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}

	all, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*models.Entry
	for _, e := range all {
		if e.MatchesSearch(query) {
			results = append(results, e)
		}
	}

	return results, nil
}

// ListEntries filters entries by category and/or tag. If both are empty, returns all.
func (s *VaultService) ListEntries(category, tag string) ([]*models.Entry, error) {
	all, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	category = strings.ToLower(strings.TrimSpace(category))
	tag = strings.TrimSpace(tag)

	var filtered []*models.Entry
	for _, e := range all {
		if category != "" && strings.ToLower(e.Category) != category {
			continue
		}
		if tag != "" && !e.HasTag(tag) {
			continue
		}
		filtered = append(filtered, e)
	}

	return filtered, nil
}

// DeleteEntry deletes an entry by ID or prefix.
func (s *VaultService) DeleteEntry(id string) error {
	return s.repo.Delete(id)
}

// GetStats computes statistics about stored entries, tags, and categories.
//
// 📚 Go concept: Maps and sorting custom structs using sort.Slice.
func (s *VaultService) GetStats() (*VaultStats, error) {
	all, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	tagFrequency := make(map[string]int)
	categoryCount := make(map[string]int)

	for _, e := range all {
		cat := e.Category
		if cat == "" {
			cat = "general"
		}
		categoryCount[cat]++

		for _, t := range e.Tags {
			tagFrequency[strings.ToLower(t)]++
		}
	}

	popular := make([]TagCount, 0, len(tagFrequency))
	for t, count := range tagFrequency {
		popular = append(popular, TagCount{Tag: t, Count: count})
	}

	// 📚 Go concept: sort.Slice sorts using a comparator function
	sort.Slice(popular, func(i, j int) bool {
		if popular[i].Count == popular[j].Count {
			return popular[i].Tag < popular[j].Tag
		}
		return popular[i].Count > popular[j].Count
	})

	return &VaultStats{
		TotalEntries: len(all),
		TotalTags:    len(tagFrequency),
		Categories:   categoryCount,
		PopularTags:  popular,
	}, nil
}

// Export dumps all entries to an external JSON file backup.
func (s *VaultService) Export(targetPath string) error {
	all, err := s.repo.FindAll()
	if err != nil {
		return fmt.Errorf("failed to fetch entries: %w", err)
	}

	data, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode entries: %w", err)
	}

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	return nil
}

// Import restores or merges entries from a JSON backup file.
// If merge is false, existing items with the same ID will be overwritten or appended.
func (s *VaultService) Import(sourcePath string) (int, error) {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read import file: %w", err)
	}

	var imported []*models.Entry
	if err := json.Unmarshal(data, &imported); err != nil {
		return 0, fmt.Errorf("invalid backup file format: %w", err)
	}

	count := 0
	for _, entry := range imported {
		if strings.TrimSpace(entry.Title) == "" || strings.TrimSpace(entry.Body) == "" {
			continue // skip corrupted or malformed items
		}
		if err := s.repo.Save(entry); err != nil {
			return count, fmt.Errorf("failed to save entry during import: %w", err)
		}
		count++
	}

	return count, nil
}
